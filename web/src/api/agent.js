import service from '@/utils/request'

// @Tags agent
// @Summary 获取配置文件内容
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /system/getSystemConfig [post]
export const getAgentInfo = () => {
  return service({
    url: '/agent/getAgentInfo',
    method: 'get'
  })
}
