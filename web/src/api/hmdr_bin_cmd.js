import service from '@/utils/request'

// @Tags bin-cmd
// @Summary 提交Web服务端二进制命令工具执行请求
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{"successful":true,"stdout":"...":"stderr":"..."},"msg":"命令执行请求成功"}"
// @Router /bin-cmd/exec [post]
export const binCMDExec = (data) => {
  return service({
    url: '/bin-cmd/exec',
    method: 'post',
    data
  })
}
