package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    wechatReq "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type GroupApi struct {
}

var groupService = service.ServiceGroupApp.WechatServiceGroup.GroupService


// CreateGroup 创建Group
// @Tags Group
// @Summary 创建Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.Group true "创建Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /group/createGroup [post]
func (groupApi *GroupApi) CreateGroup(c *gin.Context) {
	var group wechat.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := groupService.CreateGroup(&group); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteGroup 删除Group
// @Tags Group
// @Summary 删除Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.Group true "删除Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /group/deleteGroup [delete]
func (groupApi *GroupApi) DeleteGroup(c *gin.Context) {
	var group wechat.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := groupService.DeleteGroup(group); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteGroupByIds 批量删除Group
// @Tags Group
// @Summary 批量删除Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /group/deleteGroupByIds [delete]
func (groupApi *GroupApi) DeleteGroupByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := groupService.DeleteGroupByIds(IDS); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateGroup 更新Group
// @Tags Group
// @Summary 更新Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.Group true "更新Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /group/updateGroup [put]
func (groupApi *GroupApi) UpdateGroup(c *gin.Context) {
	var group wechat.Group
	err := c.ShouldBindJSON(&group)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := groupService.UpdateGroup(group); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindGroup 用id查询Group
// @Tags Group
// @Summary 用id查询Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query wechat.Group true "用id查询Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /group/findGroup [get]
func (groupApi *GroupApi) FindGroup(c *gin.Context) {
	var group wechat.Group
	err := c.ShouldBindQuery(&group)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if regroup, err := groupService.GetGroup(group.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"regroup": regroup}, c)
	}
}

// GetGroupList 分页获取Group列表
// @Tags Group
// @Summary 分页获取Group列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query wechatReq.GroupSearch true "分页获取Group列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /group/getGroupList [get]
func (groupApi *GroupApi) GetGroupList(c *gin.Context) {
	var pageInfo wechatReq.GroupSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := groupService.GetGroupInfoList(pageInfo); err != nil {
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
