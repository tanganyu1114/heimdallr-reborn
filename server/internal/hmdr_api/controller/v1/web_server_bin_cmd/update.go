package web_server_bin_cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	"go.uber.org/zap"
)

// @Tags bin-cmd
// @Summary 传入请求中命令执行参数，请求web服务端二进制命令执行命令
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{"successful":true,"stdout":"...":"stderr":"..."},"msg":"命令执行请求成功"}"
// @Router /bin-cmd/exec [post]
func (w *WebServerBinCMDController) Exec(c *gin.Context) {
	var r metav1.WebServerBinCMDExecRequest
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	resp, err := w.svc.WebServerBinCMD().Exec(c, r.WebServerOptions, r.Args...)
	if err != nil {
		global.GVA_LOG.Error("命令执行请求失败!", zap.Any("err", err))
		response.FailWithDetailed(resp, "命令执行请求失败", c)

		return
	}

	response.OkWithDetailed(resp, "命令执行请求成功", c)
}
