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

type CampusAnnouncementService struct{}

func (s *CampusAnnouncementService) GetCampusAnnouncementInfoList(ctx context.Context, info campusReq.CampusAnnouncementSearch) (list []campusModel.CampusAnnouncement, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.WithContext(ctx).
		Table("t_announcement AS a").
		Select("a.*, u.nickname AS publisher_nickname, u.phone AS publisher_phone").
		Joins("LEFT JOIN t_user AS u ON u.id = a.publisher_id")

	if title := strings.TrimSpace(info.Title); title != "" {
		db = db.Where("a.title LIKE ?", "%"+title+"%")
	}
	if publisherID, ok := parseUintFromText(info.PublisherID); ok {
		db = db.Where("a.publisher_id = ?", publisherID)
	}
	if keyword := strings.TrimSpace(info.PublisherKeyword); keyword != "" {
		db = db.Where("(u.nickname LIKE ? OR u.phone LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if info.Status != nil {
		db = db.Where("a.status = ?", *info.Status)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("a.created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []campusModel.CampusAnnouncement
	query := db.Order("a.created_at DESC")
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err = query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].StatusText = buildAnnouncementStatusText(items[i].Status)
	}
	return items, total, nil
}

func (s *CampusAnnouncementService) GetCampusAnnouncement(ctx context.Context, id uint) (announcement campusModel.CampusAnnouncement, err error) {
	err = global.GVA_DB.WithContext(ctx).
		Table("t_announcement AS a").
		Select("a.*, u.nickname AS publisher_nickname, u.phone AS publisher_phone").
		Joins("LEFT JOIN t_user AS u ON u.id = a.publisher_id").
		Where("a.id = ?", id).
		First(&announcement).Error
	if err == nil {
		announcement.StatusText = buildAnnouncementStatusText(announcement.Status)
	}
	return
}

func (s *CampusAnnouncementService) CreateCampusAnnouncement(ctx context.Context, req campusReq.CreateCampusAnnouncementReq, auditMeta campusReq.CampusAuditMeta) error {
	announcement := campusModel.CampusAnnouncement{
		Title:       strings.TrimSpace(req.Title),
		Content:     req.Content,
		PublisherID: req.PublisherID,
		Status:      *req.Status,
	}
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(announcement.TableName()).Create(&announcement).Error; err != nil {
			return err
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "announcement",
			Action:      "create_announcement",
			TargetID:    announcement.ID,
			TargetLabel: announcement.Title,
			Reason:      req.AuditReason,
			Result:      "公告创建成功",
		})
	})
}

func (s *CampusAnnouncementService) UpdateCampusAnnouncement(ctx context.Context, req campusReq.UpdateCampusAnnouncementReq, auditMeta campusReq.CampusAuditMeta) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var announcement campusModel.CampusAnnouncement
		if err := tx.Table(announcement.TableName()).Where("id = ?", req.ID).First(&announcement).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("公告不存在")
			}
			return err
		}
		updates := map[string]interface{}{
			"title":        strings.TrimSpace(req.Title),
			"content":      req.Content,
			"publisher_id": req.PublisherID,
			"status":       *req.Status,
		}
		if err := tx.Table(announcement.TableName()).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "announcement",
			Action:      "update_announcement",
			TargetID:    announcement.ID,
			TargetLabel: strings.TrimSpace(req.Title),
			Reason:      req.AuditReason,
			Result:      "公告更新成功",
		})
	})
}

func (s *CampusAnnouncementService) DeleteCampusAnnouncement(ctx context.Context, req campusReq.DeleteCampusAnnouncementReq, auditMeta campusReq.CampusAuditMeta) error {
	var announcement campusModel.CampusAnnouncement
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(announcement.TableName()).Where("id = ?", req.ID).First(&announcement).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("公告不存在")
			}
			return err
		}
		if err := tx.Table(announcement.TableName()).Delete(&announcement).Error; err != nil {
			return err
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "announcement",
			Action:      "delete_announcement",
			TargetID:    announcement.ID,
			TargetLabel: announcement.Title,
			Reason:      req.AuditReason,
			Result:      "公告删除成功",
		})
	})
}

func (s *CampusAnnouncementService) UpdateCampusAnnouncementStatus(ctx context.Context, req campusReq.UpdateCampusAnnouncementStatusReq, auditMeta campusReq.CampusAuditMeta) error {
	var announcement campusModel.CampusAnnouncement
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(announcement.TableName()).Where("id = ?", req.ID).First(&announcement).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("公告不存在")
			}
			return err
		}
		if err := tx.Table(announcement.TableName()).Where("id = ?", req.ID).Update("status", *req.Status).Error; err != nil {
			return err
		}
		action := "unpublish_announcement"
		if *req.Status == 1 {
			action = "publish_announcement"
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "announcement",
			Action:      action,
			TargetID:    announcement.ID,
			TargetLabel: announcement.Title,
			Reason:      req.AuditReason,
			Result:      "公告状态已更新为" + buildAnnouncementStatusText(*req.Status),
		})
	})
}
