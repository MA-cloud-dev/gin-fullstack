package campus

import "time"

type CampusProduct struct {
	ID                uint                 `json:"id" gorm:"column:id;primaryKey"`
	UserID            uint                 `json:"userId" gorm:"column:user_id"`
	CategoryID        uint                 `json:"categoryId" gorm:"column:category_id"`
	Title             string               `json:"title" gorm:"column:title"`
	Description       *string              `json:"description" gorm:"column:description"`
	OriginalPrice     *float64             `json:"originalPrice" gorm:"column:original_price"`
	Price             float64              `json:"price" gorm:"column:price"`
	Quality           int                  `json:"quality" gorm:"column:quality"`
	TradeMode         int                  `json:"tradeMode" gorm:"column:trade_mode"`
	ContactInfo       string               `json:"contactInfo" gorm:"column:contact_info"`
	Status            int                  `json:"status" gorm:"column:status"`
	Version           int                  `json:"version" gorm:"column:version"`
	ViewCount         int                  `json:"viewCount" gorm:"column:view_count"`
	WantCount         int                  `json:"wantCount" gorm:"column:want_count"`
	ExpireAt          time.Time            `json:"expireAt" gorm:"column:expire_at"`
	CreatedAt         time.Time            `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt         time.Time            `json:"updatedAt" gorm:"column:updated_at"`
	CategoryName      *string              `json:"categoryName" gorm:"column:category_name;->"`
	PublisherNickname *string              `json:"publisherNickname" gorm:"column:publisher_nickname;->"`
	PublisherPhone    *string              `json:"publisherPhone" gorm:"column:publisher_phone;->"`
	CoverURL          *string              `json:"coverUrl" gorm:"column:cover_url;->"`
	StatusText        string               `json:"statusText" gorm:"-"`
	TradeModeText     string               `json:"tradeModeText" gorm:"-"`
	Images            []CampusProductImage `json:"images,omitempty" gorm:"-"`
}

func (CampusProduct) TableName() string {
	return "t_product"
}

type CampusProductImage struct {
	ID        uint      `json:"id" gorm:"column:id;primaryKey"`
	ProductID uint      `json:"productId" gorm:"column:product_id"`
	ImageURL  string    `json:"imageUrl" gorm:"column:image_url"`
	SortOrder int       `json:"sortOrder" gorm:"column:sort_order"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (CampusProductImage) TableName() string {
	return "t_product_image"
}
