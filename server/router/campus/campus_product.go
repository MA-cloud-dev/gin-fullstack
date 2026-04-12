package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CampusProductRouter struct{}

func (r *CampusProductRouter) InitCampusProductRouter(Router *gin.RouterGroup) {
	productRouter := Router.Group("campusProduct").Use(middleware.OperationRecord())
	productRouterWithoutRecord := Router.Group("campusProduct")
	{
		productRouter.POST("updateCampusProductStatus", campusProductApi.UpdateCampusProductStatus)
	}
	{
		productRouterWithoutRecord.GET("getCampusProductList", campusProductApi.GetCampusProductList)
		productRouterWithoutRecord.GET("findCampusProduct", campusProductApi.FindCampusProduct)
	}
}
