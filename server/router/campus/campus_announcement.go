package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CampusAnnouncementRouter struct{}

func (r *CampusAnnouncementRouter) InitCampusAnnouncementRouter(Router *gin.RouterGroup) {
	announcementRouter := Router.Group("campusAnnouncement").Use(middleware.OperationRecord())
	announcementRouterWithoutRecord := Router.Group("campusAnnouncement")
	{
		announcementRouter.POST("createCampusAnnouncement", campusAnnouncementApi.CreateCampusAnnouncement)
		announcementRouter.PUT("updateCampusAnnouncement", campusAnnouncementApi.UpdateCampusAnnouncement)
		announcementRouter.DELETE("deleteCampusAnnouncement", campusAnnouncementApi.DeleteCampusAnnouncement)
		announcementRouter.POST("updateCampusAnnouncementStatus", campusAnnouncementApi.UpdateCampusAnnouncementStatus)
	}
	{
		announcementRouterWithoutRecord.GET("getCampusAnnouncementList", campusAnnouncementApi.GetCampusAnnouncementList)
		announcementRouterWithoutRecord.GET("findCampusAnnouncement", campusAnnouncementApi.FindCampusAnnouncement)
	}
}
