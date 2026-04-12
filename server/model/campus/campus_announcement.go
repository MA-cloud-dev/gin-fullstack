package campus

import "time"

type CampusAnnouncement struct {
	ID                uint      `json:"id" gorm:"column:id;primaryKey"`
	Title             string    `json:"title" gorm:"column:title"`
	Content           string    `json:"content" gorm:"column:content"`
	PublisherID       uint      `json:"publisherId" gorm:"column:publisher_id"`
	Status            int       `json:"status" gorm:"column:status"`
	CreatedAt         time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"column:updated_at"`
	PublisherNickname *string   `json:"publisherNickname" gorm:"column:publisher_nickname;->"`
	PublisherPhone    *string   `json:"publisherPhone" gorm:"column:publisher_phone;->"`
	StatusText        string    `json:"statusText" gorm:"-"`
}

func (CampusAnnouncement) TableName() string {
	return "t_announcement"
}
