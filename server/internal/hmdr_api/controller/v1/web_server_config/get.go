package web_server_config

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (w *WebServerConfigController) GetOptions(c *gin.Context) {
	if groups, err := w.svc.WebServerConfigs().GetOptions(c); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(groups, "获取成功", c)
	}
}

func (w *WebServerConfigController) GetConfig(c *gin.Context) {
	var r metav1.WebServerOptions
	err := c.ShouldBind(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	data, err := w.svc.WebServerConfigs().GetConfig(c, r)
	if err != nil {
		global.GVA_LOG.Error("获取失败!")
		response.FailWithMessage("获取失败", c)

		return
	}

	response.OkWithDetailed(string(data), "获取成功", c)
}
