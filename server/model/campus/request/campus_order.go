package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CampusOrderSearch struct {
	OrderNo        string      `json:"orderNo" form:"orderNo"`
	BuyerID        string      `json:"buyerId" form:"buyerId"`
	SellerID       string      `json:"sellerId" form:"sellerId"`
	ProductID      string      `json:"productId" form:"productId"`
	Status         *int        `json:"status" form:"status"`
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	commonReq.PageInfo
}
