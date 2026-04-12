package campus

import "github.com/gin-gonic/gin"

type CampusOrderRouter struct{}

func (r *CampusOrderRouter) InitCampusOrderRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("campusOrder")
	{
		orderRouter.GET("getCampusOrderList", campusOrderApi.GetCampusOrderList)
		orderRouter.GET("findCampusOrder", campusOrderApi.FindCampusOrder)
	}
}
