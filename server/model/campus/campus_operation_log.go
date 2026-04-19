package campus

import "time"

type CampusOperationLog struct {
	ID                 uint      `json:"id" gorm:"column:id;primaryKey"`
	OperatorSysUserID  uint      `json:"operatorSysUserId" gorm:"column:operator_sys_user_id"`
	OperatorUsername   string    `json:"operatorUsername" gorm:"column:operator_username"`
	OperatorNickname   string    `json:"operatorNickname" gorm:"column:operator_nickname"`
	OperatorSource     string    `json:"operatorSource" gorm:"column:operator_source"`
	Module             string    `json:"module" gorm:"column:module"`
	Action             string    `json:"action" gorm:"column:action"`
	TargetID           uint      `json:"targetId" gorm:"column:target_id"`
	TargetLabel        string    `json:"targetLabel" gorm:"column:target_label"`
	Reason             string    `json:"reason" gorm:"column:reason"`
	Result             string    `json:"result" gorm:"column:result"`
	RequestPath        string    `json:"requestPath" gorm:"column:request_path"`
	RequestMethod      string    `json:"requestMethod" gorm:"column:request_method"`
	OperatorIP         string    `json:"operatorIp" gorm:"column:operator_ip"`
	CreatedAt          time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt          time.Time `json:"updatedAt" gorm:"column:updated_at"`
	OperatorSourceText string    `json:"operatorSourceText" gorm:"-"`
	ModuleText         string    `json:"moduleText" gorm:"-"`
	ActionText         string    `json:"actionText" gorm:"-"`
}

func (CampusOperationLog) TableName() string {
	return "t_campus_operation_log"
}
