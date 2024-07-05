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
// @Router /conf/getConfinfo [get]
export const getConfInfo = (data) => {
  return service({
    url: '/conf/getConfInfo',
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
// @Router /conf/get-conf-struct [get]
export const getConfStruct = (data) => {
  return service({
    url: '/conf/get-conf-struct',
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
