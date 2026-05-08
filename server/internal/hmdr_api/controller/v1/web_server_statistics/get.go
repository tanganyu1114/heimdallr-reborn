package web_server_statistics

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (w *WebServerStatisticsController) GetProxyServiceInfo(c *gin.Context) {
	var r metav1.WebServerOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	proxyInfo, err := w.svc.WebServerStatistics().GetProxyServiceInfo(c, r)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)

		return
	}

	response.OkWithDetailed(proxyInfo, "获取成功", c)
}

func (w *WebServerStatisticsController) ConnectivityCheckOfProxyService(c *gin.Context) {
	var r metav1.ConnectivityCheckOfProxiedServersRequestOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	proxyInfo, err := w.svc.WebServerStatistics().ConnectivityCheckOfProxyService(c, r.WebServerOptions, r.ConfigContextPos)
	if err != nil {
		global.GVA_LOG.Error("网络连通性检查失败!", zap.Any("err", err))
		response.FailWithMessage("网络连通性检查失败", c)

		return
	}

	response.OkWithDetailed(proxyInfo, "网络连通性检查成功", c)
}

func (w *WebServerStatisticsController) ExportProxyServiceInfoToExcel(c *gin.Context) {
	var r metav1.WebServerOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	excelData, err := w.svc.WebServerStatistics().ExportProxyServiceInfoToExcel(c, r)
	if err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Any("err", err))
		response.FailWithMessage("导出失败", c)

		return
	}

	// 设置响应头，触发文件下载
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=proxy_service_info.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}
