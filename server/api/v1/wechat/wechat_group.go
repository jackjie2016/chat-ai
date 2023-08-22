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

type WechatGroupApi struct {
}

var wechatGroupService = service.ServiceGroupApp.WechatServiceGroup.WechatGroupService


// CreateWechatGroup 创建WechatGroup
// @Tags WechatGroup
// @Summary 创建WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.WechatGroup true "创建WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatGroup/createWechatGroup [post]
func (wechatGroupApi *WechatGroupApi) CreateWechatGroup(c *gin.Context) {
	var wechatGroup wechat.WechatGroup
	err := c.ShouldBindJSON(&wechatGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := wechatGroupService.CreateWechatGroup(&wechatGroup); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteWechatGroup 删除WechatGroup
// @Tags WechatGroup
// @Summary 删除WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.WechatGroup true "删除WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /wechatGroup/deleteWechatGroup [delete]
func (wechatGroupApi *WechatGroupApi) DeleteWechatGroup(c *gin.Context) {
	var wechatGroup wechat.WechatGroup
	err := c.ShouldBindJSON(&wechatGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := wechatGroupService.DeleteWechatGroup(wechatGroup); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteWechatGroupByIds 批量删除WechatGroup
// @Tags WechatGroup
// @Summary 批量删除WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /wechatGroup/deleteWechatGroupByIds [delete]
func (wechatGroupApi *WechatGroupApi) DeleteWechatGroupByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := wechatGroupService.DeleteWechatGroupByIds(IDS); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateWechatGroup 更新WechatGroup
// @Tags WechatGroup
// @Summary 更新WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.WechatGroup true "更新WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /wechatGroup/updateWechatGroup [put]
func (wechatGroupApi *WechatGroupApi) UpdateWechatGroup(c *gin.Context) {
	var wechatGroup wechat.WechatGroup
	err := c.ShouldBindJSON(&wechatGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := wechatGroupService.UpdateWechatGroup(wechatGroup); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindWechatGroup 用id查询WechatGroup
// @Tags WechatGroup
// @Summary 用id查询WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query wechat.WechatGroup true "用id查询WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /wechatGroup/findWechatGroup [get]
func (wechatGroupApi *WechatGroupApi) FindWechatGroup(c *gin.Context) {
	var wechatGroup wechat.WechatGroup
	err := c.ShouldBindQuery(&wechatGroup)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rewechatGroup, err := wechatGroupService.GetWechatGroup(wechatGroup.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rewechatGroup": rewechatGroup}, c)
	}
}

// GetWechatGroupList 分页获取WechatGroup列表
// @Tags WechatGroup
// @Summary 分页获取WechatGroup列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query wechatReq.WechatGroupSearch true "分页获取WechatGroup列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatGroup/getWechatGroupList [get]
func (wechatGroupApi *WechatGroupApi) GetWechatGroupList(c *gin.Context) {
	var pageInfo wechatReq.WechatGroupSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := wechatGroupService.GetWechatGroupInfoList(pageInfo); err != nil {
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
