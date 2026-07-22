package web_server_config

import (
	"github.com/gin-gonic/gin"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	"go.uber.org/zap"
)

func (w *WebServerConfigController) Remove(c *gin.Context) {
	var r metav1.WebServerConfigTargetContextOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().Remove(c, r.WebServerOptions, r.OriginalFingerprints, r.ConfigContextPos),
		"删除成功",
		"删除失败",
	)
}
