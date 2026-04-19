package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CampusAnnouncementSearch struct {
	Title            string      `json:"title" form:"title"`
	PublisherID      string      `json:"publisherId" form:"publisherId"`
	PublisherKeyword string      `json:"publisherKeyword" form:"publisherKeyword"`
	Status           *int        `json:"status" form:"status"`
	CreatedAtRange   []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	commonReq.PageInfo
}

type CreateCampusAnnouncementReq struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	PublisherID uint   `json:"publisherId" binding:"required"`
	Status      *int   `json:"status" binding:"required"`
	AuditReason string `json:"auditReason"`
}

type UpdateCampusAnnouncementReq struct {
	ID          uint   `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	PublisherID uint   `json:"publisherId" binding:"required"`
	Status      *int   `json:"status" binding:"required"`
	AuditReason string `json:"auditReason"`
}

type DeleteCampusAnnouncementReq struct {
	ID          uint   `json:"id" binding:"required"`
	AuditReason string `json:"auditReason"`
}

type UpdateCampusAnnouncementStatusReq struct {
	ID          uint   `json:"id" binding:"required"`
	Status      *int   `json:"status" binding:"required"`
	AuditReason string `json:"auditReason"`
}
