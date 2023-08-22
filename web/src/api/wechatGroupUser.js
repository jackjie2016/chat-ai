import service from '@/utils/request'

// @Tags WechatGroupUser
// @Summary 创建WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatGroupUser true "创建WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatGroupUser/createWechatGroupUser [post]
export const createWechatGroupUser = (data) => {
  return service({
    url: '/wechatGroupUser/createWechatGroupUser',
    method: 'post',
    data
  })
}

// @Tags WechatGroupUser
// @Summary 删除WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatGroupUser true "删除WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /wechatGroupUser/deleteWechatGroupUser [delete]
export const deleteWechatGroupUser = (data) => {
  return service({
    url: '/wechatGroupUser/deleteWechatGroupUser',
    method: 'delete',
    data
  })
}

// @Tags WechatGroupUser
// @Summary 删除WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /wechatGroupUser/deleteWechatGroupUser [delete]
export const deleteWechatGroupUserByIds = (data) => {
  return service({
    url: '/wechatGroupUser/deleteWechatGroupUserByIds',
    method: 'delete',
    data
  })
}

// @Tags WechatGroupUser
// @Summary 更新WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatGroupUser true "更新WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /wechatGroupUser/updateWechatGroupUser [put]
export const updateWechatGroupUser = (data) => {
  return service({
    url: '/wechatGroupUser/updateWechatGroupUser',
    method: 'put',
    data
  })
}

// @Tags WechatGroupUser
// @Summary 用id查询WechatGroupUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.WechatGroupUser true "用id查询WechatGroupUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /wechatGroupUser/findWechatGroupUser [get]
export const findWechatGroupUser = (params) => {
  return service({
    url: '/wechatGroupUser/findWechatGroupUser',
    method: 'get',
    params
  })
}

// @Tags WechatGroupUser
// @Summary 分页获取WechatGroupUser列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取WechatGroupUser列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatGroupUser/getWechatGroupUserList [get]
export const getWechatGroupUserList = (params) => {
  return service({
    url: '/wechatGroupUser/getWechatGroupUserList',
    method: 'get',
    params
  })
}
