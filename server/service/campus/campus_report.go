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

type CampusReportService struct{}

func (s *CampusReportService) GetCampusReportInfoList(ctx context.Context, info campusReq.CampusReportSearch) (list []campusModel.CampusReport, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := s.baseCampusReportQuery(ctx)

	if id, ok := parseUintFromText(info.ID); ok {
		db = db.Where("r.id = ?", id)
	}
	if reporterID, ok := parseUintFromText(info.ReporterID); ok {
		db = db.Where("r.reporter_id = ?", reporterID)
	}
	if info.TargetType != nil {
		db = db.Where("r.target_type = ?", *info.TargetType)
	}
	if targetID, ok := parseUintFromText(info.TargetID); ok {
		db = db.Where("r.target_id = ?", targetID)
	}
	if info.Reason != nil {
		db = db.Where("r.reason = ?", *info.Reason)
	}
	if info.Status != nil {
		db = db.Where("r.status = ?", *info.Status)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("r.created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []campusModel.CampusReport
	query := db.Order("r.created_at DESC")
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err = query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	s.fillCampusReportTexts(items)
	return items, total, nil
}

func (s *CampusReportService) GetCampusReport(ctx context.Context, id uint) (report campusModel.CampusReport, err error) {
	err = s.baseCampusReportQuery(ctx).Where("r.id = ?", id).First(&report).Error
	if err == nil {
		report.StatusText = buildReportStatusText(report.Status)
		report.TargetTypeText = buildReportTargetTypeText(report.TargetType)
		report.ReasonText = buildReportReasonText(report.Reason)
	}
	return
}

func (s *CampusReportService) HandleCampusReport(ctx context.Context, operator string, req campusReq.HandleCampusReportReq) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var report campusModel.CampusReport
		if err := tx.Table(report.TableName()).Where("id = ?", req.ID).First(&report).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("举报记录不存在")
			}
			return err
		}

		handlerID, err := resolveCampusOperatorUserID(ctx, operator)
		if err != nil {
			return err
		}
		if handlerID == nil {
			return errors.New("当前后台账号未绑定业务管理员，无法记录处理人")
		}

		return tx.Table(report.TableName()).Where("id = ?", req.ID).Updates(map[string]interface{}{
			"status":        *req.Status,
			"handled_by":    *handlerID,
			"handle_result": strings.TrimSpace(req.HandleResult),
		}).Error
	})
}

func (s *CampusReportService) baseCampusReportQuery(ctx context.Context) *gorm.DB {
	return global.GVA_DB.WithContext(ctx).
		Table("t_report AS r").
		Select("r.*, reporter.nickname AS reporter_nickname, reporter.phone AS reporter_phone, handled.nickname AS handled_by_nickname, handled.phone AS handled_by_phone, p.title AS target_product_title").
		Joins("LEFT JOIN t_user AS reporter ON reporter.id = r.reporter_id").
		Joins("LEFT JOIN t_user AS handled ON handled.id = r.handled_by").
		Joins("LEFT JOIN t_product AS p ON r.target_type = 1 AND p.id = r.target_id")
}

func (s *CampusReportService) fillCampusReportTexts(items []campusModel.CampusReport) {
	for i := range items {
		items[i].StatusText = buildReportStatusText(items[i].Status)
		items[i].TargetTypeText = buildReportTargetTypeText(items[i].TargetType)
		items[i].ReasonText = buildReportReasonText(items[i].Reason)
	}
}
