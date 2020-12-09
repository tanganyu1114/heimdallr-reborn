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
