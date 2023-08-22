package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/ai"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type MdjTaskApi struct {
}

var mdjTaskService = service.ServiceGroupApp.AiServiceGroup.MdjTaskService


// CreateMdjTask 创建MdjTask
// @Tags MdjTask
// @Summary 创建MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ai.MdjTask true "创建MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /mdjTask/createMdjTask [post]
func (mdjTaskApi *MdjTaskApi) CreateMdjTask(c *gin.Context) {
	var mdjTask ai.MdjTask
	err := c.ShouldBindJSON(&mdjTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := mdjTaskService.CreateMdjTask(&mdjTask); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteMdjTask 删除MdjTask
// @Tags MdjTask
// @Summary 删除MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ai.MdjTask true "删除MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mdjTask/deleteMdjTask [delete]
func (mdjTaskApi *MdjTaskApi) DeleteMdjTask(c *gin.Context) {
	var mdjTask ai.MdjTask
	err := c.ShouldBindJSON(&mdjTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := mdjTaskService.DeleteMdjTask(mdjTask); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteMdjTaskByIds 批量删除MdjTask
// @Tags MdjTask
// @Summary 批量删除MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /mdjTask/deleteMdjTaskByIds [delete]
func (mdjTaskApi *MdjTaskApi) DeleteMdjTaskByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := mdjTaskService.DeleteMdjTaskByIds(IDS); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateMdjTask 更新MdjTask
// @Tags MdjTask
// @Summary 更新MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ai.MdjTask true "更新MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /mdjTask/updateMdjTask [put]
func (mdjTaskApi *MdjTaskApi) UpdateMdjTask(c *gin.Context) {
	var mdjTask ai.MdjTask
	err := c.ShouldBindJSON(&mdjTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := mdjTaskService.UpdateMdjTask(mdjTask); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindMdjTask 用id查询MdjTask
// @Tags MdjTask
// @Summary 用id查询MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query ai.MdjTask true "用id查询MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /mdjTask/findMdjTask [get]
func (mdjTaskApi *MdjTaskApi) FindMdjTask(c *gin.Context) {
	var mdjTask ai.MdjTask
	err := c.ShouldBindQuery(&mdjTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if remdjTask, err := mdjTaskService.GetMdjTask(mdjTask.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"remdjTask": remdjTask}, c)
	}
}

// GetMdjTaskList 分页获取MdjTask列表
// @Tags MdjTask
// @Summary 分页获取MdjTask列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query aiReq.MdjTaskSearch true "分页获取MdjTask列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /mdjTask/getMdjTaskList [get]
func (mdjTaskApi *MdjTaskApi) GetMdjTaskList(c *gin.Context) {
	var pageInfo aiReq.MdjTaskSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := mdjTaskService.GetMdjTaskInfoList(pageInfo); err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败", c)
    } else {
        response.OkWithDetailed(response.PageResult{
            List:     list,
            Total:    total,
            Page:     pageInfo.Page,
            PageSize: pageInfo.PageSize,
        }, "获取成功", c)
    }
}
