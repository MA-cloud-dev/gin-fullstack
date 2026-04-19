package campus

import "github.com/gin-gonic/gin"

type CampusOperationLogRouter struct{}

func (r *CampusOperationLogRouter) InitCampusOperationLogRouter(Router *gin.RouterGroup) {
	campusOperationLogRouter := Router.Group("campusOperationLog")
	{
		campusOperationLogRouter.GET("getCampusOperationLogList", campusOperationLogApi.GetCampusOperationLogList)
		campusOperationLogRouter.GET("findCampusOperationLog", campusOperationLogApi.FindCampusOperationLog)
	}
}
