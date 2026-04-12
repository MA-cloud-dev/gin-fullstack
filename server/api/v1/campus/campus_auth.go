package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusAuthApi struct{}

// GetCampusAuthList 分页获取校园身份审核列表
// @Tags CampusAuth
// @Summary 分页获取校园身份审核列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query campusReq.CampusAuthSearch true "分页获取校园身份审核列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /campusAuth/getCampusAuthList [get]
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
// @Tags CampusAuth
// @Summary 根据ID获取校园身份审核详情
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id query uint true "审核ID"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /campusAuth/findCampusAuth [get]
func (a *CampusAuthApi) FindCampusAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.ReviewCampusAuthReq
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

// ReviewCampusAuth 审核校园身份申请
// @Tags CampusAuth
// @Summary 审核校园身份申请
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body campusReq.ReviewCampusAuthReq true "审核参数"
// @Success 200 {object} response.Response{msg=string} "审核成功"
// @Router /campusAuth/reviewCampusAuth [post]
func (a *CampusAuthApi) ReviewCampusAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.ReviewCampusAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := campusAuthService.ReviewCampusAuth(ctx, req, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("校园身份审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("审核成功", c)
}

// RevokeCampusAuth 撤回校园身份审核
// @Tags CampusAuth
// @Summary 撤回校园身份审核
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body campusReq.ReviewCampusAuthReq true "撤回参数"
// @Success 200 {object} response.Response{msg=string} "撤回成功"
// @Router /campusAuth/revokeCampusAuth [post]
func (a *CampusAuthApi) RevokeCampusAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.ReviewCampusAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := campusAuthService.RevokeCampusAuth(ctx, req.ID); err != nil {
		global.GVA_LOG.Error("撤回校园身份审核失败!", zap.Error(err))
		response.FailWithMessage("撤回失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("撤回成功", c)
}
