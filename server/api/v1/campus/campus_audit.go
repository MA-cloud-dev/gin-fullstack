package campus

import (
	"errors"
	"strings"

	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

func buildCampusAuditMeta(c *gin.Context) campusReq.CampusAuditMeta {
	userInfo := utils.GetUserInfo(c)
	meta := campusReq.CampusAuditMeta{
		OperatorSource: normalizeCampusAuditSource(c.GetHeader("x-operator-source")),
		OperatorIP:     c.ClientIP(),
		RequestPath:    c.Request.URL.Path,
		RequestMethod:  c.Request.Method,
	}
	if userInfo != nil {
		meta.OperatorSysUserID = userInfo.BaseClaims.ID
		meta.OperatorUsername = userInfo.BaseClaims.Username
		meta.OperatorNickname = userInfo.BaseClaims.NickName
	}
	return meta
}

func normalizeCampusAuditSource(source string) string {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case "cli":
		return "cli"
	default:
		return "web"
	}
}

func buildAgentReviewAuditMeta(c *gin.Context) campusReq.CampusAuditMeta {
	return campusReq.CampusAuditMeta{
		OperatorUsername: "agent",
		OperatorNickname: "Agent",
		OperatorSource:   "agent",
		OperatorIP:       c.ClientIP(),
		RequestPath:      c.Request.URL.Path,
		RequestMethod:    c.Request.Method,
	}
}

func normalizeRequiredAuditReason(reason string) (string, error) {
	trimmed := strings.TrimSpace(reason)
	if trimmed == "" {
		return "", errors.New("操作原因不能为空")
	}
	return trimmed, nil
}

func normalizeOptionalAuditReason(reason string) string {
	return strings.TrimSpace(reason)
}
