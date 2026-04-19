package campus

import (
	"context"
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"gorm.io/gorm"
)

type CampusProductService struct{}

func (s *CampusProductService) GetCampusProductInfoList(ctx context.Context, info campusReq.CampusProductSearch) (list []campusModel.CampusProduct, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.WithContext(ctx).
		Table("t_product AS p").
		Select("p.*, c.name AS category_name, u.nickname AS publisher_nickname, u.phone AS publisher_phone").
		Joins("LEFT JOIN t_category AS c ON c.id = p.category_id").
		Joins("LEFT JOIN t_user AS u ON u.id = p.user_id")

	if id, ok := parseUintFromText(info.ID); ok {
		db = db.Where("p.id = ?", id)
	}
	if title := strings.TrimSpace(info.Title); title != "" {
		db = db.Where("p.title LIKE ?", "%"+title+"%")
	}
	if userID, ok := parseUintFromText(info.UserID); ok {
		db = db.Where("p.user_id = ?", userID)
	}
	if categoryID, ok := parseUintFromText(info.CategoryID); ok {
		db = db.Where("p.category_id = ?", categoryID)
	}
	if info.Status != nil {
		db = db.Where("p.status = ?", *info.Status)
	}
	if info.TradeMode != nil {
		db = db.Where("p.trade_mode = ?", *info.TradeMode)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("p.created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []campusModel.CampusProduct
	query := db.Order("p.created_at DESC")
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err = query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	s.fillCampusProductTexts(items)
	if err = s.fillCampusProductCovers(ctx, items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *CampusProductService) GetCampusProduct(ctx context.Context, id uint) (product campusModel.CampusProduct, err error) {
	err = global.GVA_DB.WithContext(ctx).
		Table("t_product AS p").
		Select("p.*, c.name AS category_name, u.nickname AS publisher_nickname, u.phone AS publisher_phone").
		Joins("LEFT JOIN t_category AS c ON c.id = p.category_id").
		Joins("LEFT JOIN t_user AS u ON u.id = p.user_id").
		Where("p.id = ?", id).
		First(&product).Error
	if err != nil {
		return product, err
	}

	product.StatusText = buildProductStatusText(product.Status)
	product.TradeModeText = buildTradeModeText(product.TradeMode)

	var imageModel campusModel.CampusProductImage
	var images []campusModel.CampusProductImage
	if err = global.GVA_DB.WithContext(ctx).
		Table(imageModel.TableName()).
		Where("product_id = ?", product.ID).
		Order("sort_order ASC, id ASC").
		Find(&images).Error; err != nil {
		return product, err
	}
	product.Images = images
	if len(images) > 0 {
		product.CoverURL = &images[0].ImageURL
	}
	return product, nil
}

func (s *CampusProductService) UpdateCampusProductStatus(ctx context.Context, req campusReq.UpdateCampusProductStatusReq, auditMeta campusReq.CampusAuditMeta) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product campusModel.CampusProduct
		if err := tx.Table(product.TableName()).Where("id = ?", req.ID).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("商品不存在")
			}
			return err
		}
		if err := tx.Table(product.TableName()).Where("id = ?", req.ID).Update("status", *req.Status).Error; err != nil {
			return err
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "product",
			Action:      "set_product_status",
			TargetID:    product.ID,
			TargetLabel: joinCampusAuditLabel(product.Title, buildCampusAuditIDLabel("商品", product.ID)),
			Reason:      req.AuditReason,
			Result:      "商品状态已更新为" + buildProductStatusText(*req.Status),
		})
	})
}

func (s *CampusProductService) fillCampusProductTexts(items []campusModel.CampusProduct) {
	for i := range items {
		items[i].StatusText = buildProductStatusText(items[i].Status)
		items[i].TradeModeText = buildTradeModeText(items[i].TradeMode)
	}
}

func (s *CampusProductService) fillCampusProductCovers(ctx context.Context, items []campusModel.CampusProduct) error {
	if len(items) == 0 {
		return nil
	}

	productIDs := make([]uint, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ID)
	}

	var imageModel campusModel.CampusProductImage
	var images []campusModel.CampusProductImage
	if err := global.GVA_DB.WithContext(ctx).
		Table(imageModel.TableName()).
		Where("product_id IN ?", productIDs).
		Order("product_id ASC, sort_order ASC, id ASC").
		Find(&images).Error; err != nil {
		return err
	}

	coverMap := make(map[uint]string, len(images))
	for _, image := range images {
		if _, ok := coverMap[image.ProductID]; !ok {
			coverMap[image.ProductID] = image.ImageURL
		}
	}
	for i := range items {
		if cover, ok := coverMap[items[i].ID]; ok {
			items[i].CoverURL = &cover
		}
	}
	return nil
}
