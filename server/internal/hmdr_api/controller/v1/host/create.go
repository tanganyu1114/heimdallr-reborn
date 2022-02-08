package host

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HostController) Create(c *gin.Context) {
	var r v1.Host
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	if err = h.svc.Hosts().Create(c, r); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)

		return
	}

	response.OkWithMessage("创建成功", c)
}
