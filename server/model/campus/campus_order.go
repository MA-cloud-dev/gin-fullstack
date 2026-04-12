package campus

import "time"

type CampusOrder struct {
	ID                 uint       `json:"id" gorm:"column:id;primaryKey"`
	OrderNo            string     `json:"orderNo" gorm:"column:order_no"`
	BuyerID            uint       `json:"buyerId" gorm:"column:buyer_id"`
	SellerID           uint       `json:"sellerId" gorm:"column:seller_id"`
	ProductID          uint       `json:"productId" gorm:"column:product_id"`
	ProductTitle       string     `json:"productTitle" gorm:"column:product_title"`
	ProductImage       *string    `json:"productImage" gorm:"column:product_image"`
	Price              float64    `json:"price" gorm:"column:price"`
	Status             int        `json:"status" gorm:"column:status"`
	CloseReason        *string    `json:"closeReason" gorm:"column:close_reason"`
	CloseBy            *int       `json:"closeBy" gorm:"column:close_by"`
	CloseConfirmed     *int       `json:"closeConfirmed" gorm:"column:close_confirmed"`
	Remark             *string    `json:"remark" gorm:"column:remark"`
	ConfirmedAt        *time.Time `json:"confirmedAt" gorm:"column:confirmed_at"`
	CompletedAt        *time.Time `json:"completedAt" gorm:"column:completed_at"`
	CancelledAt        *time.Time `json:"cancelledAt" gorm:"column:cancelled_at"`
	CreatedAt          time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt          time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	BuyerNickname      *string    `json:"buyerNickname" gorm:"column:buyer_nickname;->"`
	BuyerPhone         *string    `json:"buyerPhone" gorm:"column:buyer_phone;->"`
	SellerNickname     *string    `json:"sellerNickname" gorm:"column:seller_nickname;->"`
	SellerPhone        *string    `json:"sellerPhone" gorm:"column:seller_phone;->"`
	StatusText         string     `json:"statusText" gorm:"-"`
	CloseByText        string     `json:"closeByText" gorm:"-"`
	CloseConfirmedText string     `json:"closeConfirmedText" gorm:"-"`
}

func (CampusOrder) TableName() string {
	return "t_order"
}
