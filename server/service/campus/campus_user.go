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

type CampusUserService struct{}

func (s *CampusUserService) GetCampusUserInfoList(ctx context.Context, info campusReq.CampusUserSearch) (list []campusModel.CampusUser, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.WithContext(ctx).
		Table("t_user AS u").
		Select("u.*, ca.id AS auth_record_id, ca.student_id, ca.real_name, ca.college, ca.review_status, ca.review_source, ca.review_remark, ca.reviewed_at, ca.reviewed_by, su.nick_name AS reviewed_by_name").
		Joins("LEFT JOIN (?) AS ca ON ca.user_id = u.id", latestCampusAuthSubQuery(ctx)).
		Joins("LEFT JOIN sys_users AS su ON su.id = ca.reviewed_by")

	if id, ok := parseUintFromText(info.ID); ok {
		db = db.Where("u.id = ?", id)
	}
	if phone := strings.TrimSpace(info.Phone); phone != "" {
		db = db.Where("u.phone LIKE ?", "%"+phone+"%")
	}
	if nickname := strings.TrimSpace(info.Nickname); nickname != "" {
		db = db.Where("u.nickname LIKE ?", "%"+nickname+"%")
	}
	if info.Role != nil {
		db = db.Where("u.role = ?", *info.Role)
	}
	if info.Status != nil {
		db = db.Where("u.status = ?", *info.Status)
	}
	if info.AuthStatus != nil {
		db = db.Where("u.auth_status = ?", *info.AuthStatus)
	}
	if studentID := strings.TrimSpace(info.StudentID); studentID != "" {
		db = db.Where("ca.student_id LIKE ?", "%"+studentID+"%")
	}
	if realName := strings.TrimSpace(info.RealName); realName != "" {
		db = db.Where("ca.real_name LIKE ?", "%"+realName+"%")
	}
	if college := strings.TrimSpace(info.College); college != "" {
		db = db.Where("ca.college LIKE ?", "%"+college+"%")
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("u.created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []campusModel.CampusUser
	query := db.Order("u.created_at DESC")
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err = query.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].StatusText = buildStatusText(items[i].Status)
		items[i].RoleText = buildUserRoleText(items[i].Role)
		items[i].AuthStatusText = buildAuthStatusText(items[i])
	}
	return items, total, nil
}

func (s *CampusUserService) GetCampusUser(ctx context.Context, id uint) (user campusModel.CampusUser, err error) {
	err = global.GVA_DB.WithContext(ctx).
		Table("t_user AS u").
		Select("u.*, ca.id AS auth_record_id, ca.student_id, ca.real_name, ca.college, ca.review_status, ca.review_source, ca.review_remark, ca.reviewed_at, ca.reviewed_by, su.nick_name AS reviewed_by_name").
		Joins("LEFT JOIN (?) AS ca ON ca.user_id = u.id", latestCampusAuthSubQuery(ctx)).
		Joins("LEFT JOIN sys_users AS su ON su.id = ca.reviewed_by").
		Where("u.id = ?", id).
		First(&user).Error
	if err == nil {
		user.StatusText = buildStatusText(user.Status)
		user.RoleText = buildUserRoleText(user.Role)
		user.AuthStatusText = buildAuthStatusText(user)
	}
	return
}

func (s *CampusUserService) UpdateCampusUserStatus(ctx context.Context, req campusReq.UpdateCampusUserStatusReq, auditMeta campusReq.CampusAuditMeta) error {
	return global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user campusModel.CampusUser
		if err := tx.Table(user.TableName()).Where("id = ?", req.ID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("校园用户不存在")
			}
			return err
		}
		if err := tx.Table(user.TableName()).Where("id = ?", req.ID).Update("status", *req.Status).Error; err != nil {
			return err
		}
		action := "disable_user"
		if *req.Status == 0 {
			action = "enable_user"
		}
		nickname := ""
		if user.Nickname != nil {
			nickname = *user.Nickname
		}
		return createCampusOperationLogWithTx(tx, campusOperationAuditInput{
			Meta:        auditMeta,
			Module:      "user",
			Action:      action,
			TargetID:    user.ID,
			TargetLabel: joinCampusAuditLabel(strings.TrimSpace(nickname), user.Phone),
			Reason:      req.AuditReason,
			Result:      "用户状态已更新为" + buildStatusText(*req.Status),
		})
	})
}
