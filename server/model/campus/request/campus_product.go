package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CampusProductSearch struct {
	ID             string      `json:"id" form:"id"`
	Title          string      `json:"title" form:"title"`
	UserID         string      `json:"userId" form:"userId"`
	CategoryID     string      `json:"categoryId" form:"categoryId"`
	Status         *int        `json:"status" form:"status"`
	TradeMode      *int        `json:"tradeMode" form:"tradeMode"`
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	commonReq.PageInfo
}

type UpdateCampusProductStatusReq struct {
	ID     uint `json:"id" form:"id" binding:"required"`
	Status *int `json:"status" form:"status" binding:"required"`
}
