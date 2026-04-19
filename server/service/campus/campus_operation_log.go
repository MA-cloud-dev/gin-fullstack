package campus

import (
	"context"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"gorm.io/gorm"
)

type CampusOperationLogService struct{}

type campusOperationAuditInput struct {
	Meta        campusReq.CampusAuditMeta
	Module      string
	Action      string
	TargetID    uint
	TargetLabel string
	Reason      string
	Result      string
}

func (s *CampusOperationLogService) GetCampusOperationLogInfoList(ctx context.Context, info campusReq.CampusOperationLogSearch) (list []campusModel.CampusOperationLog, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.WithContext(ctx).Table("t_campus_operation_log")

	if keyword := strings.TrimSpace(info.OperatorKeyword); keyword != "" {
		db = db.Where("(operator_username LIKE ? OR operator_nickname LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if source := normalizeCampusOperationSource(info.OperatorSource); source != "" {
		db = db.Where("operator_source = ?", source)
	}
	if module := strings.TrimSpace(info.Module); module != "" {
		db = db.Where("module = ?", module)
	}
	if action := strings.TrimSpace(info.Action); action != "" {
		db = db.Where("action = ?", action)
	}
	if targetID, ok := parseUintFromText(info.TargetID); ok {
		db = db.Where("target_id = ?", targetID)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []campusModel.CampusOperationLog
	query := db.Order("created_at DESC, id DESC")
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err = query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	fillCampusOperationLogListTexts(items)
	return items, total, nil
}

func (s *CampusOperationLogService) GetCampusOperationLog(ctx context.Context, id uint) (item campusModel.CampusOperationLog, err error) {
	err = global.GVA_DB.WithContext(ctx).Table(item.TableName()).Where("id = ?", id).First(&item).Error
	if err != nil {
		return item, err
	}
	fillCampusOperationLogTexts(&item)
	return item, nil
}

func createCampusOperationLogWithTx(tx *gorm.DB, input campusOperationAuditInput) error {
	record := campusModel.CampusOperationLog{
		OperatorSysUserID: input.Meta.OperatorSysUserID,
		OperatorUsername:  strings.TrimSpace(input.Meta.OperatorUsername),
		OperatorNickname:  strings.TrimSpace(input.Meta.OperatorNickname),
		OperatorSource:    normalizeCampusOperationSource(input.Meta.OperatorSource),
		Module:            strings.TrimSpace(input.Module),
		Action:            strings.TrimSpace(input.Action),
		TargetID:          input.TargetID,
		TargetLabel:       strings.TrimSpace(input.TargetLabel),
		Reason:            strings.TrimSpace(input.Reason),
		Result:            strings.TrimSpace(input.Result),
		RequestPath:       strings.TrimSpace(input.Meta.RequestPath),
		RequestMethod:     strings.ToUpper(strings.TrimSpace(input.Meta.RequestMethod)),
		OperatorIP:        strings.TrimSpace(input.Meta.OperatorIP),
	}
	return tx.Table(record.TableName()).Create(&record).Error
}

func fillCampusOperationLogListTexts(items []campusModel.CampusOperationLog) {
	for i := range items {
		fillCampusOperationLogTexts(&items[i])
	}
}

func fillCampusOperationLogTexts(item *campusModel.CampusOperationLog) {
	item.OperatorSourceText = buildCampusOperationSourceText(item.OperatorSource)
	item.ModuleText = buildCampusOperationModuleText(item.Module)
	item.ActionText = buildCampusOperationActionText(item.Action)
}

func normalizeCampusOperationSource(source string) string {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case "web", "cli", "agent":
		return strings.ToLower(strings.TrimSpace(source))
	default:
		return ""
	}
}

func joinCampusAuditLabel(parts ...string) string {
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return strings.Join(items, " / ")
}

func buildCampusAuditIDLabel(prefix string, id uint) string {
	return fmt.Sprintf("%s#%d", prefix, id)
}
