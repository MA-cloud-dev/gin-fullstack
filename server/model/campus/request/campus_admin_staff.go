package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CampusAdminStaffSearch struct {
	ID             string      `json:"id" form:"id"`
	Username       string      `json:"username" form:"username"`
	DisplayName    string      `json:"displayName" form:"displayName"`
	RoleType       *int        `json:"roleType" form:"roleType"`
	Status         *int        `json:"status" form:"status"`
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	commonReq.PageInfo
}

type UpdateCampusAdminStaffStatusReq struct {
	ID          uint   `json:"id" binding:"required"`
	Status      *int   `json:"status" binding:"required"`
	AuditReason string `json:"auditReason" binding:"required"`
}
