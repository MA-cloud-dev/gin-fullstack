package campus

import (
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusAuthApi struct{}

const agentReviewCallbackTokenHeader = "X-Agent-Callback-Token"

// GetCampusAuthList 分页获取校园身份审核列表
func (a *CampusAuthApi) GetCampusAuthList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo campusReq.CampusAuthSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := campusAuthService.GetCampusAuthInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取校园身份审核列表失败!", zap.Error(err))
		response.FailWithMessage("获取校园身份审核列表失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// FindCampusAuth 根据ID获取校园身份审核详情
func (a *CampusAuthApi) FindCampusAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	auth, err := campusAuthService.GetCampusAuth(ctx, req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取校园身份审核详情失败!", zap.Error(err))
		response.FailWithMessage("获取校园身份审核详情失败:"+err.Error(), c)
		return
	}

	response.OkWithData(auth, c)
}

// SubmitCampusAuthTest 是 B 端模拟 C 端校园认证提交流程的测试入口。
func (a *CampusAuthApi) SubmitCampusAuthTest(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.SubmitCampusAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := campusAuthService.SubmitCampusAuth(ctx, req)
	if err != nil {
		global.GVA_LOG.Error("提交校园身份审核申请失败!", zap.Error(err))
		response.FailWithMessage("提交失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(resp, "申请已提交，审核中", c)
}

// ReviewCampusAuth 人工通过校园身份申请
func (a *CampusAuthApi) ReviewCampusAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.ReviewCampusAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	auditReason, err := normalizeRequiredAuditReason(req.AuditReason)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.AuditReason = auditReason
	req.ReviewRemark = auditReason

	if err := campusAuthService.ReviewCampusAuth(ctx, req, utils.GetUserID(c), buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("校园身份审核通过失败!", zap.Error(err))
		response.FailWithMessage("审核失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("审核成功", c)
}

func (a *CampusAuthApi) RejectCampusAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.ReviewCampusAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	auditReason, err := normalizeRequiredAuditReason(req.AuditReason)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.AuditReason = auditReason
	req.ReviewRemark = auditReason

	if err := campusAuthService.RejectCampusAuth(ctx, req, utils.GetUserID(c), buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("校园身份审核拒绝失败!", zap.Error(err))
		response.FailWithMessage("审核失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("拒绝成功", c)
}

// RevokeCampusAuth 撤回校园身份审核
func (a *CampusAuthApi) RevokeCampusAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.ReviewCampusAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	auditReason, err := normalizeRequiredAuditReason(req.AuditReason)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.AuditReason = auditReason

	if err := campusAuthService.RevokeCampusAuth(ctx, req, buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("撤回校园身份审核失败!", zap.Error(err))
		response.FailWithMessage("撤回失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("撤回成功", c)
}

func (a *CampusAuthApi) AgentReviewCallback(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.AgentReviewCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := campusAuthService.HandleAgentReviewCallback(ctx, c.GetHeader(agentReviewCallbackTokenHeader), req, buildAgentReviewAuditMeta(c)); err != nil {
		if strings.Contains(err.Error(), "callback token") {
			response.NoAuth(err.Error(), c)
			return
		}
		global.GVA_LOG.Error("处理 Agent 审核回调失败!", zap.Error(err))
		response.FailWithMessage("处理失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("回调处理成功", c)
}
