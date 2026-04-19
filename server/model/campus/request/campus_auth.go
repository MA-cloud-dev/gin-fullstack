package request

import (
	"errors"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"strings"
	"time"
)

type CampusAuthSearch struct {
	StudentID       string      `json:"studentId" form:"studentId"`
	RealName        string      `json:"realName" form:"realName"`
	College         string      `json:"college" form:"college"`
	ReviewStatus    string      `json:"reviewStatus" form:"reviewStatus"`
	CreatedAtRange  []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	ReviewedAtRange []time.Time `json:"reviewedAtRange" form:"reviewedAtRange[]"`
	commonReq.PageInfo
}

type ReviewCampusAuthReq struct {
	ID           uint   `json:"id" form:"id" binding:"required"`
	ReviewRemark string `json:"reviewRemark" form:"reviewRemark"`
	AuditReason  string `json:"auditReason" form:"auditReason" binding:"required"`
}

type SubmitCampusAuthReq struct {
	UserID    uint64 `json:"userId" binding:"required"`
	StudentID string `json:"studentId" binding:"required"`
	RealName  string `json:"realName" binding:"required"`
	College   string `json:"college" binding:"required"`
}

type AgentReviewTaskPayload struct {
	UserID    string `json:"userId"`
	StudentID string `json:"studentId"`
	RealName  string `json:"realName"`
	College   string `json:"college"`
}

type AgentReviewCallbackReq struct {
	TaskID       string         `json:"taskId" binding:"required"`
	BizType      string         `json:"bizType" binding:"required"`
	BizID        string         `json:"bizId" binding:"required"`
	TraceID      string         `json:"traceId" binding:"required"`
	TaskStatus   string         `json:"taskStatus" binding:"required"`
	FinalAction  string         `json:"finalAction"`
	ReviewRemark string         `json:"reviewRemark"`
	Result       map[string]any `json:"result"`
}

func (r *AgentReviewCallbackReq) Normalize() {
	r.TaskID = strings.TrimSpace(r.TaskID)
	r.BizType = strings.TrimSpace(r.BizType)
	r.BizID = strings.TrimSpace(r.BizID)
	r.TraceID = strings.TrimSpace(r.TraceID)
	r.TaskStatus = strings.ToLower(strings.TrimSpace(r.TaskStatus))
	r.FinalAction = strings.ToLower(strings.TrimSpace(r.FinalAction))
	r.ReviewRemark = strings.TrimSpace(r.ReviewRemark)
}

func (r AgentReviewCallbackReq) Validate() error {
	switch r.TaskStatus {
	case campusModel.AgentReviewCallbackTaskStatusProcessing:
		if r.FinalAction != "" {
			return errors.New("processing 阶段不允许携带 finalAction")
		}
	case campusModel.AgentReviewCallbackTaskStatusCompleted:
		switch r.FinalAction {
		case campusModel.AgentReviewFinalActionApprove, campusModel.AgentReviewFinalActionReject, campusModel.AgentReviewFinalActionEscalate:
		default:
			return errors.New("completed 阶段必须携带有效的 finalAction")
		}
	case campusModel.AgentReviewCallbackTaskStatusFailed:
		if r.FinalAction != "" {
			return errors.New("failed 阶段不允许携带 finalAction")
		}
	default:
		return errors.New("taskStatus 不合法")
	}
	return nil
}
