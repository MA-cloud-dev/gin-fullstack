package campus

import (
	"context"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"gorm.io/gorm"
)

type CampusOrderService struct{}

func (s *CampusOrderService) GetCampusOrderInfoList(ctx context.Context, info campusReq.CampusOrderSearch) (list []campusModel.CampusOrder, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := s.baseCampusOrderQuery(ctx)

	if orderNo := strings.TrimSpace(info.OrderNo); orderNo != "" {
		db = db.Where("o.order_no LIKE ?", "%"+orderNo+"%")
	}
	if buyerID, ok := parseUintFromText(info.BuyerID); ok {
		db = db.Where("o.buyer_id = ?", buyerID)
	}
	if sellerID, ok := parseUintFromText(info.SellerID); ok {
		db = db.Where("o.seller_id = ?", sellerID)
	}
	if productID, ok := parseUintFromText(info.ProductID); ok {
		db = db.Where("o.product_id = ?", productID)
	}
	if info.Status != nil {
		db = db.Where("o.status = ?", *info.Status)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("o.created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []campusModel.CampusOrder
	query := db.Order("o.created_at DESC")
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err = query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	s.fillCampusOrderTexts(items)
	return items, total, nil
}

func (s *CampusOrderService) GetCampusOrder(ctx context.Context, id uint) (order campusModel.CampusOrder, err error) {
	err = s.baseCampusOrderQuery(ctx).Where("o.id = ?", id).First(&order).Error
	if err == nil {
		order.StatusText = buildOrderStatusText(order.Status)
		order.CloseByText = buildOrderCloseByText(order.CloseBy)
		order.CloseConfirmedText = buildOrderCloseConfirmedText(order.CloseConfirmed)
	}
	return
}

func (s *CampusOrderService) baseCampusOrderQuery(ctx context.Context) *gorm.DB {
	return global.GVA_DB.WithContext(ctx).
		Table("t_order AS o").
		Select("o.*, buyer.nickname AS buyer_nickname, buyer.phone AS buyer_phone, seller.nickname AS seller_nickname, seller.phone AS seller_phone").
		Joins("LEFT JOIN t_user AS buyer ON buyer.id = o.buyer_id").
		Joins("LEFT JOIN t_user AS seller ON seller.id = o.seller_id")
}

func (s *CampusOrderService) fillCampusOrderTexts(items []campusModel.CampusOrder) {
	for i := range items {
		items[i].StatusText = buildOrderStatusText(items[i].Status)
		items[i].CloseByText = buildOrderCloseByText(items[i].CloseBy)
		items[i].CloseConfirmedText = buildOrderCloseConfirmedText(items[i].CloseConfirmed)
	}
}
