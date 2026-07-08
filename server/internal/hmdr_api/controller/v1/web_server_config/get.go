package web_server_config

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
)

func getErrorHandle(c *gin.Context, err error, okdetailobj interface{}, okmsg, failuremsg string) {
	if err != nil {
		if errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010) {
			global.GVA_LOG.Info("查询时，指纹校验失败!", zap.Error(err))
			response.FailWithMessage("指纹校验失败, 请重新查询, 刷新配置文件!", c)

			return
		}

		global.GVA_LOG.Info(failuremsg+"!", zap.Any("err", err))
		response.FailWithMessage(failuremsg, c)

		return
	}

	response.OkWithDetailed(okdetailobj, okmsg, c)
}

func (w *WebServerConfigController) GetOptions(c *gin.Context) {
	if groups, err := w.svc.WebServerConfigs().GetOptions(c); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(groups, "获取成功", c)
	}
}

func (w *WebServerConfigController) GetConfigTextLines(c *gin.Context) {
	var r metav1.WebServerOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	configmeta, err := w.svc.WebServerConfigs().GetConfig(c, r)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)

		return
	}

	lines, err := configmeta.Config.ConfigLines(false)
	if err != nil {
		global.GVA_LOG.Error("解析配置文本失败!", zap.Any("err", err))
		response.FailWithMessage("解析配置文本失败", c)

		return
	}
	response.OkWithDetailed(strings.Join(lines, "\n"), "获取成功", c)
}

func (w *WebServerConfigController) GetContextTextLines(c *gin.Context) {
	var r metav1.WebServerConfigTargetContextOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	getErrorHandle(func() (ginctx *gin.Context, err error, detail interface{}, okmsg, failuremsg string) {
		ginctx = c
		detail = ""
		okmsg = "获取成功"
		ctx, e := w.svc.WebServerConfigs().GetContext(c, r.WebServerOptions, r.OriginalFingerprints, r.ConfigContextPos)
		if e != nil {
			failuremsg = "获取失败"
			err = e
			return
		}

		lines, e := ctx.ConfigLines(false)
		if e != nil {
			failuremsg = "解析上下文配置文本失败"
			err = e
			return
		}
		detail = strings.Join(lines, "\n")
		return
	}())
}

func (w *WebServerConfigController) GetConfig(c *gin.Context) {
	var r metav1.WebServerOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	configmeta, err := w.svc.WebServerConfigs().GetConfig(c, r)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)

		return
	}

	response.OkWithDetailed(configmeta, "获取成功", c)
}

func (w *WebServerConfigController) GetIncludedConfigs(c *gin.Context) {
	var r metav1.WebServerConfigTargetContextOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	includes, err := w.svc.WebServerConfigs().GetIncludedConfigs(c, r.WebServerOptions, r.OriginalFingerprints, r.ConfigContextPos)
	getErrorHandle(c, err, includes, "获取包含的配置文件成功", "获取包含的配置文件失败")
}

func (w *WebServerConfigController) SearchContextPositions(c *gin.Context) {
	var r metav1.WebServerConfigContextPosSearchOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	poslist, err := w.svc.WebServerConfigs().SearchContextPositions(c, r.WebServerOptions, r.OriginalFingerprints, r.SearchKeywordsMeta)
	getErrorHandle(c, err, poslist, "搜索成功", "搜索失败")
}
