package campus

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusOrderApi struct{}

func (a *CampusOrderApi) GetCampusOrderList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo campusReq.CampusOrderSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := campusOrderService.GetCampusOrderInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取交易信息列表失败!", zap.Error(err))
		response.FailWithMessage("获取交易信息列表失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *CampusOrderApi) FindCampusOrder(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("订单ID不合法", c)
		return
	}

	order, err := campusOrderService.GetCampusOrder(ctx, uint(id))
	if err != nil {
		global.GVA_LOG.Error("获取交易信息详情失败!", zap.Error(err))
		response.FailWithMessage("获取交易信息详情失败:"+err.Error(), c)
		return
	}
	response.OkWithData(order, c)
}
