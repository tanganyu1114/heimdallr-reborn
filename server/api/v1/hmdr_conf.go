package v1

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/response"
	"gin-vue-admin/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags conf
// @Summary 获取服务器信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /conf/getOptions [get]
func GetOptions(c *gin.Context) {
	if data, err := service.GetOptions(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}

// @Tags conf
// @Summary 获取配置文件信息
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /conf/getConfInfo [get]
func GetConfInfo(c *gin.Context) {
	var sf model.SearchConf
	_ = c.ShouldBind(&sf)
	data := service.GetConfInfo(sf)
	if data == nil {
		global.GVA_LOG.Error("获取失败!")
		response.FailWithMessage("获取失败", c)
		return
	} else {
		response.OkWithDetailed(string(*data), "获取成功", c)
	}
}
