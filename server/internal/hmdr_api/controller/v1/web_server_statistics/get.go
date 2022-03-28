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
	err := c.ShouldBind(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	proxyInfo, err := w.svc.WebServerStatistics().GetProxyServiceInfo(c, r)
	if err != nil {
		global.GVA_LOG.Error("获取失败!")
		response.FailWithMessage("获取失败", c)

		return
	}
	//jsonData, err := json.Marshal(proxyInfo)
	//if err != nil {
	//	global.GVA_LOG.Error("解析失败!")
	//	response.FailWithMessage("解析失败", c)
	//
	//	return
	//}

	response.OkWithDetailed(proxyInfo, "获取成功", c)
	//response.OkWithDetailed(string(jsonData), "获取成功", c)
}
