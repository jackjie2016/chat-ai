package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type GroupRouter struct {
}

// InitGroupRouter 初始化 Group 路由信息
func (s *GroupRouter) InitGroupRouter(Router *gin.RouterGroup) {
	groupRouter := Router.Group("group").Use(middleware.OperationRecord())
	groupRouterWithoutRecord := Router.Group("group")
	var groupApi = v1.ApiGroupApp.WechatApiGroup.GroupApi
	{
		groupRouter.POST("createGroup", groupApi.CreateGroup)   // 新建Group
		groupRouter.DELETE("deleteGroup", groupApi.DeleteGroup) // 删除Group
		groupRouter.DELETE("deleteGroupByIds", groupApi.DeleteGroupByIds) // 批量删除Group
		groupRouter.PUT("updateGroup", groupApi.UpdateGroup)    // 更新Group
	}
	{
		groupRouterWithoutRecord.GET("findGroup", groupApi.FindGroup)        // 根据ID获取Group
		groupRouterWithoutRecord.GET("getGroupList", groupApi.GetGroupList)  // 获取Group列表
	}
}
