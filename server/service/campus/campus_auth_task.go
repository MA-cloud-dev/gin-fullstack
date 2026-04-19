package campus

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	campusResp "github.com/flipped-aurora/gin-vue-admin/server/model/campus/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	agentReviewRequestTokenHeader  = "X-Agent-Request-Token"
	agentReviewCallbackTokenHeader = "X-Agent-Callback-Token"
)

var campusStudentIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{4,32}$`)

type agentReviewDispatchRequest struct {
	TaskID      string                           `json:"taskId"`
	BizType     string                           `json:"bizType"`
	BizID       string                           `json:"bizId"`
	TraceID     string                           `json:"traceId"`
	CallbackURL string                           `json:"callbackUrl"`
	Payload     campusReq.AgentReviewTaskPayload `json:"payload"`
}

func (s *CampusAuthService) SubmitCampusAuth(ctx context.Context, req campusReq.SubmitCampusAuthReq) (campusResp.SubmitCampusAuthResp, error) {
	normalizedReq, err := normalizeSubmitCampusAuthReq(req)
	if err != nil {
		return campusResp.SubmitCampusAuthResp{}, err
	}

	var result campusResp.SubmitCampusAuthResp
	var task campusModel.AgentReviewTask
	var payload campusReq.AgentReviewTaskPayload

	err = global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user campusModel.CampusUser
		if err := tx.Table(user.TableName()).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", normalizedReq.UserID).
			First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("校园用户不存在")
			}
			return err
		}

		if hasActive, err := hasActiveCampusAuthTask(tx, normalizedReq.UserID, ""); err != nil {
			return err
		} else if hasActive {
			return errors.New("已有审核中的申请，请勿重复提交")
		}
		if blocked, message, err := shouldBlockCampusAuthResubmission(tx, normalizedReq.UserID); err != nil {
			return err
		} else if blocked {
			return errors.New(message)
		}

		previousStatus, err := computeCampusUserStableAuthStatus(tx, uint(normalizedReq.UserID))
		if err != nil {
			return err
		}
		if previousStatus == campusModel.CampusUserAuthStatusApproved {
			return errors.New("当前用户已通过校园认证，无需重复提交")
		}

		if err := supersedeFailedCampusAuthTasks(tx, normalizedReq.UserID); err != nil {
			return err
		}

		auth, err := upsertCampusAuthSubmissionRecord(tx, uint(normalizedReq.UserID), normalizedReq)
		if err != nil {
			return err
		}

		payload = buildCampusAuthTaskPayload(normalizedReq)
		payloadJSON, err := marshalJSONValue(payload)
		if err != nil {
			return err
		}
		callbackURL, err := buildCampusAuthCallbackURL()
		if err != nil {
			return err
		}
		taskID, traceID := newCampusAuthTaskIdentifiers()
		task = campusModel.AgentReviewTask{
			TaskID:                 taskID,
			BizType:                campusModel.AgentReviewBizTypeCampusAuth,
			BizID:                  uint64(auth.ID),
			TraceID:                traceID,
			CallbackURL:            callbackURL,
			PayloadJSON:            payloadJSON,
			TaskStatus:             campusModel.AgentReviewTaskStatusPendingDispatch,
			FinalAction:            "",
			ReviewRemark:           "",
			DispatchAttemptCount:   0,
			PreviousUserAuthStatus: previousStatus,
			AuthRecordID:           uint64(auth.ID),
			UserID:                 normalizedReq.UserID,
		}
		if err := tx.Table(task.TableName()).Create(&task).Error; err != nil {
			return err
		}
		if err := updateCampusUserAuthStatus(tx, auth.UserID, campusModel.CampusUserAuthStatusProcessing); err != nil {
			return err
		}

		result = campusResp.SubmitCampusAuthResp{
			AuthRecordID: uint64(auth.ID),
			TaskID:       task.TaskID,
			ReviewStatus: campusModel.CampusAuthReviewStatusProcessing,
			AuthStatus:   campusModel.CampusUserAuthStatusProcessing,
		}
		return nil
	})
	if err != nil {
		return campusResp.SubmitCampusAuthResp{}, err
	}

	if err := s.dispatchAgentReviewTask(ctx, task, payload); err != nil {
		if markErr := s.markAgentReviewTaskDispatchFailed(ctx, task.TaskID, err.Error()); markErr != nil {
			global.GVA_LOG.Error("标记校园认证任务派发失败状态异常", zap.Error(markErr))
		}
		return campusResp.SubmitCampusAuthResp{}, fmt.Errorf("提交审核任务失败: %w", err)
	}

	if err := s.markAgentReviewTaskSubmitted(ctx, task.TaskID); err != nil {
		return campusResp.SubmitCampusAuthResp{}, fmt.Errorf("外部审核任务已接单，但本地状态更新失败，请按 taskId=%s 补偿: %w", task.TaskID, err)
	}

	return result, nil
}

func (s *CampusAuthService) HandleAgentReviewCallback(ctx context.Context, callbackToken string, req campusReq.AgentReviewCallbackReq, auditMeta campusReq.CampusAuditMeta) error {
	expectedToken := strings.TrimSpace(global.GVA_CONFIG.AgentReview.CallbackToken)
	if expectedToken != "" && subtle.ConstantTimeCompare([]byte(expectedToken), []byte(strings.TrimSpace(callbackToken))) != 1 {
		return errors.New("callback token 校验失败")
	}

	req.Normalize()
	if err := req.Validate(); err != nil {
		return err
	}

	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var task campusModel.AgentReviewTask
		if err := tx.Table(task.TableName()).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("task_id = ?", req.TaskID).
			First(&task).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("任务不存在")
			}
			return err
		}

		if task.BizType != req.BizType {
			return errors.New("bizType 不匹配")
		}
		if strconv.FormatUint(task.BizID, 10) != req.BizID {
			return errors.New("bizId 不匹配")
		}
		if task.TraceID != req.TraceID {
			return errors.New("traceId 不匹配")
		}

		nextTaskStatus, nextFinalAction := mapCallbackToLocalTaskState(req)
		if isSameTerminalTaskState(task, nextTaskStatus, nextFinalAction) || isTerminalTaskStatus(task.TaskStatus) {
			return nil
		}

		resultJSON, err := marshalJSONValue(req.Result)
		if err != nil {
			return err
		}

		auth, err := loadCampusAuthForUpdate(tx, uint(task.AuthRecordID))
		if err != nil {
			return err
		}

		now := time.Now()
		taskUpdates := map[string]interface{}{
			"task_status":    nextTaskStatus,
			"final_action":   nextFinalAction,
			"review_remark":  req.ReviewRemark,
			"result_json":    resultJSON,
			"callback_at":    now,
			"callback_error": "",
		}
		authUpdates := map[string]interface{}{
			"review_source": campusModel.CampusAuthReviewSourceAgent,
		}

		action := "agent_review_processing"
		resultText := "Agent审核处理中"
		userAuthStatus := campusModel.CampusUserAuthStatusProcessing

		switch req.TaskStatus {
		case campusModel.AgentReviewCallbackTaskStatusProcessing:
			authUpdates["review_status"] = campusModel.CampusAuthReviewStatusProcessing
			if req.ReviewRemark != "" {
				authUpdates["review_remark"] = req.ReviewRemark
			}
		case campusModel.AgentReviewCallbackTaskStatusCompleted:
			switch req.FinalAction {
			case campusModel.AgentReviewFinalActionApprove:
				action = "agent_review_approve"
				resultText = "Agent审核通过"
				userAuthStatus = campusModel.CampusUserAuthStatusApproved
				authUpdates["review_status"] = campusModel.CampusAuthReviewStatusApproved
				authUpdates["review_remark"] = req.ReviewRemark
				authUpdates["reviewed_at"] = now
			case campusModel.AgentReviewFinalActionReject:
				action = "agent_review_reject"
				resultText = "Agent审核拒绝"
				userAuthStatus = campusModel.CampusUserAuthStatusRejected
				authUpdates["review_status"] = campusModel.CampusAuthReviewStatusRejected
				authUpdates["review_remark"] = req.ReviewRemark
				authUpdates["reviewed_at"] = now
			case campusModel.AgentReviewFinalActionEscalate:
				action = "agent_review_escalate"
				resultText = "Agent已转人工复核"
				userAuthStatus = campusModel.CampusUserAuthStatusProcessing
				authUpdates["review_status"] = campusModel.CampusAuthReviewStatusProcessing
				authUpdates["review_remark"] = req.ReviewRemark
				authUpdates["reviewed_at"] = nil
				authUpdates["reviewed_by"] = nil
			}
		case campusModel.AgentReviewCallbackTaskStatusFailed:
			action = "agent_review_failed"
			resultText = "Agent审核失败"
			taskUpdates["callback_error"] = buildCallbackErrorMessage(req)
			authUpdates["review_status"] = campusModel.CampusAuthReviewStatusProcessing
			userAuthStatus, err = resolveRollbackCampusUserAuthStatus(tx, uint(task.UserID), task.PreviousUserAuthStatus, task.TaskID)
			if err != nil {
				return err
			}
		}

		if err := tx.Table(task.TableName()).Where("id = ?", task.ID).Updates(taskUpdates).Error; err != nil {
			return err
		}
		if err := tx.Table(auth.TableName()).Where("id = ?", auth.ID).Updates(authUpdates).Error; err != nil {
			return err
		}
		if err := updateCampusUserAuthStatus(tx, auth.UserID, userAuthStatus); err != nil {
			return err
		}

		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "auth",
			Action:      action,
			TargetID:    auth.ID,
			TargetLabel: joinCampusAuditLabel(auth.RealName, auth.StudentID),
			Reason:      req.ReviewRemark,
			Result:      resultText,
		})
	})
}

func normalizeSubmitCampusAuthReq(req campusReq.SubmitCampusAuthReq) (campusReq.SubmitCampusAuthReq, error) {
	req.StudentID = strings.TrimSpace(req.StudentID)
	req.RealName = strings.TrimSpace(req.RealName)
	req.College = strings.TrimSpace(req.College)

	if req.UserID == 0 {
		return req, errors.New("userId 不能为空")
	}
	if req.StudentID == "" || req.RealName == "" || req.College == "" {
		return req, errors.New("学号、姓名、学院不能为空")
	}
	if !campusStudentIDPattern.MatchString(req.StudentID) {
		return req, errors.New("学号格式不正确，仅支持 4-32 位字母、数字、下划线或中划线")
	}
	if len([]rune(req.RealName)) > 32 {
		return req, errors.New("姓名长度不能超过 32 个字符")
	}
	if len([]rune(req.College)) > 64 {
		return req, errors.New("学院长度不能超过 64 个字符")
	}
	return req, nil
}

func buildCampusAuthTaskPayload(req campusReq.SubmitCampusAuthReq) campusReq.AgentReviewTaskPayload {
	return campusReq.AgentReviewTaskPayload{
		UserID:    strconv.FormatUint(req.UserID, 10),
		StudentID: req.StudentID,
		RealName:  req.RealName,
		College:   req.College,
	}
}

func buildCampusAuthCallbackURL() (string, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(global.GVA_CONFIG.AgentReview.PublicBaseURL), "/")
	callbackPath := strings.TrimSpace(global.GVA_CONFIG.AgentReview.CallbackPath)
	if baseURL == "" {
		return "", errors.New("agent-review.public-base-url 未配置")
	}
	if callbackPath == "" {
		return "", errors.New("agent-review.callback-path 未配置")
	}
	if !strings.HasPrefix(callbackPath, "/") {
		callbackPath = "/" + callbackPath
	}
	return baseURL + callbackPath, nil
}

func newCampusAuthTaskIdentifiers() (string, string) {
	now := time.Now()
	taskID := fmt.Sprintf("task_%s_%s", now.Format("20060102"), strings.ReplaceAll(uuid.NewString()[:8], "-", ""))
	traceID := "trace_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	return taskID, traceID
}

func marshalJSONValue(value any) (datatypes.JSON, error) {
	if value == nil {
		return datatypes.JSON([]byte("{}")), nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(raw), nil
}

func hasActiveCampusAuthTask(tx *gorm.DB, userID uint64, excludeTaskID string) (bool, error) {
	var count int64
	query := tx.Table(campusModel.AgentReviewTask{}.TableName()).
		Where("user_id = ? AND biz_type = ? AND task_status IN ?", userID, campusModel.AgentReviewBizTypeCampusAuth, []string{
			campusModel.AgentReviewTaskStatusPendingDispatch,
			campusModel.AgentReviewTaskStatusSubmitted,
			campusModel.AgentReviewTaskStatusProcessing,
		})
	if excludeTaskID != "" {
		query = query.Where("task_id <> ?", excludeTaskID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func supersedeFailedCampusAuthTasks(tx *gorm.DB, userID uint64) error {
	return tx.Table(campusModel.AgentReviewTask{}.TableName()).
		Where("user_id = ? AND biz_type = ? AND task_status = ?", userID, campusModel.AgentReviewBizTypeCampusAuth, campusModel.AgentReviewTaskStatusFailed).
		Updates(map[string]interface{}{
			"task_status": campusModel.AgentReviewTaskStatusSuperseded,
		}).Error
}

func shouldBlockCampusAuthResubmission(tx *gorm.DB, userID uint64) (bool, string, error) {
	auth, hasAuth, err := loadLatestCampusAuthSubmissionRecord(tx, uint(userID))
	if err != nil {
		return false, "", err
	}
	if !hasAuth || normalizeCampusAuthStatusWithLegacy(auth.ReviewStatus, auth.ReviewedAt) != campusModel.CampusAuthReviewStatusProcessing {
		return false, "", nil
	}

	task, hasTask, err := loadLatestCampusAuthTask(tx, userID)
	if err != nil {
		return false, "", err
	}
	if !hasTask {
		return true, "已有审核中的申请，请勿重复提交", nil
	}

	switch strings.TrimSpace(task.TaskStatus) {
	case campusModel.AgentReviewTaskStatusFailed, campusModel.AgentReviewTaskStatusSuperseded:
		return false, "", nil
	case campusModel.AgentReviewTaskStatusCompleted:
		if strings.TrimSpace(task.FinalAction) == campusModel.AgentReviewFinalActionEscalate {
			return true, "已有人工复核中的申请，请勿重复提交", nil
		}
		return true, "已有审核中的申请，请勿重复提交", nil
	default:
		return true, "已有审核中的申请，请勿重复提交", nil
	}
}

func loadLatestCampusAuthSubmissionRecord(tx *gorm.DB, userID uint) (campusModel.CampusAuth, bool, error) {
	var auth campusModel.CampusAuth
	err := tx.Table(auth.TableName()).
		Where("user_id = ?", userID).
		Order("id DESC").
		First(&auth).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return campusModel.CampusAuth{}, false, nil
	}
	if err != nil {
		return campusModel.CampusAuth{}, false, err
	}
	return auth, true, nil
}

func loadLatestCampusAuthTask(tx *gorm.DB, userID uint64) (campusModel.AgentReviewTask, bool, error) {
	var task campusModel.AgentReviewTask
	err := tx.Table(task.TableName()).
		Where("user_id = ? AND biz_type = ?", userID, campusModel.AgentReviewBizTypeCampusAuth).
		Order("id DESC").
		First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return campusModel.AgentReviewTask{}, false, nil
	}
	if err != nil {
		return campusModel.AgentReviewTask{}, false, err
	}
	return task, true, nil
}

func upsertCampusAuthSubmissionRecord(tx *gorm.DB, userID uint, req campusReq.SubmitCampusAuthReq) (campusModel.CampusAuth, error) {
	var auth campusModel.CampusAuth
	err := tx.Table(auth.TableName()).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		Order("id DESC").
		First(&auth).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return auth, err
	}

	if err == nil && normalizeCampusAuthStatusWithLegacy(auth.ReviewStatus, auth.ReviewedAt) == campusModel.CampusAuthReviewStatusProcessing {
		if err := tx.Table(auth.TableName()).Where("id = ?", auth.ID).Updates(map[string]interface{}{
			"student_id":    req.StudentID,
			"real_name":     req.RealName,
			"college":       req.College,
			"review_status": campusModel.CampusAuthReviewStatusProcessing,
			"review_source": nil,
			"review_remark": nil,
			"reviewed_at":   nil,
			"reviewed_by":   nil,
		}).Error; err != nil {
			return auth, err
		}
		auth.StudentID = req.StudentID
		auth.RealName = req.RealName
		auth.College = req.College
		auth.ReviewStatus = campusModel.CampusAuthReviewStatusProcessing
		auth.ReviewSource = nil
		auth.ReviewRemark = nil
		auth.ReviewedAt = nil
		auth.ReviewedBy = nil
		return auth, nil
	}

	auth = campusModel.CampusAuth{
		UserID:       userID,
		StudentID:    req.StudentID,
		RealName:     req.RealName,
		College:      req.College,
		ReviewStatus: campusModel.CampusAuthReviewStatusProcessing,
	}
	if err := tx.Table(auth.TableName()).Create(&auth).Error; err != nil {
		return auth, err
	}
	return auth, nil
}

func computeCampusUserStableAuthStatus(tx *gorm.DB, userID uint) (int, error) {
	if hasActive, err := hasActiveCampusAuthTask(tx, uint64(userID), ""); err != nil {
		return campusModel.CampusUserAuthStatusUnverified, err
	} else if hasActive {
		return campusModel.CampusUserAuthStatusProcessing, nil
	}

	var auth campusModel.CampusAuth
	err := tx.Table(auth.TableName()).
		Where("user_id = ? AND review_status IN ?", userID, []string{
			campusModel.CampusAuthReviewStatusApproved,
			campusModel.CampusAuthReviewStatusRejected,
		}).
		Order("reviewed_at DESC, id DESC").
		First(&auth).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return campusModel.CampusUserAuthStatusUnverified, nil
	}
	if err != nil {
		return campusModel.CampusUserAuthStatusUnverified, err
	}

	switch normalizeCampusAuthReviewStatus(auth.ReviewStatus) {
	case campusModel.CampusAuthReviewStatusApproved:
		return campusModel.CampusUserAuthStatusApproved, nil
	case campusModel.CampusAuthReviewStatusRejected:
		return campusModel.CampusUserAuthStatusRejected, nil
	default:
		return campusModel.CampusUserAuthStatusUnverified, nil
	}
}

func resolveRollbackCampusUserAuthStatus(tx *gorm.DB, userID uint, previousStatus int, excludeTaskID string) (int, error) {
	if hasActive, err := hasActiveCampusAuthTask(tx, uint64(userID), excludeTaskID); err != nil {
		return campusModel.CampusUserAuthStatusUnverified, err
	} else if hasActive {
		return campusModel.CampusUserAuthStatusProcessing, nil
	}

	switch previousStatus {
	case campusModel.CampusUserAuthStatusUnverified, campusModel.CampusUserAuthStatusRejected, campusModel.CampusUserAuthStatusApproved:
		return previousStatus, nil
	default:
		return computeCampusUserStableAuthStatus(tx, userID)
	}
}

func mapCallbackToLocalTaskState(req campusReq.AgentReviewCallbackReq) (string, string) {
	switch req.TaskStatus {
	case campusModel.AgentReviewCallbackTaskStatusProcessing:
		return campusModel.AgentReviewTaskStatusProcessing, ""
	case campusModel.AgentReviewCallbackTaskStatusCompleted:
		return campusModel.AgentReviewTaskStatusCompleted, req.FinalAction
	case campusModel.AgentReviewCallbackTaskStatusFailed:
		return campusModel.AgentReviewTaskStatusFailed, ""
	default:
		return campusModel.AgentReviewTaskStatusFailed, ""
	}
}

func isSameTerminalTaskState(task campusModel.AgentReviewTask, nextTaskStatus, nextFinalAction string) bool {
	if !isTerminalTaskStatus(task.TaskStatus) {
		return false
	}
	return task.TaskStatus == nextTaskStatus && strings.TrimSpace(task.FinalAction) == strings.TrimSpace(nextFinalAction)
}

func isTerminalTaskStatus(taskStatus string) bool {
	switch strings.TrimSpace(taskStatus) {
	case campusModel.AgentReviewTaskStatusCompleted, campusModel.AgentReviewTaskStatusFailed, campusModel.AgentReviewTaskStatusSuperseded:
		return true
	default:
		return false
	}
}

func buildCallbackErrorMessage(req campusReq.AgentReviewCallbackReq) string {
	if req.ReviewRemark != "" {
		return req.ReviewRemark
	}
	if len(req.Result) == 0 {
		return "agent callback failed"
	}
	raw, err := json.Marshal(req.Result)
	if err != nil {
		return "agent callback failed"
	}
	return string(raw)
}

func (s *CampusAuthService) dispatchAgentReviewTask(ctx context.Context, task campusModel.AgentReviewTask, payload campusReq.AgentReviewTaskPayload) error {
	cfg := global.GVA_CONFIG.AgentReview
	if !cfg.Enabled {
		return errors.New("agent-review 未启用")
	}
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return errors.New("agent-review.endpoint 未配置")
	}

	if err := global.GVA_DB.WithContext(ctx).
		Table(task.TableName()).
		Where("task_id = ?", task.TaskID).
		UpdateColumn("dispatch_attempt_count", gorm.Expr("dispatch_attempt_count + ?", 1)).Error; err != nil {
		return err
	}

	body := agentReviewDispatchRequest{
		TaskID:      task.TaskID,
		BizType:     task.BizType,
		BizID:       strconv.FormatUint(task.BizID, 10),
		TraceID:     task.TraceID,
		CallbackURL: task.CallbackURL,
		Payload:     payload,
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	timeoutSeconds := cfg.RequestTimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = 10
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodPost, endpoint, bytes.NewReader(rawBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if token := strings.TrimSpace(cfg.RequestToken); token != "" {
		req.Header.Set(agentReviewRequestTokenHeader, token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if len(responseBody) == 0 {
			return fmt.Errorf("agent 接口返回状态码 %d", resp.StatusCode)
		}
		return fmt.Errorf("agent 接口返回状态码 %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}
	return nil
}

func (s *CampusAuthService) markAgentReviewTaskSubmitted(ctx context.Context, taskID string) error {
	result := global.GVA_DB.WithContext(ctx).
		Table(campusModel.AgentReviewTask{}.TableName()).
		Where("task_id = ? AND task_status = ?", taskID, campusModel.AgentReviewTaskStatusPendingDispatch).
		Updates(map[string]interface{}{
			"task_status":    campusModel.AgentReviewTaskStatusSubmitted,
			"dispatch_error": "",
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("任务状态不是 pending_dispatch")
	}
	return nil
}

func (s *CampusAuthService) markAgentReviewTaskDispatchFailed(ctx context.Context, taskID string, dispatchErr string) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var task campusModel.AgentReviewTask
		if err := tx.Table(task.TableName()).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("task_id = ?", taskID).
			First(&task).Error; err != nil {
			return err
		}
		if isTerminalTaskStatus(task.TaskStatus) {
			return nil
		}

		if err := tx.Table(task.TableName()).Where("id = ?", task.ID).Updates(map[string]interface{}{
			"task_status":    campusModel.AgentReviewTaskStatusFailed,
			"dispatch_error": strings.TrimSpace(dispatchErr),
		}).Error; err != nil {
			return err
		}

		authStatus, err := resolveRollbackCampusUserAuthStatus(tx, uint(task.UserID), task.PreviousUserAuthStatus, task.TaskID)
		if err != nil {
			return err
		}
		return updateCampusUserAuthStatus(tx, uint(task.UserID), authStatus)
	})
}
