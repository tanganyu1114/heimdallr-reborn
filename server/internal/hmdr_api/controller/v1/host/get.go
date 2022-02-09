package host

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HostController) Get(c *gin.Context) {
	var r v1.Host
	err := c.ShouldBindQuery(&r)
	if err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Any("err", err))
		response.FailWithMessage("参数异常", c)

		return
	}

	if host, err := h.svc.Hosts().Get(c, r.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rehmdrHost": host}, c)
	}
}
