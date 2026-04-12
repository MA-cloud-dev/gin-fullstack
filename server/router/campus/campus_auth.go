package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CampusAuthRouter struct{}

func (r *CampusAuthRouter) InitCampusAuthRouter(Router *gin.RouterGroup) {
	campusAuthRouter := Router.Group("campusAuth").Use(middleware.OperationRecord())
	campusAuthRouterWithoutRecord := Router.Group("campusAuth")
	{
		campusAuthRouter.POST("reviewCampusAuth", campusAuthApi.ReviewCampusAuth)
		campusAuthRouter.POST("revokeCampusAuth", campusAuthApi.RevokeCampusAuth)
	}
	{
		campusAuthRouterWithoutRecord.GET("getCampusAuthList", campusAuthApi.GetCampusAuthList)
		campusAuthRouterWithoutRecord.GET("findCampusAuth", campusAuthApi.FindCampusAuth)
	}
}
