package campus

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusCategoryApi struct{}

func (a *CampusCategoryApi) GetCampusCategoryList(c *gin.Context) {
	ctx := c.Request.Context()

	var searchInfo campusReq.CampusCategorySearch
	if err := c.ShouldBindQuery(&searchInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, err := campusCategoryService.GetCampusCategoryTree(ctx, searchInfo)
	if err != nil {
		global.GVA_LOG.Error("获取分类列表失败!", zap.Error(err))
		response.FailWithMessage("获取分类列表失败:"+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

func (a *CampusCategoryApi) FindCampusCategory(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("分类ID不合法", c)
		return
	}

	category, err := campusCategoryService.GetCampusCategory(ctx, uint(id))
	if err != nil {
		global.GVA_LOG.Error("获取分类详情失败!", zap.Error(err))
		response.FailWithMessage("获取分类详情失败:"+err.Error(), c)
		return
	}
	response.OkWithData(category, c)
}

func (a *CampusCategoryApi) CreateCampusCategory(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.CreateCampusCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.AuditReason = normalizeOptionalAuditReason(req.AuditReason)

	if err := campusCategoryService.CreateCampusCategory(ctx, req, buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("创建分类失败!", zap.Error(err))
		response.FailWithMessage("创建分类失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *CampusCategoryApi) UpdateCampusCategory(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.UpdateCampusCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.AuditReason = normalizeOptionalAuditReason(req.AuditReason)

	if err := campusCategoryService.UpdateCampusCategory(ctx, req, buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("更新分类失败!", zap.Error(err))
		response.FailWithMessage("更新分类失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *CampusCategoryApi) UpdateCampusCategoryStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.UpdateCampusCategoryStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.AuditReason = normalizeOptionalAuditReason(req.AuditReason)

	if err := campusCategoryService.UpdateCampusCategoryStatus(ctx, req, buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("更新分类状态失败!", zap.Error(err))
		response.FailWithMessage("更新分类状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
