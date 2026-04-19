package campus

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	CampusOverviewApi
	CampusAuthApi
	CampusProductApi
	CampusCategoryApi
	CampusUserApi
	CampusAdminStaffApi
	CampusAnnouncementApi
	CampusOrderApi
	CampusReportApi
	CampusOperationLogApi
}

var (
	campusOverviewService     = service.ServiceGroupApp.CampusServiceGroup.CampusOverviewService
	campusAuthService         = service.ServiceGroupApp.CampusServiceGroup.CampusAuthService
	campusProductService      = service.ServiceGroupApp.CampusServiceGroup.CampusProductService
	campusCategoryService     = service.ServiceGroupApp.CampusServiceGroup.CampusCategoryService
	campusUserService         = service.ServiceGroupApp.CampusServiceGroup.CampusUserService
	campusAdminStaffService   = service.ServiceGroupApp.CampusServiceGroup.CampusAdminStaffService
	campusAnnouncementService = service.ServiceGroupApp.CampusServiceGroup.CampusAnnouncementService
	campusOrderService        = service.ServiceGroupApp.CampusServiceGroup.CampusOrderService
	campusReportService       = service.ServiceGroupApp.CampusServiceGroup.CampusReportService
	campusOperationLogService = service.ServiceGroupApp.CampusServiceGroup.CampusOperationLogService
)
