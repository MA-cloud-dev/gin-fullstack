package campus

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	CampusAuthRouter
	CampusProductRouter
	CampusCategoryRouter
	CampusUserRouter
	CampusAdminStaffRouter
	CampusAnnouncementRouter
	CampusOrderRouter
	CampusReportRouter
}

var (
	campusAuthApi         = api.ApiGroupApp.CampusApiGroup.CampusAuthApi
	campusProductApi      = api.ApiGroupApp.CampusApiGroup.CampusProductApi
	campusCategoryApi     = api.ApiGroupApp.CampusApiGroup.CampusCategoryApi
	campusUserApi         = api.ApiGroupApp.CampusApiGroup.CampusUserApi
	campusAdminStaffApi   = api.ApiGroupApp.CampusApiGroup.CampusAdminStaffApi
	campusAnnouncementApi = api.ApiGroupApp.CampusApiGroup.CampusAnnouncementApi
	campusOrderApi        = api.ApiGroupApp.CampusApiGroup.CampusOrderApi
	campusReportApi       = api.ApiGroupApp.CampusApiGroup.CampusReportApi
)
