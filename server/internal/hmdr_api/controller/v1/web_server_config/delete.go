package web_server_config

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
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
