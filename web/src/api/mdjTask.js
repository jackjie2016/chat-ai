import service from '@/utils/request'

// @Tags MdjTask
// @Summary 创建MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.MdjTask true "创建MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /mdjTask/createMdjTask [post]
export const createMdjTask = (data) => {
  return service({
    url: '/mdjTask/createMdjTask',
    method: 'post',
    data
  })
}

// @Tags MdjTask
// @Summary 删除MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.MdjTask true "删除MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mdjTask/deleteMdjTask [delete]
export const deleteMdjTask = (data) => {
  return service({
    url: '/mdjTask/deleteMdjTask',
    method: 'delete',
    data
  })
}

// @Tags MdjTask
// @Summary 删除MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /mdjTask/deleteMdjTask [delete]
export const deleteMdjTaskByIds = (data) => {
  return service({
    url: '/mdjTask/deleteMdjTaskByIds',
    method: 'delete',
    data
  })
}

// @Tags MdjTask
// @Summary 更新MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.MdjTask true "更新MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /mdjTask/updateMdjTask [put]
export const updateMdjTask = (data) => {
  return service({
    url: '/mdjTask/updateMdjTask',
    method: 'put',
    data
  })
}

// @Tags MdjTask
// @Summary 用id查询MdjTask
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.MdjTask true "用id查询MdjTask"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /mdjTask/findMdjTask [get]
export const findMdjTask = (params) => {
  return service({
    url: '/mdjTask/findMdjTask',
    method: 'get',
    params
  })
}

// @Tags MdjTask
// @Summary 分页获取MdjTask列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取MdjTask列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /mdjTask/getMdjTaskList [get]
export const getMdjTaskList = (params) => {
  return service({
    url: '/mdjTask/getMdjTaskList',
    method: 'get',
    params
  })
}
