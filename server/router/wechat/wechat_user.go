package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type WechatUserRouter struct {
}

// InitWechatUserRouter 初始化 WechatUser 路由信息
func (s *WechatUserRouter) InitWechatUserRouter(Router *gin.RouterGroup) {
	wechatUserRouter := Router.Group("wechatUser").Use(middleware.OperationRecord())
	wechatUserRouterWithoutRecord := Router.Group("wechatUser")
	var wechatUserApi = v1.ApiGroupApp.WechatApiGroup.WechatUserApi
	{
		wechatUserRouter.POST("createWechatUser", wechatUserApi.CreateWechatUser)   // 新建WechatUser
		wechatUserRouter.DELETE("deleteWechatUser", wechatUserApi.DeleteWechatUser) // 删除WechatUser
		wechatUserRouter.DELETE("deleteWechatUserByIds", wechatUserApi.DeleteWechatUserByIds) // 批量删除WechatUser
		wechatUserRouter.PUT("updateWechatUser", wechatUserApi.UpdateWechatUser)    // 更新WechatUser
	}
	{
		wechatUserRouterWithoutRecord.GET("findWechatUser", wechatUserApi.FindWechatUser)        // 根据ID获取WechatUser
		wechatUserRouterWithoutRecord.GET("getWechatUserList", wechatUserApi.GetWechatUserList)  // 获取WechatUser列表
	}
}
