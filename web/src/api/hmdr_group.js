import service from '@/utils/request'

// @Tags HmdrGroup
// @Summary 创建HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrGroup true "创建HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hmdrGroup/createHmdrGroup [post]
export const createHmdrGroup = (data) => {
  return service({
    url: '/hmdrGroup/createHmdrGroup',
    method: 'post',
    data
  })
}

// @Tags HmdrGroup
// @Summary 删除HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrGroup true "删除HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /hmdrGroup/deleteHmdrGroup [delete]
export const deleteHmdrGroup = (data) => {
  return service({
    url: '/hmdrGroup/deleteHmdrGroup',
    method: 'delete',
    data
  })
}

// @Tags HmdrGroup
// @Summary 删除HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /hmdrGroup/deleteHmdrGroup [delete]
export const deleteHmdrGroupByIds = (data) => {
  return service({
    url: '/hmdrGroup/deleteHmdrGroupByIds',
    method: 'delete',
    data
  })
}

// @Tags HmdrGroup
// @Summary 更新HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrGroup true "更新HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /hmdrGroup/updateHmdrGroup [put]
export const updateHmdrGroup = (data) => {
  return service({
    url: '/hmdrGroup/updateHmdrGroup',
    method: 'put',
    data
  })
}

// @Tags HmdrGroup
// @Summary 用id查询HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrGroup true "用id查询HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /hmdrGroup/findHmdrGroup [get]
export const findHmdrGroup = (params) => {
  return service({
    url: '/hmdrGroup/findHmdrGroup',
    method: 'get',
    params
  })
}

// @Tags HmdrGroup
// @Summary 分页获取HmdrGroup列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取HmdrGroup列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hmdrGroup/getHmdrGroupList [get]
export const getHmdrGroupList = (params) => {
  return service({
    url: '/hmdrGroup/getHmdrGroupList',
    method: 'get',
    params
  })
}
