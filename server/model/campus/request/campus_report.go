package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CampusReportSearch struct {
	ID             string      `json:"id" form:"id"`
	ReporterID     string      `json:"reporterId" form:"reporterId"`
	TargetType     *int        `json:"targetType" form:"targetType"`
	TargetID       string      `json:"targetId" form:"targetId"`
	Reason         *int        `json:"reason" form:"reason"`
	Status         *int        `json:"status" form:"status"`
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	commonReq.PageInfo
}

type HandleCampusReportReq struct {
	ID           uint   `json:"id" binding:"required"`
	Status       *int   `json:"status" binding:"required"`
	HandleResult string `json:"handleResult" binding:"required"`
}
