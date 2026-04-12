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
	if user.ReviewedAt != nil {
		return "已认证"
	}
	if user.StudentID != nil && strings.TrimSpace(*user.StudentID) != "" {
		return "待审核"
	}
	switch user.AuthStatus {
	case 2:
		return "待审核"
	case 3:
		return "已认证"
	default:
		return "未认证"
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
