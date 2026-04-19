package campus

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CampusAuthRouter struct{}

func (r *CampusAuthRouter) InitCampusAuthRouter(privateRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	campusAuthRouter := privateRouter.Group("campusAuth").Use(middleware.OperationRecord())
	campusAuthRouterWithoutRecord := privateRouter.Group("campusAuth")
	campusAuthTestRouter := privateRouter.Group("campusAgentReviewTest").Use(middleware.OperationRecord())
	agentReviewPublicRouter := publicRouter.Group("api/agent/review")
	{
		campusAuthRouter.POST("reviewCampusAuth", campusAuthApi.ReviewCampusAuth)
		campusAuthRouter.POST("rejectCampusAuth", campusAuthApi.RejectCampusAuth)
		campusAuthRouter.POST("revokeCampusAuth", campusAuthApi.RevokeCampusAuth)
	}
	{
		campusAuthRouterWithoutRecord.GET("getCampusAuthList", campusAuthApi.GetCampusAuthList)
		campusAuthRouterWithoutRecord.GET("findCampusAuth", campusAuthApi.FindCampusAuth)
	}
	{
		campusAuthTestRouter.POST("submit", campusAuthApi.SubmitCampusAuthTest)
		agentReviewPublicRouter.POST("callback", campusAuthApi.AgentReviewCallback)
	}
}
