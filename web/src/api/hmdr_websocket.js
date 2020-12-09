import service from '@/utils/request'

// @Tags websocket
// @Summary 获取配置文件内容
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /conf/getLogsInfo [post]
export const getLogsInfo = (data) => {
  return service({
    url: '/hmdrWebSocket/getLogsInfo',
    method: 'post',
    data
  })
}
