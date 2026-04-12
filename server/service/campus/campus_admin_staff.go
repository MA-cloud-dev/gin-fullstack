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

type CampusAdminStaffService struct{}

func (s *CampusAdminStaffService) GetCampusAdminStaffInfoList(ctx context.Context, info campusReq.CampusAdminStaffSearch) (list []campusModel.CampusAdminStaff, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.WithContext(ctx).Table("t_admin_staff")

	if id, ok := parseUintFromText(info.ID); ok {
		db = db.Where("id = ?", id)
	}
	if username := strings.TrimSpace(info.Username); username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if displayName := strings.TrimSpace(info.DisplayName); displayName != "" {
		db = db.Where("display_name LIKE ?", "%"+displayName+"%")
	}
	if info.RoleType != nil {
		db = db.Where("role_type = ?", *info.RoleType)
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []campusModel.CampusAdminStaff
	query := db.Order("created_at DESC")
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err = query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].StatusText = buildStatusText(items[i].Status)
		items[i].RoleTypeText = buildAdminRoleTypeText(items[i].RoleType)
	}
	return items, total, nil
}

func (s *CampusAdminStaffService) GetCampusAdminStaff(ctx context.Context, id uint) (admin campusModel.CampusAdminStaff, err error) {
	err = global.GVA_DB.WithContext(ctx).Table(admin.TableName()).Where("id = ?", id).First(&admin).Error
	if err == nil {
		admin.StatusText = buildStatusText(admin.Status)
		admin.RoleTypeText = buildAdminRoleTypeText(admin.RoleType)
	}
	return
}

func (s *CampusAdminStaffService) UpdateCampusAdminStaffStatus(ctx context.Context, req campusReq.UpdateCampusAdminStaffStatusReq) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var admin campusModel.CampusAdminStaff
		if err := tx.Table(admin.TableName()).Where("id = ?", req.ID).First(&admin).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("B端管理员不存在")
			}
			return err
		}
		return tx.Table(admin.TableName()).Where("id = ?", req.ID).Update("status", *req.Status).Error
	})
}
