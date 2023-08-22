import service from '@/utils/request'

// @Tags WechatGroup
// @Summary 创建WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatGroup true "创建WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatGroup/createWechatGroup [post]
export const createWechatGroup = (data) => {
  return service({
    url: '/wechatGroup/createWechatGroup',
    method: 'post',
    data
  })
}

// @Tags WechatGroup
// @Summary 删除WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatGroup true "删除WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /wechatGroup/deleteWechatGroup [delete]
export const deleteWechatGroup = (data) => {
  return service({
    url: '/wechatGroup/deleteWechatGroup',
    method: 'delete',
    data
  })
}

// @Tags WechatGroup
// @Summary 删除WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /wechatGroup/deleteWechatGroup [delete]
export const deleteWechatGroupByIds = (data) => {
  return service({
    url: '/wechatGroup/deleteWechatGroupByIds',
    method: 'delete',
    data
  })
}

// @Tags WechatGroup
// @Summary 更新WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatGroup true "更新WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /wechatGroup/updateWechatGroup [put]
export const updateWechatGroup = (data) => {
  return service({
    url: '/wechatGroup/updateWechatGroup',
    method: 'put',
    data
  })
}

// @Tags WechatGroup
// @Summary 用id查询WechatGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.WechatGroup true "用id查询WechatGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /wechatGroup/findWechatGroup [get]
export const findWechatGroup = (params) => {
  return service({
    url: '/wechatGroup/findWechatGroup',
    method: 'get',
    params
  })
}

// @Tags WechatGroup
// @Summary 分页获取WechatGroup列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取WechatGroup列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatGroup/getWechatGroupList [get]
export const getWechatGroupList = (params) => {
  return service({
    url: '/wechatGroup/getWechatGroupList',
    method: 'get',
    params
  })
}
