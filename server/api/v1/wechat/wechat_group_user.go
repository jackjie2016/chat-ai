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

type WechatGroupUserApi struct {
}

var wechatGroupUserService = service.ServiceGroupApp.WechatServiceGroup.WechatGroupUserService


// CreateWechatGroupUser 创建WechatGroupUser
// @Tags WechatGroupUser
// @Summary 创建WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.WechatGroupUser true "创建WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatGroupUser/createWechatGroupUser [post]
func (wechatGroupUserApi *WechatGroupUserApi) CreateWechatGroupUser(c *gin.Context) {
	var wechatGroupUser wechat.WechatGroupUser
	err := c.ShouldBindJSON(&wechatGroupUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := wechatGroupUserService.CreateWechatGroupUser(&wechatGroupUser); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteWechatGroupUser 删除WechatGroupUser
// @Tags WechatGroupUser
// @Summary 删除WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.WechatGroupUser true "删除WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /wechatGroupUser/deleteWechatGroupUser [delete]
func (wechatGroupUserApi *WechatGroupUserApi) DeleteWechatGroupUser(c *gin.Context) {
	var wechatGroupUser wechat.WechatGroupUser
	err := c.ShouldBindJSON(&wechatGroupUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := wechatGroupUserService.DeleteWechatGroupUser(wechatGroupUser); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteWechatGroupUserByIds 批量删除WechatGroupUser
// @Tags WechatGroupUser
// @Summary 批量删除WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /wechatGroupUser/deleteWechatGroupUserByIds [delete]
func (wechatGroupUserApi *WechatGroupUserApi) DeleteWechatGroupUserByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := wechatGroupUserService.DeleteWechatGroupUserByIds(IDS); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateWechatGroupUser 更新WechatGroupUser
// @Tags WechatGroupUser
// @Summary 更新WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body wechat.WechatGroupUser true "更新WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /wechatGroupUser/updateWechatGroupUser [put]
func (wechatGroupUserApi *WechatGroupUserApi) UpdateWechatGroupUser(c *gin.Context) {
	var wechatGroupUser wechat.WechatGroupUser
	err := c.ShouldBindJSON(&wechatGroupUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := wechatGroupUserService.UpdateWechatGroupUser(wechatGroupUser); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindWechatGroupUser 用id查询WechatGroupUser
// @Tags WechatGroupUser
// @Summary 用id查询WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query wechat.WechatGroupUser true "用id查询WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /wechatGroupUser/findWechatGroupUser [get]
func (wechatGroupUserApi *WechatGroupUserApi) FindWechatGroupUser(c *gin.Context) {
	var wechatGroupUser wechat.WechatGroupUser
	err := c.ShouldBindQuery(&wechatGroupUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rewechatGroupUser, err := wechatGroupUserService.GetWechatGroupUser(wechatGroupUser.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rewechatGroupUser": rewechatGroupUser}, c)
	}
}

// GetWechatGroupUserList 分页获取WechatGroupUser列表
// @Tags WechatGroupUser
// @Summary 分页获取WechatGroupUser列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query wechatReq.WechatGroupUserSearch true "分页获取WechatGroupUser列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatGroupUser/getWechatGroupUserList [get]
func (wechatGroupUserApi *WechatGroupUserApi) GetWechatGroupUserList(c *gin.Context) {
	var pageInfo wechatReq.WechatGroupUserSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := wechatGroupUserService.GetWechatGroupUserInfoList(pageInfo); err != nil {
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
