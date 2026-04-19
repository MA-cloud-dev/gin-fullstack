package campus

import (
	"time"

	"gorm.io/datatypes"
)

const (
	AgentReviewBizTypeCampusAuth = "campus_auth"

	AgentReviewTaskStatusPendingDispatch = "pending_dispatch"
	AgentReviewTaskStatusSubmitted       = "submitted"
	AgentReviewTaskStatusProcessing      = "processing"
	AgentReviewTaskStatusCompleted       = "completed"
	AgentReviewTaskStatusFailed          = "failed"
	AgentReviewTaskStatusSuperseded      = "superseded"

	AgentReviewFinalActionApprove  = "approve"
	AgentReviewFinalActionReject   = "reject"
	AgentReviewFinalActionEscalate = "escalate"

	AgentReviewCallbackTaskStatusProcessing = "processing"
	AgentReviewCallbackTaskStatusCompleted  = "completed"
	AgentReviewCallbackTaskStatusFailed     = "failed"
)

type AgentReviewTask struct {
	ID                     uint           `json:"id" gorm:"column:id;primaryKey"`
	TaskID                 string         `json:"taskId" gorm:"column:task_id;size:64;uniqueIndex"`
	BizType                string         `json:"bizType" gorm:"column:biz_type;size:64;index"`
	BizID                  uint64         `json:"bizId" gorm:"column:biz_id;index"`
	TraceID                string         `json:"traceId" gorm:"column:trace_id;size:64;index"`
	CallbackURL            string         `json:"callbackUrl" gorm:"column:callback_url;size:512"`
	PayloadJSON            datatypes.JSON `json:"payloadJson" gorm:"column:payload_json;type:json"`
	TaskStatus             string         `json:"taskStatus" gorm:"column:task_status;size:32;index"`
	FinalAction            string         `json:"finalAction" gorm:"column:final_action;size:32"`
	ReviewRemark           string         `json:"reviewRemark" gorm:"column:review_remark;type:text"`
	ResultJSON             datatypes.JSON `json:"resultJson" gorm:"column:result_json;type:json"`
	DispatchAttemptCount   int            `json:"dispatchAttemptCount" gorm:"column:dispatch_attempt_count"`
	DispatchError          string         `json:"dispatchError" gorm:"column:dispatch_error;type:text"`
	CallbackError          string         `json:"callbackError" gorm:"column:callback_error;type:text"`
	PreviousUserAuthStatus int            `json:"previousUserAuthStatus" gorm:"column:previous_user_auth_status"`
	AuthRecordID           uint64         `json:"authRecordId" gorm:"column:auth_record_id;index"`
	UserID                 uint64         `json:"userId" gorm:"column:user_id;index"`
	CreatedAt              time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt              time.Time      `json:"updatedAt" gorm:"column:updated_at"`
	CallbackAt             *time.Time     `json:"callbackAt" gorm:"column:callback_at"`
}

func (AgentReviewTask) TableName() string {
	return "t_agent_review_task"
}
