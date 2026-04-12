package campus

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CampusAuthService struct{}

func (s *CampusAuthService) GetCampusAuthInfoList(ctx context.Context, info campusReq.CampusAuthSearch) (list []campus.CampusAuth, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.WithContext(ctx).
		Table("t_campus_auth AS ca").
		Select("ca.*, su.nick_name AS reviewed_by_name").
		Joins("LEFT JOIN sys_users AS su ON su.id = ca.reviewed_by")

	if info.StudentID != "" {
		db = db.Where("ca.student_id LIKE ?", "%"+strings.TrimSpace(info.StudentID)+"%")
	}
	if info.RealName != "" {
		db = db.Where("ca.real_name LIKE ?", "%"+strings.TrimSpace(info.RealName)+"%")
	}
	if info.College != "" {
		db = db.Where("ca.college LIKE ?", "%"+strings.TrimSpace(info.College)+"%")
	}
	if info.Reviewed != nil {
		if *info.Reviewed {
			db = db.Where("ca.reviewed_at IS NOT NULL")
		} else {
			db = db.Where("ca.reviewed_at IS NULL")
		}
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("ca.created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if len(info.ReviewedAtRange) == 2 {
		db = db.Where("ca.reviewed_at BETWEEN ? AND ?", info.ReviewedAtRange[0], info.ReviewedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("ca.created_at DESC")
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	var items []campus.CampusAuth
	err = db.Find(&items).Error
	return items, total, err
}

func (s *CampusAuthService) GetCampusAuth(ctx context.Context, id uint) (auth campus.CampusAuth, err error) {
	err = global.GVA_DB.WithContext(ctx).
		Table("t_campus_auth AS ca").
		Select("ca.*, su.nick_name AS reviewed_by_name").
		Joins("LEFT JOIN sys_users AS su ON su.id = ca.reviewed_by").
		Where("ca.id = ?", id).
		First(&auth).Error
	return
}

func (s *CampusAuthService) ReviewCampusAuth(ctx context.Context, req campusReq.ReviewCampusAuthReq, reviewerID uint) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var auth campus.CampusAuth
		if err := tx.Table(auth.TableName()).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", req.ID).
			First(&auth).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("审核记录不存在")
			}
			return err
		}

		if auth.ReviewedAt != nil {
			return errors.New("该记录已审核，请勿重复操作")
		}

		updates := map[string]interface{}{
			"review_remark": strings.TrimSpace(req.ReviewRemark),
			"reviewed_at":   time.Now(),
			"reviewed_by":   reviewerID,
		}
		return tx.Table(auth.TableName()).Where("id = ?", req.ID).Updates(updates).Error
	})
}

func (s *CampusAuthService) RevokeCampusAuth(ctx context.Context, id uint) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var auth campus.CampusAuth
		if err := tx.Table(auth.TableName()).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", id).
			First(&auth).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("审核记录不存在")
			}
			return err
		}

		if auth.ReviewedAt == nil {
			return errors.New("该记录尚未审核，无需撤回")
		}

		updates := map[string]interface{}{
			"review_remark": nil,
			"reviewed_at":   nil,
			"reviewed_by":   nil,
		}
		return tx.Table(auth.TableName()).Where("id = ?", id).Updates(updates).Error
	})
}
