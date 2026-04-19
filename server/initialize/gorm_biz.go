package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	"gorm.io/gorm"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(&campusModel.CampusOperationLog{}, &campusModel.CampusAuth{}, &campusModel.AgentReviewTask{})
	if err != nil {
		return err
	}
	if err = db.Exec(`
		UPDATE t_campus_auth
		SET review_status = CASE
			WHEN reviewed_at IS NOT NULL THEN 'approved'
			ELSE 'processing'
		END
		WHERE review_status IS NULL OR review_status = ''
	`).Error; err != nil {
		return err
	}
	if err = db.Exec(`
		UPDATE t_campus_auth
		SET review_source = 'manual'
		WHERE reviewed_at IS NOT NULL
		  AND reviewed_by IS NOT NULL
		  AND (review_source IS NULL OR review_source = '')
	`).Error; err != nil {
		return err
	}
	if err = syncCampusUserAuthStatuses(db); err != nil {
		return err
	}
	return nil
}
func syncCampusUserAuthStatuses(db *gorm.DB) error {
	type latestCampusAuthStatus struct {
		UserID       uint   `gorm:"column:user_id"`
		ReviewStatus string `gorm:"column:review_status"`
	}
	statusByUser := map[uint]int{}
	var latestStatuses []latestCampusAuthStatus
	if err := db.Table("t_campus_auth AS ca").Select("ca.user_id, ca.review_status").Joins("INNER JOIN (SELECT user_id, MAX(id) AS max_id FROM t_campus_auth GROUP BY user_id) AS latest ON latest.max_id = ca.id").Scan(&latestStatuses).Error; err != nil {
		return err
	}
	for _, item := range latestStatuses {
		switch item.ReviewStatus {
		case campusModel.CampusAuthReviewStatusApproved:
			statusByUser[item.UserID] = campusModel.CampusUserAuthStatusApproved
		case campusModel.CampusAuthReviewStatusRejected:
			statusByUser[item.UserID] = campusModel.CampusUserAuthStatusRejected
		case campusModel.CampusAuthReviewStatusProcessing:
			statusByUser[item.UserID] = campusModel.CampusUserAuthStatusProcessing
		default:
			statusByUser[item.UserID] = campusModel.CampusUserAuthStatusUnverified
		}
	}
	var activeTaskUserIDs []uint64
	if err := db.Table(campusModel.AgentReviewTask{}.TableName()).Distinct("user_id").Where("biz_type = ? AND task_status IN ?", campusModel.AgentReviewBizTypeCampusAuth, []string{campusModel.AgentReviewTaskStatusPendingDispatch, campusModel.AgentReviewTaskStatusSubmitted, campusModel.AgentReviewTaskStatusProcessing}).Pluck("user_id", &activeTaskUserIDs).Error; err != nil {
		return err
	}
	for _, userID := range activeTaskUserIDs {
		statusByUser[uint(userID)] = campusModel.CampusUserAuthStatusProcessing
	}
	var users []campusModel.CampusUser
	if err := db.Table(campusModel.CampusUser{}.TableName()).Select("id, role, auth_status").Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		if user.Role == 1 {
			continue
		}
		targetStatus, ok := statusByUser[user.ID]
		if !ok {
			targetStatus = campusModel.CampusUserAuthStatusUnverified
		}
		if user.AuthStatus == targetStatus {
			continue
		}
		if err := db.Table(campusModel.CampusUser{}.TableName()).Where("id = ?", user.ID).Update("auth_status", targetStatus).Error; err != nil {
			return err
		}
	}
	return nil
}
