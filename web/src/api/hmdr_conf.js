import service from '@/utils/request'

// @Tags conf
// @Summary 获取配置文件内容
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /conf/getOptions [get]
export const getOptions = () => {
  return service({
    url: '/conf/getOptions',
    method: 'get'
  })
}

// @Tags conf
// @Summary 获取配置文件内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /conf/getConfinfo [post]
export const getConfInfo = (data) => {
  return service({
    url: '/conf/getConfInfo',
    method: 'post',
    data
  })
}

// @Tags conf
// @Summary 获取配置上下文配置文本内容
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /conf/get-context-text [post]
export const getContextText = (data) => {
  return service({
    url: '/conf/get-context-text',
    method: 'post',
    data
  })
}

// @Tags conf
// @Summary 获取配置文件json全量数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /conf/get-conf-struct [post]
export const getConfStruct = (data) => {
  return service({
    url: '/conf/get-conf-struct',
    method: 'post',
    data
  })
}

// @Tags conf
// @Summary 获取包含的配置文件路径列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{["some/included_config.conf"]},"msg":"获取包含的配置文件成功"}"
// @Router /conf/get-conf-struct [post]
export const getIncludes = (data) => {
  return service({
    url: '/conf/get-includes',
    method: 'post',
    data
  })
}

// @Tags conf
// @Summary 删除指定配置上下文
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"删除成功"}"
// @Router /conf/remove-ctx [delete]
export const removeCtx = (data) => {
  return service({
    url: '/conf/remove-ctx',
    method: 'delete',
    data
  })
}

// @Tags conf
// @Summary 修改指定配置上下文启用状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"修改启用状态成功"}"
// @Router /conf/change-ctx-enabled-state [post]
export const changeCtxEnabledState = (data) => {
  return service({
    url: '/conf/change-ctx-enabled-state',
    method: 'post',
    data
  })
}

// @Tags conf
// @Summary 修改指定配置上下文配置参数值
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"修改成功"}"
// @Router /conf/modify-ctx-value [post]
export const modifyCtxValue = (data) => {
  return service({
    url: '/conf/modify-ctx-value',
    method: 'post',
    data
  })
}

// @Tags conf
// @Summary 插入需新增的配置上下文
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"新增成功"}"
// @Router /conf/insert-new-ctx [post]
export const insertNewCtx = (data) => {
  return service({
    url: '/conf/insert-new-ctx',
    method: 'post',
    data
  })
}

// @Tags conf
// @Summary 插入需被克隆的配置上下文
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"新增成功"}"
// @Router /conf/insert-clone-ctx [post]
export const insertCloneCtx = (data) => {
  return service({
    url: '/conf/insert-clone-ctx',
    method: 'post',
    data
  })
}

// @Tags conf
// @Summary 移动指定配置上下文
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"修改成功"}"
// @Router /conf/move-ctx [post]
export const moveCtx = (data) => {
  return service({
    url: '/conf/move-ctx',
    method: 'post',
    data
  })
}
