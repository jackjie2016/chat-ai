package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MidjourneyRouter struct{}

func (s *MdjRouter) InitMidjourneyApiRouter(Router *gin.RouterGroup) {
	MidjourneyRouter := Router.Group("trigger").Use(middleware.OperationRecord())
	MidjourneyRouterWithoutRecord := Router.Group("trigger")
	mdjRouterApi := v1.ApiGroupApp.AiGroup.MidBotApi
	{
		MidjourneyRouter.POST("/health", mdjRouterApi.Health)
	}
	{

		MidjourneyRouterWithoutRecord.POST("/upload", mdjRouterApi.Upload)
		MidjourneyRouterWithoutRecord.POST("/midjourney-bot", mdjRouterApi.MidjourneyBot)
		MidjourneyRouterWithoutRecord.GET("/messages", mdjRouterApi.Messages)
		MidjourneyRouterWithoutRecord.POST("/midjourney-list", mdjRouterApi.MidjourneyList)
		MidjourneyRouterWithoutRecord.POST("/list", mdjRouterApi.List)
		MidjourneyRouterWithoutRecord.POST("/find", mdjRouterApi.Find)
		MidjourneyRouterWithoutRecord.POST("/ok", mdjRouterApi.Ok)

	}
}
