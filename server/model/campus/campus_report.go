package campus

import "time"

type CampusReport struct {
	ID                 uint      `json:"id" gorm:"column:id;primaryKey"`
	ReporterID         uint      `json:"reporterId" gorm:"column:reporter_id"`
	TargetType         int       `json:"targetType" gorm:"column:target_type"`
	TargetID           uint      `json:"targetId" gorm:"column:target_id"`
	Reason             int       `json:"reason" gorm:"column:reason"`
	Description        *string   `json:"description" gorm:"column:description"`
	Status             int       `json:"status" gorm:"column:status"`
	HandledBy          *uint     `json:"handledBy" gorm:"column:handled_by"`
	HandleResult       *string   `json:"handleResult" gorm:"column:handle_result"`
	CreatedAt          time.Time `json:"createdAt" gorm:"column:created_at"`
	ReporterNickname   *string   `json:"reporterNickname" gorm:"column:reporter_nickname;->"`
	ReporterPhone      *string   `json:"reporterPhone" gorm:"column:reporter_phone;->"`
	HandledByNickname  *string   `json:"handledByNickname" gorm:"column:handled_by_nickname;->"`
	HandledByPhone     *string   `json:"handledByPhone" gorm:"column:handled_by_phone;->"`
	TargetProductTitle *string   `json:"targetProductTitle" gorm:"column:target_product_title;->"`
	StatusText         string    `json:"statusText" gorm:"-"`
	TargetTypeText     string    `json:"targetTypeText" gorm:"-"`
	ReasonText         string    `json:"reasonText" gorm:"-"`
}

func (CampusReport) TableName() string {
	return "t_report"
}
