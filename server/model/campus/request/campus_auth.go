package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CampusAuthSearch struct {
	StudentID       string      `json:"studentId" form:"studentId"`
	RealName        string      `json:"realName" form:"realName"`
	College         string      `json:"college" form:"college"`
	Reviewed        *bool       `json:"reviewed" form:"reviewed"`
	CreatedAtRange  []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	ReviewedAtRange []time.Time `json:"reviewedAtRange" form:"reviewedAtRange[]"`
	commonReq.PageInfo
}

type ReviewCampusAuthReq struct {
	ID           uint   `json:"id" form:"id" binding:"required"`
	ReviewRemark string `json:"reviewRemark" form:"reviewRemark"`
}
