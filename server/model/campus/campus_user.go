package campus

import "time"

type CampusUser struct {
	ID             uint       `json:"id" gorm:"column:id;primaryKey"`
	Phone          string     `json:"phone" gorm:"column:phone"`
	OpenID         *string    `json:"openid" gorm:"column:openid"`
	Nickname       *string    `json:"nickname" gorm:"column:nickname"`
	AvatarURL      *string    `json:"avatarUrl" gorm:"column:avatar_url"`
	WechatID       *string    `json:"wechatId" gorm:"column:wechat_id"`
	Grade          *string    `json:"grade" gorm:"column:grade"`
	Dormitory      *string    `json:"dormitory" gorm:"column:dormitory"`
	Username       *string    `json:"username" gorm:"column:username"`
	Role           int        `json:"role" gorm:"column:role"`
	Status         int        `json:"status" gorm:"column:status"`
	AuthStatus     int        `json:"authStatus" gorm:"column:auth_status"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	StudentID      *string    `json:"studentId" gorm:"column:student_id;->"`
	RealName       *string    `json:"realName" gorm:"column:real_name;->"`
	College        *string    `json:"college" gorm:"column:college;->"`
	ReviewRemark   *string    `json:"reviewRemark" gorm:"column:review_remark;->"`
	ReviewedAt     *time.Time `json:"reviewedAt" gorm:"column:reviewed_at;->"`
	ReviewedBy     *uint      `json:"reviewedBy" gorm:"column:reviewed_by;->"`
	ReviewedByName *string    `json:"reviewedByName" gorm:"column:reviewed_by_name;->"`
	StatusText     string     `json:"statusText" gorm:"-"`
	RoleText       string     `json:"roleText" gorm:"-"`
	AuthStatusText string     `json:"authStatusText" gorm:"-"`
}

func (CampusUser) TableName() string {
	return "t_user"
}
