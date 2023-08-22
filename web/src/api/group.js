import service from '@/utils/request'

// @Tags Group
// @Summary 创建Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Group true "创建Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /group/createGroup [post]
export const createGroup = (data) => {
  return service({
    url: '/group/createGroup',
    method: 'post',
    data
  })
}

// @Tags Group
// @Summary 删除Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Group true "删除Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /group/deleteGroup [delete]
export const deleteGroup = (data) => {
  return service({
    url: '/group/deleteGroup',
    method: 'delete',
    data
  })
}

// @Tags Group
// @Summary 删除Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /group/deleteGroup [delete]
export const deleteGroupByIds = (data) => {
  return service({
    url: '/group/deleteGroupByIds',
    method: 'delete',
    data
  })
}

// @Tags Group
// @Summary 更新Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Group true "更新Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /group/updateGroup [put]
export const updateGroup = (data) => {
  return service({
    url: '/group/updateGroup',
    method: 'put',
    data
  })
}

// @Tags Group
// @Summary 用id查询Group
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Group true "用id查询Group"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /group/findGroup [get]
export const findGroup = (params) => {
  return service({
    url: '/group/findGroup',
    method: 'get',
    params
  })
}

// @Tags Group
// @Summary 分页获取Group列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Group列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /group/getGroupList [get]
export const getGroupList = (params) => {
  return service({
    url: '/group/getGroupList',
    method: 'get',
    params
  })
}
