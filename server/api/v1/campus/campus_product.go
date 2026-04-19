package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusProductApi struct{}

func (a *CampusProductApi) GetCampusProductList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo campusReq.CampusProductSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := campusProductService.GetCampusProductInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取商品管理列表失败!", zap.Error(err))
		response.FailWithMessage("获取商品管理列表失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *CampusProductApi) FindCampusProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	product, err := campusProductService.GetCampusProduct(ctx, req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取商品详情失败!", zap.Error(err))
		response.FailWithMessage("获取商品详情失败:"+err.Error(), c)
		return
	}
	response.OkWithData(product, c)
}

func (a *CampusProductApi) UpdateCampusProductStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.UpdateCampusProductStatusReq
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
	if err := campusProductService.UpdateCampusProductStatus(ctx, req, buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("更新商品状态失败!", zap.Error(err))
		response.FailWithMessage("更新商品状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
