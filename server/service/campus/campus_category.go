package campus

import (
	"context"
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
)

type CampusCategoryService struct{}

func (s *CampusCategoryService) GetCampusCategoryTree(ctx context.Context, info campusReq.CampusCategorySearch) (list []campusModel.CampusCategory, err error) {
	db := global.GVA_DB.WithContext(ctx).
		Table("t_category AS c").
		Select("c.*, pc.name AS parent_name").
		Joins("LEFT JOIN t_category AS pc ON pc.id = c.parent_id")

	if name := strings.TrimSpace(info.Name); name != "" {
		db = db.Where("c.name LIKE ?", "%"+name+"%")
	}
	if info.Status != nil {
		db = db.Where("c.status = ?", *info.Status)
	}

	var items []campusModel.CampusCategory
	if err = db.Order("COALESCE(c.parent_id, 0) ASC, c.sort_order ASC, c.id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	for i := range items {
		items[i].StatusText = buildStatusText(items[i].Status)
	}
	return buildCategoryTree(items), nil
}

func (s *CampusCategoryService) GetCampusCategory(ctx context.Context, id uint) (category campusModel.CampusCategory, err error) {
	err = global.GVA_DB.WithContext(ctx).
		Table("t_category AS c").
		Select("c.*, pc.name AS parent_name").
		Joins("LEFT JOIN t_category AS pc ON pc.id = c.parent_id").
		Where("c.id = ?", id).
		First(&category).Error
	if err == nil {
		category.StatusText = buildStatusText(category.Status)
	}
	return
}

func (s *CampusCategoryService) CreateCampusCategory(ctx context.Context, req campusReq.CreateCampusCategoryReq) error {
	if req.ParentID != nil && *req.ParentID > 0 {
		var count int64
		if err := global.GVA_DB.WithContext(ctx).Table("t_category").Where("id = ?", *req.ParentID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("父级分类不存在")
		}
	}

	category := campusModel.CampusCategory{
		Name:      strings.TrimSpace(req.Name),
		ParentID:  sanitizeParentID(req.ParentID),
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}
	if icon := strings.TrimSpace(req.Icon); icon != "" {
		category.Icon = &icon
	}
	return global.GVA_DB.WithContext(ctx).Table(category.TableName()).Create(&category).Error
}

func (s *CampusCategoryService) UpdateCampusCategory(ctx context.Context, req campusReq.UpdateCampusCategoryReq) error {
	if req.ParentID != nil {
		if *req.ParentID == req.ID {
			return errors.New("父级分类不能选择自己")
		}
		if *req.ParentID > 0 {
			var count int64
			if err := global.GVA_DB.WithContext(ctx).Table("t_category").Where("id = ?", *req.ParentID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return errors.New("父级分类不存在")
			}
		}
	}

	updates := map[string]interface{}{
		"name":       strings.TrimSpace(req.Name),
		"parent_id":  sanitizeParentID(req.ParentID),
		"sort_order": req.SortOrder,
		"icon":       nil,
		"status":     req.Status,
	}
	if icon := strings.TrimSpace(req.Icon); icon != "" {
		updates["icon"] = icon
	}
	return global.GVA_DB.WithContext(ctx).Table("t_category").Where("id = ?", req.ID).Updates(updates).Error
}

func (s *CampusCategoryService) UpdateCampusCategoryStatus(ctx context.Context, req campusReq.UpdateCampusCategoryStatusReq) error {
	return global.GVA_DB.WithContext(ctx).Table("t_category").Where("id = ?", req.ID).Update("status", *req.Status).Error
}

func sanitizeParentID(parentID *uint) *uint {
	if parentID == nil || *parentID == 0 {
		return nil
	}
	return parentID
}

func buildCategoryTree(items []campusModel.CampusCategory) []campusModel.CampusCategory {
	childrenMap := make(map[uint][]campusModel.CampusCategory)
	roots := make([]campusModel.CampusCategory, 0)
	itemIDs := make(map[uint]struct{}, len(items))

	for _, item := range items {
		itemIDs[item.ID] = struct{}{}
	}

	for i := range items {
		items[i].Children = nil
		if items[i].ParentID != nil {
			if _, ok := itemIDs[*items[i].ParentID]; !ok {
				roots = append(roots, items[i])
				continue
			}
			childrenMap[*items[i].ParentID] = append(childrenMap[*items[i].ParentID], items[i])
			continue
		}
		roots = append(roots, items[i])
	}

	var attachChildren func(nodes []campusModel.CampusCategory) []campusModel.CampusCategory
	attachChildren = func(nodes []campusModel.CampusCategory) []campusModel.CampusCategory {
		for i := range nodes {
			if children, ok := childrenMap[nodes[i].ID]; ok {
				nodes[i].Children = attachChildren(children)
			}
		}
		return nodes
	}

	return attachChildren(roots)
}
