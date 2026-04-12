package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusAdminStaffApi struct{}

func (a *CampusAdminStaffApi) GetCampusAdminStaffList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo campusReq.CampusAdminStaffSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := campusAdminStaffService.GetCampusAdminStaffInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取B端管理员列表失败!", zap.Error(err))
		response.FailWithMessage("获取B端管理员列表失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *CampusAdminStaffApi) FindCampusAdminStaff(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		ID uint `form:"id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	admin, err := campusAdminStaffService.GetCampusAdminStaff(ctx, req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取B端管理员详情失败!", zap.Error(err))
		response.FailWithMessage("获取B端管理员详情失败:"+err.Error(), c)
		return
	}
	response.OkWithData(admin, c)
}

func (a *CampusAdminStaffApi) UpdateCampusAdminStaffStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.UpdateCampusAdminStaffStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := campusAdminStaffService.UpdateCampusAdminStaffStatus(ctx, req); err != nil {
		global.GVA_LOG.Error("更新B端管理员状态失败!", zap.Error(err))
		response.FailWithMessage("更新B端管理员状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
