package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CampusAdminStaffRouter struct{}

func (r *CampusAdminStaffRouter) InitCampusAdminStaffRouter(Router *gin.RouterGroup) {
	adminStaffRouter := Router.Group("campusAdminStaff").Use(middleware.OperationRecord())
	adminStaffRouterWithoutRecord := Router.Group("campusAdminStaff")
	{
		adminStaffRouter.POST("updateCampusAdminStaffStatus", campusAdminStaffApi.UpdateCampusAdminStaffStatus)
	}
	{
		adminStaffRouterWithoutRecord.GET("getCampusAdminStaffList", campusAdminStaffApi.GetCampusAdminStaffList)
		adminStaffRouterWithoutRecord.GET("findCampusAdminStaff", campusAdminStaffApi.FindCampusAdminStaff)
	}
}
