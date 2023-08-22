package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MdjRouter struct{}

func (s *MdjRouter) InitMdjApiRouter(Router *gin.RouterGroup) {
	mdjRouter := Router.Group("mdj").Use(middleware.OperationRecord())
	mdjRouterWithoutRecord := Router.Group("mdj")
	mdjRouterApi := v1.ApiGroupApp.AiGroup.MdjApi
	{
		mdjRouter.POST("createApi", mdjRouterApi.CreateApi) // 创建Api

	}
	{

		mdjRouterWithoutRecord.POST("/retrieveMessages", mdjRouterApi.RetrieveMessages)
		mdjRouterWithoutRecord.POST("/msglist", mdjRouterApi.Msglist)
		mdjRouterWithoutRecord.POST("/imagine", mdjRouterApi.Imagine)
		mdjRouterWithoutRecord.POST("/upscale", mdjRouterApi.Upscale)
		mdjRouterWithoutRecord.POST("/variation", mdjRouterApi.Variation)
	}
}
