package request

import (
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type CampusAuditMeta struct {
	OperatorSysUserID uint   `json:"operatorSysUserId"`
	OperatorUsername  string `json:"operatorUsername"`
	OperatorNickname  string `json:"operatorNickname"`
	OperatorSource    string `json:"operatorSource"`
	OperatorIP        string `json:"operatorIp"`
	RequestPath       string `json:"requestPath"`
	RequestMethod     string `json:"requestMethod"`
}

type CampusOperationLogSearch struct {
	OperatorKeyword string      `json:"operatorKeyword" form:"operatorKeyword"`
	OperatorSource  string      `json:"operatorSource" form:"operatorSource"`
	Module          string      `json:"module" form:"module"`
	Action          string      `json:"action" form:"action"`
	TargetID        string      `json:"targetId" form:"targetId"`
	CreatedAtRange  []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	commonReq.PageInfo
}
