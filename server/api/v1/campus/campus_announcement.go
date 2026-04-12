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

type CampusAnnouncementApi struct{}

func (a *CampusAnnouncementApi) GetCampusAnnouncementList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo campusReq.CampusAnnouncementSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := campusAnnouncementService.GetCampusAnnouncementInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取公告列表失败!", zap.Error(err))
		response.FailWithMessage("获取公告列表失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *CampusAnnouncementApi) FindCampusAnnouncement(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("公告ID不合法", c)
		return
	}

	announcement, err := campusAnnouncementService.GetCampusAnnouncement(ctx, uint(id))
	if err != nil {
		global.GVA_LOG.Error("获取公告详情失败!", zap.Error(err))
		response.FailWithMessage("获取公告详情失败:"+err.Error(), c)
		return
	}
	response.OkWithData(announcement, c)
}

func (a *CampusAnnouncementApi) CreateCampusAnnouncement(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.CreateCampusAnnouncementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := campusAnnouncementService.CreateCampusAnnouncement(ctx, req); err != nil {
		global.GVA_LOG.Error("创建公告失败!", zap.Error(err))
		response.FailWithMessage("创建公告失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *CampusAnnouncementApi) UpdateCampusAnnouncement(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.UpdateCampusAnnouncementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := campusAnnouncementService.UpdateCampusAnnouncement(ctx, req); err != nil {
		global.GVA_LOG.Error("更新公告失败!", zap.Error(err))
		response.FailWithMessage("更新公告失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *CampusAnnouncementApi) DeleteCampusAnnouncement(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.DeleteCampusAnnouncementReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := campusAnnouncementService.DeleteCampusAnnouncement(ctx, req.ID); err != nil {
		global.GVA_LOG.Error("删除公告失败!", zap.Error(err))
		response.FailWithMessage("删除公告失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (a *CampusAnnouncementApi) UpdateCampusAnnouncementStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.UpdateCampusAnnouncementStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := campusAnnouncementService.UpdateCampusAnnouncementStatus(ctx, req); err != nil {
		global.GVA_LOG.Error("更新公告状态失败!", zap.Error(err))
		response.FailWithMessage("更新公告状态失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
