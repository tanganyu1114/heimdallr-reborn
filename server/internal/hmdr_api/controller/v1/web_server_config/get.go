package web_server_config

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

func (w *WebServerConfigController) GetOptions(c *gin.Context) {
	if groups, err := w.svc.WebServerConfigs().GetOptions(c); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(groups, "获取成功", c)
	}
}

func (w *WebServerConfigController) GetConfigTextLines(c *gin.Context) {
	// TODO: 新增支持指定配置上下文对象文本格式获取的功能
	var r metav1.WebServerOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	config, err := w.svc.WebServerConfigs().GetConfig(c, r)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)

		return
	}
	response.OkWithDetailed(strings.Join(config.TextLines(), "\n"), "获取成功", c)
}

func (w *WebServerConfigController) GetContextTextLines(c *gin.Context) {
	var r metav1.WebServerConfigTargetContextOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	ctx, err := w.svc.WebServerConfigs().GetContext(c, r.WebServerOptions, r.ConfigContextPos)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)

		return
	}
	lines, err := ctx.ConfigLines(false)
	if err != nil {
		global.GVA_LOG.Error("解析上下文配置文本失败！", zap.Any("err", err))
		response.FailWithMessage("解析上下文配置文本失败", c)
	}
	response.OkWithDetailed(strings.Join(lines, "\n"), "获取成功", c)
}

func (w *WebServerConfigController) GetConfig(c *gin.Context) {
	var r metav1.WebServerOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	config, err := w.svc.WebServerConfigs().GetConfig(c, r)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)

		return
	}

	response.OkWithDetailed(config.Main(), "获取成功", c)
}
