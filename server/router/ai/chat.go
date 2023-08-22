package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct{}

func (s *ChatRouter) InitChatApiRouter(Router *gin.RouterGroup) {
	chatRouter := Router.Group("chat").Use(middleware.OperationRecord())
	chatRouterWithoutRecord := Router.Group("chat")
	chatRouterApi := v1.ApiGroupApp.AiGroup.ChatApi
	{
		chatRouter.POST("createApi", chatRouterApi.CreateApi) // 创建Api

	}
	{
		chatRouterWithoutRecord.POST("completion", chatRouterApi.Completion) // 获取所有api
		chatRouterWithoutRecord.POST("context", chatRouterApi.Context)       // 获取所有api
		chatRouterWithoutRecord.POST("reply", chatRouterApi.Reply)           // 获取所有api
		chatRouterWithoutRecord.GET("login", chatRouterApi.Login)            // 登录
		chatRouterWithoutRecord.GET("wechatWs", chatRouterApi.WechatWs)      // 登录
		chatRouterWithoutRecord.GET("sse", chatRouterApi.Sse)                // 获取所有api

	}
}
