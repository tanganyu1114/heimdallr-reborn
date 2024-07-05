import service from '@/utils/request'

// @Tags HmdrStatistics
// @Summary 获取代理配置简要信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /hmdr-statistics/proxy-svc-brief [post]
export const getProxyServiceInfo = (data) => {
  return service({
    url: '/hmdr-statistics/proxy-svc-brief',
    method: 'post',
    data
  })
}
