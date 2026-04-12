package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CampusUserRouter struct{}

func (r *CampusUserRouter) InitCampusUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("campusUser").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("campusUser")
	{
		userRouter.POST("updateCampusUserStatus", campusUserApi.UpdateCampusUserStatus)
	}
	{
		userRouterWithoutRecord.GET("getCampusUserList", campusUserApi.GetCampusUserList)
		userRouterWithoutRecord.GET("findCampusUser", campusUserApi.FindCampusUser)
	}
}
