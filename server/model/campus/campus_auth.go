package campus

import "time"

// CampusAuth 映射校园身份审核申请表，只做现有业务表读取和审核回写。
type CampusAuth struct {
	ID             uint       `json:"id" gorm:"column:id;primaryKey"`
	UserID         uint       `json:"userId" gorm:"column:user_id"`
	StudentID      string     `json:"studentId" gorm:"column:student_id"`
	RealName       string     `json:"realName" gorm:"column:real_name"`
	College        string     `json:"college" gorm:"column:college"`
	ReviewRemark   *string    `json:"reviewRemark" gorm:"column:review_remark"`
	ReviewedAt     *time.Time `json:"reviewedAt" gorm:"column:reviewed_at"`
	ReviewedBy     *uint      `json:"reviewedBy" gorm:"column:reviewed_by"`
	ReviewedByName *string    `json:"reviewedByName" gorm:"column:reviewed_by_name;->"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at"`
}

func (CampusAuth) TableName() string {
	return "t_campus_auth"
}
