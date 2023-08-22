package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type WechatGroupUserRouter struct {
}

// InitWechatGroupUserRouter 初始化 WechatGroupUser 路由信息
func (s *WechatGroupUserRouter) InitWechatGroupUserRouter(Router *gin.RouterGroup) {
	wechatGroupUserRouter := Router.Group("wechatGroupUser").Use(middleware.OperationRecord())
	wechatGroupUserRouterWithoutRecord := Router.Group("wechatGroupUser")
	var wechatGroupUserApi = v1.ApiGroupApp.WechatApiGroup.WechatGroupUserApi
	{
		wechatGroupUserRouter.POST("createWechatGroupUser", wechatGroupUserApi.CreateWechatGroupUser)   // 新建WechatGroupUser
		wechatGroupUserRouter.DELETE("deleteWechatGroupUser", wechatGroupUserApi.DeleteWechatGroupUser) // 删除WechatGroupUser
		wechatGroupUserRouter.DELETE("deleteWechatGroupUserByIds", wechatGroupUserApi.DeleteWechatGroupUserByIds) // 批量删除WechatGroupUser
		wechatGroupUserRouter.PUT("updateWechatGroupUser", wechatGroupUserApi.UpdateWechatGroupUser)    // 更新WechatGroupUser
	}
	{
		wechatGroupUserRouterWithoutRecord.GET("findWechatGroupUser", wechatGroupUserApi.FindWechatGroupUser)        // 根据ID获取WechatGroupUser
		wechatGroupUserRouterWithoutRecord.GET("getWechatGroupUserList", wechatGroupUserApi.GetWechatGroupUserList)  // 获取WechatGroupUser列表
	}
}
