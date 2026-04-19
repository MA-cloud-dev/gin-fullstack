package campus

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	CampusOverviewRouter
	CampusAuthRouter
	CampusProductRouter
	CampusCategoryRouter
	CampusUserRouter
	CampusAdminStaffRouter
	CampusAnnouncementRouter
	CampusOrderRouter
	CampusReportRouter
	CampusOperationLogRouter
}

var (
	campusOverviewApi     = api.ApiGroupApp.CampusApiGroup.CampusOverviewApi
	campusAuthApi         = api.ApiGroupApp.CampusApiGroup.CampusAuthApi
	campusProductApi      = api.ApiGroupApp.CampusApiGroup.CampusProductApi
	campusCategoryApi     = api.ApiGroupApp.CampusApiGroup.CampusCategoryApi
	campusUserApi         = api.ApiGroupApp.CampusApiGroup.CampusUserApi
	campusAdminStaffApi   = api.ApiGroupApp.CampusApiGroup.CampusAdminStaffApi
	campusAnnouncementApi = api.ApiGroupApp.CampusApiGroup.CampusAnnouncementApi
	campusOrderApi        = api.ApiGroupApp.CampusApiGroup.CampusOrderApi
	campusReportApi       = api.ApiGroupApp.CampusApiGroup.CampusReportApi
	campusOperationLogApi = api.ApiGroupApp.CampusApiGroup.CampusOperationLogApi
)
