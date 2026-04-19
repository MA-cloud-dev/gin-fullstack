package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusOperationLogApi struct{}

func (a *CampusOperationLogApi) GetCampusOperationLogList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo campusReq.CampusOperationLogSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := campusOperationLogService.GetCampusOperationLogInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取校园操作审计列表失败!", zap.Error(err))
		response.FailWithMessage("获取校园操作审计列表失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *CampusOperationLogApi) FindCampusOperationLog(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	item, err := campusOperationLogService.GetCampusOperationLog(ctx, req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取校园操作审计详情失败!", zap.Error(err))
		response.FailWithMessage("获取校园操作审计详情失败:"+err.Error(), c)
		return
	}
	response.OkWithData(item, c)
}
