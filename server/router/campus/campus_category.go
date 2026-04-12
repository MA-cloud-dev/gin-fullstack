package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CampusCategoryRouter struct{}

func (r *CampusCategoryRouter) InitCampusCategoryRouter(Router *gin.RouterGroup) {
	categoryRouter := Router.Group("campusCategory").Use(middleware.OperationRecord())
	categoryRouterWithoutRecord := Router.Group("campusCategory")
	{
		categoryRouter.POST("createCampusCategory", campusCategoryApi.CreateCampusCategory)
		categoryRouter.PUT("updateCampusCategory", campusCategoryApi.UpdateCampusCategory)
		categoryRouter.POST("updateCampusCategoryStatus", campusCategoryApi.UpdateCampusCategoryStatus)
	}
	{
		categoryRouterWithoutRecord.GET("getCampusCategoryList", campusCategoryApi.GetCampusCategoryList)
		categoryRouterWithoutRecord.GET("findCampusCategory", campusCategoryApi.FindCampusCategory)
	}
}
