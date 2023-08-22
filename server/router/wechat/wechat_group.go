package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type WechatGroupRouter struct {
}

// InitWechatGroupRouter 初始化 WechatGroup 路由信息
func (s *WechatGroupRouter) InitWechatGroupRouter(Router *gin.RouterGroup) {
	wechatGroupRouter := Router.Group("wechatGroup").Use(middleware.OperationRecord())
	wechatGroupRouterWithoutRecord := Router.Group("wechatGroup")
	var wechatGroupApi = v1.ApiGroupApp.WechatApiGroup.WechatGroupApi
	{
		wechatGroupRouter.POST("createWechatGroup", wechatGroupApi.CreateWechatGroup)   // 新建WechatGroup
		wechatGroupRouter.DELETE("deleteWechatGroup", wechatGroupApi.DeleteWechatGroup) // 删除WechatGroup
		wechatGroupRouter.DELETE("deleteWechatGroupByIds", wechatGroupApi.DeleteWechatGroupByIds) // 批量删除WechatGroup
		wechatGroupRouter.PUT("updateWechatGroup", wechatGroupApi.UpdateWechatGroup)    // 更新WechatGroup
	}
	{
		wechatGroupRouterWithoutRecord.GET("findWechatGroup", wechatGroupApi.FindWechatGroup)        // 根据ID获取WechatGroup
		wechatGroupRouterWithoutRecord.GET("getWechatGroupList", wechatGroupApi.GetWechatGroupList)  // 获取WechatGroup列表
	}
}
