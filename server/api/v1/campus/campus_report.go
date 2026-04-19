package campus

import (
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusReportApi struct{}

func (a *CampusReportApi) GetCampusReportList(c *gin.Context) {
	ctx := c.Request.Context()

	var pageInfo campusReq.CampusReportSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := campusReportService.GetCampusReportInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取举报列表失败!", zap.Error(err))
		response.FailWithMessage("获取举报列表失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *CampusReportApi) FindCampusReport(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("举报ID不合法", c)
		return
	}

	report, err := campusReportService.GetCampusReport(ctx, uint(id))
	if err != nil {
		global.GVA_LOG.Error("获取举报详情失败!", zap.Error(err))
		response.FailWithMessage("获取举报详情失败:"+err.Error(), c)
		return
	}
	response.OkWithData(report, c)
}

func (a *CampusReportApi) HandleCampusReport(c *gin.Context) {
	ctx := c.Request.Context()

	var req campusReq.HandleCampusReportReq
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
	req.HandleResult = strings.TrimSpace(req.HandleResult)
	if req.HandleResult == "" {
		response.FailWithMessage("处理结果不能为空", c)
		return
	}

	if err := campusReportService.HandleCampusReport(ctx, utils.GetUserName(c), req, buildCampusAuditMeta(c)); err != nil {
		global.GVA_LOG.Error("处理举报失败!", zap.Error(err))
		response.FailWithMessage("处理举报失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("处理成功", c)
}
