package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusUserApi struct{}

func (a *CampusUserApi) GetCampusUserList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo campusReq.CampusUserSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := campusUserService.GetCampusUserInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取校园用户列表失败!", zap.Error(err))
		response.FailWithMessage("获取校园用户列表失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *CampusUserApi) FindCampusUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	user, err := campusUserService.GetCampusUser(ctx, req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取校园用户详情失败!", zap.Error(err))
		response.FailWithMessage("获取校园用户详情失败:"+err.Error(), c)
		return
	}
	response.OkWithData(user, c)
}

func (a *CampusUserApi) UpdateCampusUserStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.UpdateCampusUserStatusReq
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
	if err := campusUserService.UpdateCampusUserStatus(ctx, req, buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("更新校园用户状态失败!", zap.Error(err))
		response.FailWithMessage("更新校园用户状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
