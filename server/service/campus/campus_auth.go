package campus

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CampusAuthService struct{}

func (s *CampusAuthService) GetCampusAuthInfoList(ctx context.Context, info campusReq.CampusAuthSearch) (list []campusModel.CampusAuth, total int64, err error) {
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
	if reviewStatus := strings.TrimSpace(info.ReviewStatus); reviewStatus != "" {
		db = db.Where("ca.review_status = ?", normalizeCampusAuthReviewStatus(reviewStatus))
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

	db = db.Order("ca.created_at DESC, ca.id DESC")
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	var items []campusModel.CampusAuth
	if err = db.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	fillCampusAuthListViewFields(items)
	return items, total, nil
}

func (s *CampusAuthService) GetCampusAuth(ctx context.Context, id uint) (auth campusModel.CampusAuth, err error) {
	err = global.GVA_DB.WithContext(ctx).
		Table("t_campus_auth AS ca").
		Select("ca.*, su.nick_name AS reviewed_by_name").
		Joins("LEFT JOIN sys_users AS su ON su.id = ca.reviewed_by").
		Where("ca.id = ?", id).
		First(&auth).Error
	if err != nil {
		return auth, err
	}
	applyCampusAuthViewFields(&auth)
	return auth, nil
}

func (s *CampusAuthService) ReviewCampusAuth(ctx context.Context, req campusReq.ReviewCampusAuthReq, reviewerID uint, auditMeta campusReq.CampusAuditMeta) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		auth, err := loadCampusAuthForUpdate(tx, req.ID)
		if err != nil {
			return err
		}

		if normalizeCampusAuthStatusWithLegacy(auth.ReviewStatus, auth.ReviewedAt) == campusModel.CampusAuthReviewStatusApproved {
			return errors.New("该记录已通过审核，请勿重复操作")
		}

		now := time.Now()
		reviewRemark := strings.TrimSpace(req.ReviewRemark)
		reviewSource := campusModel.CampusAuthReviewSourceManual
		updates := map[string]interface{}{
			"review_status": campusModel.CampusAuthReviewStatusApproved,
			"review_source": reviewSource,
			"review_remark": reviewRemark,
			"reviewed_at":   now,
			"reviewed_by":   reviewerID,
		}
		if err := tx.Table(auth.TableName()).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := updateCampusUserAuthStatus(tx, auth.UserID, campusModel.CampusUserAuthStatusApproved); err != nil {
			return err
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "auth",
			Action:      "approve_auth",
			TargetID:    auth.ID,
			TargetLabel: joinCampusAuditLabel(auth.RealName, auth.StudentID),
			Reason:      req.AuditReason,
			Result:      "已通过校园身份审核",
		})
	})
}

func (s *CampusAuthService) RejectCampusAuth(ctx context.Context, req campusReq.ReviewCampusAuthReq, reviewerID uint, auditMeta campusReq.CampusAuditMeta) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		auth, err := loadCampusAuthForUpdate(tx, req.ID)
		if err != nil {
			return err
		}

		if normalizeCampusAuthStatusWithLegacy(auth.ReviewStatus, auth.ReviewedAt) == campusModel.CampusAuthReviewStatusRejected {
			return errors.New("该记录已拒绝，请勿重复操作")
		}

		now := time.Now()
		reviewRemark := strings.TrimSpace(req.ReviewRemark)
		reviewSource := campusModel.CampusAuthReviewSourceManual
		updates := map[string]interface{}{
			"review_status": campusModel.CampusAuthReviewStatusRejected,
			"review_source": reviewSource,
			"review_remark": reviewRemark,
			"reviewed_at":   now,
			"reviewed_by":   reviewerID,
		}
		if err := tx.Table(auth.TableName()).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := updateCampusUserAuthStatus(tx, auth.UserID, campusModel.CampusUserAuthStatusRejected); err != nil {
			return err
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "auth",
			Action:      "reject_auth",
			TargetID:    auth.ID,
			TargetLabel: joinCampusAuditLabel(auth.RealName, auth.StudentID),
			Reason:      req.AuditReason,
			Result:      "已拒绝校园身份审核",
		})
	})
}

func (s *CampusAuthService) RevokeCampusAuth(ctx context.Context, req campusReq.ReviewCampusAuthReq, auditMeta campusReq.CampusAuditMeta) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		auth, err := loadCampusAuthForUpdate(tx, req.ID)
		if err != nil {
			return err
		}

		if normalizeCampusAuthStatusWithLegacy(auth.ReviewStatus, auth.ReviewedAt) == campusModel.CampusAuthReviewStatusProcessing {
			return errors.New("该记录已处于审核中，无需撤回")
		}

		updates := map[string]interface{}{
			"review_status": campusModel.CampusAuthReviewStatusProcessing,
			"review_source": nil,
			"review_remark": nil,
			"reviewed_at":   nil,
			"reviewed_by":   nil,
		}
		if err := tx.Table(auth.TableName()).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return err
		}
		if err := updateCampusUserAuthStatus(tx, auth.UserID, campusModel.CampusUserAuthStatusProcessing); err != nil {
			return err
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "auth",
			Action:      "revoke_auth",
			TargetID:    auth.ID,
			TargetLabel: joinCampusAuditLabel(auth.RealName, auth.StudentID),
			Reason:      req.AuditReason,
			Result:      "已撤回校园身份审核并恢复为审核中",
		})
	})
}

func loadCampusAuthForUpdate(tx *gorm.DB, id uint) (campusModel.CampusAuth, error) {
	var auth campusModel.CampusAuth
	if err := tx.Table(auth.TableName()).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return auth, errors.New("审核记录不存在")
		}
		return auth, err
	}
	return auth, nil
}

func updateCampusUserAuthStatus(tx *gorm.DB, userID uint, authStatus int) error {
	return tx.Table(campusModel.CampusUser{}.TableName()).Where("id = ?", userID).Update("auth_status", authStatus).Error
}

func normalizeCampusAuthStatusWithLegacy(reviewStatus string, reviewedAt *time.Time) string {
	trimmed := strings.TrimSpace(reviewStatus)
	if trimmed != "" {
		return normalizeCampusAuthReviewStatus(trimmed)
	}
	if reviewedAt != nil {
		return campusModel.CampusAuthReviewStatusApproved
	}
	return campusModel.CampusAuthReviewStatusProcessing
}

func applyCampusAuthViewFields(auth *campusModel.CampusAuth) {
	auth.ReviewStatus = normalizeCampusAuthStatusWithLegacy(auth.ReviewStatus, auth.ReviewedAt)
	auth.ReviewStatusText = buildCampusAuthReviewStatusText(auth.ReviewStatus)
}

func fillCampusAuthListViewFields(items []campusModel.CampusAuth) {
	for i := range items {
		applyCampusAuthViewFields(&items[i])
	}
}
