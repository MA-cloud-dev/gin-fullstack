package campus

import "time"

type CampusAdminStaff struct {
	ID           uint       `json:"id" gorm:"column:id;primaryKey"`
	Username     string     `json:"username" gorm:"column:username"`
	DisplayName  string     `json:"displayName" gorm:"column:display_name"`
	RoleType     int        `json:"roleType" gorm:"column:role_type"`
	Status       int        `json:"status" gorm:"column:status"`
	CreatedBy    *uint      `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy    *uint      `json:"updatedBy" gorm:"column:updated_by"`
	LastLoginAt  *time.Time `json:"lastLoginAt" gorm:"column:last_login_at"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	RoleTypeText string     `json:"roleTypeText" gorm:"-"`
	StatusText   string     `json:"statusText" gorm:"-"`
}

func (CampusAdminStaff) TableName() string {
	return "t_admin_staff"
}
