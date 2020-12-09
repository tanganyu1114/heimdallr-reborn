import service from '@/utils/request'

// @Tags HmdrHost
// @Summary 创建HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrHost true "创建HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hmdrHost/createHmdrHost [post]
export const createHmdrHost = (data) => {
  return service({
    url: '/hmdrHost/createHmdrHost',
    method: 'post',
    data
  })
}

// @Tags HmdrHost
// @Summary 删除HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrHost true "删除HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /hmdrHost/deleteHmdrHost [delete]
export const deleteHmdrHost = (data) => {
  return service({
    url: '/hmdrHost/deleteHmdrHost',
    method: 'delete',
    data
  })
}

// @Tags HmdrHost
// @Summary 删除HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /hmdrHost/deleteHmdrHost [delete]
export const deleteHmdrHostByIds = (data) => {
  return service({
    url: '/hmdrHost/deleteHmdrHostByIds',
    method: 'delete',
    data
  })
}

// @Tags HmdrHost
// @Summary 更新HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrHost true "更新HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /hmdrHost/updateHmdrHost [put]
export const updateHmdrHost = (data) => {
  return service({
    url: '/hmdrHost/updateHmdrHost',
    method: 'put',
    data
  })
}

// @Tags HmdrHost
// @Summary 用id查询HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrHost true "用id查询HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /hmdrHost/findHmdrHost [get]
export const findHmdrHost = (params) => {
  return service({
    url: '/hmdrHost/findHmdrHost',
    method: 'get',
    params
  })
}

// @Tags HmdrHost
// @Summary 分页获取HmdrHost列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取HmdrHost列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hmdrHost/getHmdrHostList [get]
export const getHmdrHostList = (params) => {
  return service({
    url: '/hmdrHost/getHmdrHostList',
    method: 'get',
    params
  })
}
