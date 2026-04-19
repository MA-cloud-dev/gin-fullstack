package campus

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	"gorm.io/gorm"
)

func parseUintFromText(raw string) (uint64, bool) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, false
	}
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func buildStatusText(status int) string {
	switch status {
	case 0:
		return "启用"
	case 1:
		return "禁用"
	default:
		return fmt.Sprintf("状态(%d)", status)
	}
}

func buildProductStatusText(status int) string {
	switch status {
	case 0:
		return "在售"
	case 1:
		return "待审核"
	case 2:
		return "交易中"
	case 3:
		return "已下架"
	case 4:
		return "已售出"
	default:
		return fmt.Sprintf("状态(%d)", status)
	}
}

func buildTradeModeText(mode int) string {
	switch mode {
	case 0:
		return "线下自提"
	case 1:
		return "快递邮寄"
	default:
		return fmt.Sprintf("方式(%d)", mode)
	}
}

func buildUserRoleText(role int) string {
	switch role {
	case 0:
		return "普通用户"
	case 1:
		return "管理员"
	default:
		return fmt.Sprintf("角色(%d)", role)
	}
}

func buildAdminRoleTypeText(roleType int) string {
	switch roleType {
	case 1:
		return "超级管理员"
	case 2:
		return "运营管理员"
	default:
		return fmt.Sprintf("角色(%d)", roleType)
	}
}

func buildAuthStatusText(user campusModel.CampusUser) string {
	if user.Role == 1 {
		return "管理员"
	}
	switch user.AuthStatus {
	case campusModel.CampusUserAuthStatusRejected:
		return "已拒绝"
	case campusModel.CampusUserAuthStatusProcessing:
		return "审核中"
	case campusModel.CampusUserAuthStatusApproved:
		return "已认证"
	case campusModel.CampusUserAuthStatusUnverified:
		return "未认证"
	default:
		return fmt.Sprintf("认证状态(%d)", user.AuthStatus)
	}
}

func buildCampusAuthReviewStatusText(status string) string {
	switch normalizeCampusAuthReviewStatus(status) {
	case campusModel.CampusAuthReviewStatusApproved:
		return "已通过"
	case campusModel.CampusAuthReviewStatusRejected:
		return "已拒绝"
	case campusModel.CampusAuthReviewStatusProcessing:
		return "审核中"
	default:
		return "审核中"
	}
}

func normalizeCampusAuthReviewStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case campusModel.CampusAuthReviewStatusApproved:
		return campusModel.CampusAuthReviewStatusApproved
	case campusModel.CampusAuthReviewStatusRejected:
		return campusModel.CampusAuthReviewStatusRejected
	case campusModel.CampusAuthReviewStatusProcessing:
		return campusModel.CampusAuthReviewStatusProcessing
	default:
		return campusModel.CampusAuthReviewStatusProcessing
	}
}

func buildAnnouncementStatusText(status int) string {
	switch status {
	case 0:
		return "下线"
	case 1:
		return "上线"
	default:
		return fmt.Sprintf("状态(%d)", status)
	}
}

func buildOrderStatusText(status int) string {
	switch status {
	case 1:
		return "待付款"
	case 2:
		return "待确认"
	case 3:
		return "已完成"
	case 4:
		return "已取消"
	case 5:
		return "已关闭"
	default:
		return fmt.Sprintf("状态(%d)", status)
	}
}

func buildOrderCloseByText(closeBy *int) string {
	if closeBy == nil || *closeBy == 0 {
		return "-"
	}
	switch *closeBy {
	case 1:
		return "买家"
	case 2:
		return "卖家"
	case 3:
		return "系统"
	default:
		return fmt.Sprintf("关闭方(%d)", *closeBy)
	}
}

func buildOrderCloseConfirmedText(closeConfirmed *int) string {
	if closeConfirmed == nil || *closeConfirmed == 0 {
		return "-"
	}
	switch *closeConfirmed {
	case 1:
		return "已确认"
	case 2:
		return "未确认"
	default:
		return fmt.Sprintf("确认状态(%d)", *closeConfirmed)
	}
}

func buildReportStatusText(status int) string {
	switch status {
	case 0:
		return "待处理"
	case 1:
		return "已处理"
	default:
		return fmt.Sprintf("状态(%d)", status)
	}
}

func buildReportTargetTypeText(targetType int) string {
	switch targetType {
	case 1:
		return "商品"
	default:
		return fmt.Sprintf("目标(%d)", targetType)
	}
}

func buildReportReasonText(reason int) string {
	switch reason {
	case 1:
		return "虚假信息"
	case 2:
		return "疑似违规"
	case 3:
		return "疑似诈骗"
	default:
		return fmt.Sprintf("原因(%d)", reason)
	}
}

func buildCampusOperationSourceText(source string) string {
	switch strings.TrimSpace(source) {
	case "cli":
		return "CLI工具"
	case "agent":
		return "Agent回调"
	case "web":
		return "B端页面"
	default:
		return source
	}
}

func buildCampusOperationModuleText(module string) string {
	switch strings.TrimSpace(module) {
	case "auth":
		return "校园身份审核"
	case "user":
		return "校园用户"
	case "staff":
		return "B端管理员"
	case "product":
		return "商品"
	case "report":
		return "举报"
	case "category":
		return "分类"
	case "announcement":
		return "公告"
	default:
		return module
	}
}

func buildCampusOperationActionText(action string) string {
	switch strings.TrimSpace(action) {
	case "approve_auth":
		return "通过校园审核"
	case "reject_auth":
		return "拒绝校园审核"
	case "revoke_auth":
		return "撤回校园审核"
	case "agent_review_processing":
		return "Agent审核处理中"
	case "agent_review_approve":
		return "Agent审核通过"
	case "agent_review_reject":
		return "Agent审核拒绝"
	case "agent_review_escalate":
		return "Agent转人工复核"
	case "agent_review_failed":
		return "Agent审核失败"
	case "enable_user":
		return "启用用户"
	case "disable_user":
		return "禁用用户"
	case "enable_staff":
		return "启用B端管理员"
	case "disable_staff":
		return "禁用B端管理员"
	case "set_product_status":
		return "调整商品状态"
	case "handle_report":
		return "处理举报"
	case "create_category":
		return "创建分类"
	case "update_category":
		return "更新分类"
	case "enable_category":
		return "启用分类"
	case "disable_category":
		return "停用分类"
	case "create_announcement":
		return "创建公告"
	case "update_announcement":
		return "更新公告"
	case "publish_announcement":
		return "上线公告"
	case "unpublish_announcement":
		return "下线公告"
	case "delete_announcement":
		return "删除公告"
	default:
		return action
	}
}

func resolveCampusOperatorUserID(ctx context.Context, username string) (*uint, error) {
	operator := strings.TrimSpace(username)
	if operator == "" {
		return nil, nil
	}

	var user campusModel.CampusUser
	err := global.GVA_DB.WithContext(ctx).
		Table(user.TableName()).
		Select("id").
		Where("phone = ? OR username = ?", operator, operator).
		Order("role DESC, id ASC").
		First(&user).Error
	if err == nil {
		id := user.ID
		return &id, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func latestCampusAuthSubQuery(ctx context.Context) *gorm.DB {
	return global.GVA_DB.WithContext(ctx).
		Table("t_campus_auth AS ca").
		Select("ca.*").
		Joins("INNER JOIN (SELECT user_id, MAX(id) AS max_id FROM t_campus_auth GROUP BY user_id) AS latest ON latest.max_id = ca.id")
}
