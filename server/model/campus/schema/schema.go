package schema

import (
	"time"

	"gorm.io/datatypes"
)

type CampusUser struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	Phone      string    `gorm:"column:phone;size:32;index"`
	OpenID     *string   `gorm:"column:openid;size:128;index"`
	Nickname   *string   `gorm:"column:nickname;size:128"`
	AvatarURL  *string   `gorm:"column:avatar_url;size:512"`
	WechatID   *string   `gorm:"column:wechat_id;size:128;index"`
	Grade      *string   `gorm:"column:grade;size:64"`
	Dormitory  *string   `gorm:"column:dormitory;size:128"`
	Username   *string   `gorm:"column:username;size:128;index"`
	Role       int       `gorm:"column:role;index"`
	Status     int       `gorm:"column:status;index"`
	AuthStatus int       `gorm:"column:auth_status;index"`
	CreatedAt  time.Time `gorm:"column:created_at;index"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (CampusUser) TableName() string {
	return "t_user"
}

type CampusProduct struct {
	ID            uint      `gorm:"column:id;primaryKey"`
	UserID        uint      `gorm:"column:user_id;index"`
	CategoryID    uint      `gorm:"column:category_id;index"`
	Title         string    `gorm:"column:title;size:255;index"`
	Description   *string   `gorm:"column:description;type:text"`
	OriginalPrice *float64  `gorm:"column:original_price"`
	Price         float64   `gorm:"column:price"`
	Quality       int       `gorm:"column:quality"`
	TradeMode     int       `gorm:"column:trade_mode;index"`
	ContactInfo   string    `gorm:"column:contact_info;size:255"`
	Status        int       `gorm:"column:status;index"`
	Version       int       `gorm:"column:version"`
	ViewCount     int       `gorm:"column:view_count"`
	WantCount     int       `gorm:"column:want_count"`
	ExpireAt      time.Time `gorm:"column:expire_at;index"`
	CreatedAt     time.Time `gorm:"column:created_at;index"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (CampusProduct) TableName() string {
	return "t_product"
}

type CampusProductImage struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	ProductID uint      `gorm:"column:product_id;index"`
	ImageURL  string    `gorm:"column:image_url;size:512"`
	SortOrder int       `gorm:"column:sort_order;index"`
	CreatedAt time.Time `gorm:"column:created_at;index"`
}

func (CampusProductImage) TableName() string {
	return "t_product_image"
}

type CampusCategory struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;size:128;index"`
	ParentID  *uint     `gorm:"column:parent_id;index"`
	SortOrder int       `gorm:"column:sort_order;index"`
	Icon      *string   `gorm:"column:icon;size:255"`
	Status    int       `gorm:"column:status;index"`
	CreatedAt time.Time `gorm:"column:created_at;index"`
}

func (CampusCategory) TableName() string {
	return "t_category"
}

type CampusAdminStaff struct {
	ID          uint       `gorm:"column:id;primaryKey"`
	Username    string     `gorm:"column:username;size:128;index"`
	DisplayName string     `gorm:"column:display_name;size:128;index"`
	RoleType    int        `gorm:"column:role_type;index"`
	Status      int        `gorm:"column:status;index"`
	CreatedBy   *uint      `gorm:"column:created_by;index"`
	UpdatedBy   *uint      `gorm:"column:updated_by;index"`
	LastLoginAt *time.Time `gorm:"column:last_login_at;index"`
	CreatedAt   time.Time  `gorm:"column:created_at;index"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (CampusAdminStaff) TableName() string {
	return "t_admin_staff"
}

type CampusAnnouncement struct {
	ID          uint      `gorm:"column:id;primaryKey"`
	Title       string    `gorm:"column:title;size:255;index"`
	Content     string    `gorm:"column:content;type:text"`
	PublisherID uint      `gorm:"column:publisher_id;index"`
	Status      int       `gorm:"column:status;index"`
	CreatedAt   time.Time `gorm:"column:created_at;index"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (CampusAnnouncement) TableName() string {
	return "t_announcement"
}

type CampusOrder struct {
	ID             uint       `gorm:"column:id;primaryKey"`
	OrderNo        string     `gorm:"column:order_no;size:64;index"`
	BuyerID        uint       `gorm:"column:buyer_id;index"`
	SellerID       uint       `gorm:"column:seller_id;index"`
	ProductID      uint       `gorm:"column:product_id;index"`
	ProductTitle   string     `gorm:"column:product_title;size:255"`
	ProductImage   *string    `gorm:"column:product_image;size:512"`
	Price          float64    `gorm:"column:price"`
	Status         int        `gorm:"column:status;index"`
	CloseReason    *string    `gorm:"column:close_reason;type:text"`
	CloseBy        *int       `gorm:"column:close_by;index"`
	CloseConfirmed *int       `gorm:"column:close_confirmed;index"`
	Remark         *string    `gorm:"column:remark;type:text"`
	ConfirmedAt    *time.Time `gorm:"column:confirmed_at;index"`
	CompletedAt    *time.Time `gorm:"column:completed_at;index"`
	CancelledAt    *time.Time `gorm:"column:cancelled_at;index"`
	CreatedAt      time.Time  `gorm:"column:created_at;index"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
}

func (CampusOrder) TableName() string {
	return "t_order"
}

type CampusReport struct {
	ID           uint      `gorm:"column:id;primaryKey"`
	ReporterID   uint      `gorm:"column:reporter_id;index"`
	TargetType   int       `gorm:"column:target_type;index"`
	TargetID     uint      `gorm:"column:target_id;index"`
	Reason       int       `gorm:"column:reason;index"`
	Description  *string   `gorm:"column:description;type:text"`
	Status       int       `gorm:"column:status;index"`
	HandledBy    *uint     `gorm:"column:handled_by;index"`
	HandleResult *string   `gorm:"column:handle_result;type:text"`
	CreatedAt    time.Time `gorm:"column:created_at;index"`
}

func (CampusReport) TableName() string {
	return "t_report"
}

type CampusAuth struct {
	ID           uint       `gorm:"column:id;primaryKey"`
	UserID       uint       `gorm:"column:user_id;index"`
	StudentID    string     `gorm:"column:student_id;size:64;index"`
	RealName     string     `gorm:"column:real_name;size:128;index"`
	College      string     `gorm:"column:college;size:255;index"`
	ReviewStatus string     `gorm:"column:review_status;size:32;index"`
	ReviewSource *string    `gorm:"column:review_source;size:32"`
	ReviewRemark *string    `gorm:"column:review_remark;type:text"`
	ReviewedAt   *time.Time `gorm:"column:reviewed_at;index"`
	ReviewedBy   *uint      `gorm:"column:reviewed_by;index"`
	CreatedAt    time.Time  `gorm:"column:created_at;index"`
}

func (CampusAuth) TableName() string {
	return "t_campus_auth"
}

type CampusOperationLog struct {
	ID                uint      `gorm:"column:id;primaryKey"`
	OperatorSysUserID uint      `gorm:"column:operator_sys_user_id;index"`
	OperatorUsername  string    `gorm:"column:operator_username;size:128;index"`
	OperatorNickname  string    `gorm:"column:operator_nickname;size:128"`
	OperatorSource    string    `gorm:"column:operator_source;size:32;index"`
	Module            string    `gorm:"column:module;size:64;index"`
	Action            string    `gorm:"column:action;size:64;index"`
	TargetID          uint      `gorm:"column:target_id;index"`
	TargetLabel       string    `gorm:"column:target_label;size:255"`
	Reason            string    `gorm:"column:reason;type:text"`
	Result            string    `gorm:"column:result;type:text"`
	RequestPath       string    `gorm:"column:request_path;size:255"`
	RequestMethod     string    `gorm:"column:request_method;size:16"`
	OperatorIP        string    `gorm:"column:operator_ip;size:64"`
	CreatedAt         time.Time `gorm:"column:created_at;index"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (CampusOperationLog) TableName() string {
	return "t_campus_operation_log"
}

type AgentReviewTask struct {
	ID                     uint           `gorm:"column:id;primaryKey"`
	TaskID                 string         `gorm:"column:task_id;size:64;uniqueIndex"`
	BizType                string         `gorm:"column:biz_type;size:64;index"`
	BizID                  uint64         `gorm:"column:biz_id;index"`
	TraceID                string         `gorm:"column:trace_id;size:64;index"`
	CallbackURL            string         `gorm:"column:callback_url;size:512"`
	PayloadJSON            datatypes.JSON `gorm:"column:payload_json;type:json"`
	TaskStatus             string         `gorm:"column:task_status;size:32;index"`
	FinalAction            string         `gorm:"column:final_action;size:32"`
	ReviewRemark           string         `gorm:"column:review_remark;type:text"`
	ResultJSON             datatypes.JSON `gorm:"column:result_json;type:json"`
	DispatchAttemptCount   int            `gorm:"column:dispatch_attempt_count"`
	DispatchError          string         `gorm:"column:dispatch_error;type:text"`
	CallbackError          string         `gorm:"column:callback_error;type:text"`
	PreviousUserAuthStatus int            `gorm:"column:previous_user_auth_status"`
	AuthRecordID           uint64         `gorm:"column:auth_record_id;index"`
	UserID                 uint64         `gorm:"column:user_id;index"`
	CreatedAt              time.Time      `gorm:"column:created_at;index"`
	UpdatedAt              time.Time      `gorm:"column:updated_at"`
	CallbackAt             *time.Time     `gorm:"column:callback_at;index"`
}

func (AgentReviewTask) TableName() string {
	return "t_agent_review_task"
}

func Tables() []interface{} {
	return []interface{}{
		CampusUser{},
		CampusProduct{},
		CampusProductImage{},
		CampusCategory{},
		CampusAdminStaff{},
		CampusAnnouncement{},
		CampusOrder{},
		CampusReport{},
		CampusAuth{},
		CampusOperationLog{},
		AgentReviewTask{},
	}
}
