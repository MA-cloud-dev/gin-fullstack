package campus

import "github.com/gin-gonic/gin"

type CampusOverviewRouter struct{}

func (r *CampusOverviewRouter) InitCampusOverviewRouter(Router *gin.RouterGroup) {
	campusOverviewRouter := Router.Group("campusOverview")
	{
		campusOverviewRouter.GET("getCampusOverview", campusOverviewApi.GetCampusOverview)
	}
}
