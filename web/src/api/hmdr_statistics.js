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

// @Tags HmdrStatistics
// @Summary 代理服务网络连通性检查
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"网络连通性检查成功"}"
// @Router /hmdr-statistics/conn-check-of-proxy-svc [post]
export const connectivityCheckOfProxyService = (data) => {
  return service({
    url: '/hmdr-statistics/conn-check-of-proxy-svc',
    method: 'post',
    data
  })
}

// @Tags HmdrStatistics
// @Summary 导出代理服务信息为Excel
// @Security ApiKeyAuth
// @Produce  application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} file "Excel文件"
// @Router /hmdr-statistics/export-proxy-svc-excel [post]
export const exportProxyServiceInfoToExcel = (data) => {
  return service({
    url: '/hmdr-statistics/export-proxy-svc-excel',
    method: 'post',
    data,
    responseType: 'arraybuffer' // 使用 arraybuffer 而不是 blob，便于判断响应类型
  })
}
