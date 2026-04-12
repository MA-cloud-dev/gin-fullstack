package campus

import "time"

type CampusCategory struct {
	ID         uint             `json:"id" gorm:"column:id;primaryKey"`
	Name       string           `json:"name" gorm:"column:name"`
	ParentID   *uint            `json:"parentId" gorm:"column:parent_id"`
	SortOrder  int              `json:"sortOrder" gorm:"column:sort_order"`
	Icon       *string          `json:"icon" gorm:"column:icon"`
	Status     int              `json:"status" gorm:"column:status"`
	CreatedAt  time.Time        `json:"createdAt" gorm:"column:created_at"`
	ParentName *string          `json:"parentName" gorm:"column:parent_name;->"`
	StatusText string           `json:"statusText" gorm:"-"`
	Children   []CampusCategory `json:"children,omitempty" gorm:"-"`
}

func (CampusCategory) TableName() string {
	return "t_category"
}
