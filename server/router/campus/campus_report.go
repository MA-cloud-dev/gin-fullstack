package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CampusReportRouter struct{}

func (r *CampusReportRouter) InitCampusReportRouter(Router *gin.RouterGroup) {
	reportRouter := Router.Group("campusReport").Use(middleware.OperationRecord())
	reportRouterWithoutRecord := Router.Group("campusReport")
	{
		reportRouter.POST("handleCampusReport", campusReportApi.HandleCampusReport)
	}
	{
		reportRouterWithoutRecord.GET("getCampusReportList", campusReportApi.GetCampusReportList)
		reportRouterWithoutRecord.GET("findCampusReport", campusReportApi.FindCampusReport)
	}
}
