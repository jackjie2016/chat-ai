package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type MdjTaskRouter struct {
}

// InitMdjTaskRouter 初始化 MdjTask 路由信息
func (s *MdjTaskRouter) InitMdjTaskRouter(Router *gin.RouterGroup) {
	mdjTaskRouter := Router.Group("mdjTask").Use(middleware.OperationRecord())
	mdjTaskRouterWithoutRecord := Router.Group("mdjTask")
	var mdjTaskApi = v1.ApiGroupApp.AiApiGroup.MdjTaskApi
	{
		mdjTaskRouter.POST("createMdjTask", mdjTaskApi.CreateMdjTask)   // 新建MdjTask
		mdjTaskRouter.DELETE("deleteMdjTask", mdjTaskApi.DeleteMdjTask) // 删除MdjTask
		mdjTaskRouter.DELETE("deleteMdjTaskByIds", mdjTaskApi.DeleteMdjTaskByIds) // 批量删除MdjTask
		mdjTaskRouter.PUT("updateMdjTask", mdjTaskApi.UpdateMdjTask)    // 更新MdjTask
	}
	{
		mdjTaskRouterWithoutRecord.GET("findMdjTask", mdjTaskApi.FindMdjTask)        // 根据ID获取MdjTask
		mdjTaskRouterWithoutRecord.GET("getMdjTaskList", mdjTaskApi.GetMdjTaskList)  // 获取MdjTask列表
	}
}
