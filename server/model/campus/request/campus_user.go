package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CampusUserSearch struct {
	ID             string      `json:"id" form:"id"`
	Phone          string      `json:"phone" form:"phone"`
	Nickname       string      `json:"nickname" form:"nickname"`
	Role           *int        `json:"role" form:"role"`
	Status         *int        `json:"status" form:"status"`
	AuthStatus     *int        `json:"authStatus" form:"authStatus"`
	StudentID      string      `json:"studentId" form:"studentId"`
	RealName       string      `json:"realName" form:"realName"`
	College        string      `json:"college" form:"college"`
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	commonReq.PageInfo
}

type UpdateCampusUserStatusReq struct {
	ID     uint `json:"id" binding:"required"`
	Status *int `json:"status" binding:"required"`
}
