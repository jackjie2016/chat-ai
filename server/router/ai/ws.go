package ai

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type WsRouter struct {
}

// InitMdjTaskRouter 初始化 MdjTask 路由信息
func (s *MdjTaskRouter) InitWsRouter(Router *gin.RouterGroup) {
	wsRouter := Router.Group("ws")

	var wsApi = v1.ApiGroupApp.AiApiGroup.WsApi

	{

		wsRouter.GET("talk", wsApi.Talk) // 获取MdjTask列表
	}
}
