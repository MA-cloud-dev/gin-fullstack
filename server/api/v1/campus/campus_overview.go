package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CampusOverviewApi struct{}

func (a *CampusOverviewApi) GetCampusOverview(c *gin.Context) {
	data, err := campusOverviewService.GetCampusOverview(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("获取校园数据总揽失败!", zap.Error(err))
		response.FailWithMessage("获取校园数据总揽失败:"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}
