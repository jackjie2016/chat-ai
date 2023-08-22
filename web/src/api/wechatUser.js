import service from '@/utils/request'

// @Tags WechatUser
// @Summary 创建WechatUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatUser true "创建WechatUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatUser/createWechatUser [post]
export const createWechatUser = (data) => {
  return service({
    url: '/wechatUser/createWechatUser',
    method: 'post',
    data
  })
}

// @Tags WechatUser
// @Summary 删除WechatUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatUser true "删除WechatUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /wechatUser/deleteWechatUser [delete]
export const deleteWechatUser = (data) => {
  return service({
    url: '/wechatUser/deleteWechatUser',
    method: 'delete',
    data
  })
}

// @Tags WechatUser
// @Summary 删除WechatUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除WechatUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /wechatUser/deleteWechatUser [delete]
export const deleteWechatUserByIds = (data) => {
  return service({
    url: '/wechatUser/deleteWechatUserByIds',
    method: 'delete',
    data
  })
}

// @Tags WechatUser
// @Summary 更新WechatUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.WechatUser true "更新WechatUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /wechatUser/updateWechatUser [put]
export const updateWechatUser = (data) => {
  return service({
    url: '/wechatUser/updateWechatUser',
    method: 'put',
    data
  })
}

// @Tags WechatUser
// @Summary 用id查询WechatUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.WechatUser true "用id查询WechatUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /wechatUser/findWechatUser [get]
export const findWechatUser = (params) => {
  return service({
    url: '/wechatUser/findWechatUser',
    method: 'get',
    params
  })
}

// @Tags WechatUser
// @Summary 分页获取WechatUser列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取WechatUser列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /wechatUser/getWechatUserList [get]
export const getWechatUserList = (params) => {
  return service({
    url: '/wechatUser/getWechatUserList',
    method: 'get',
    params
  })
}
