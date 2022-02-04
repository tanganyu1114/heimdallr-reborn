package host

import (
	"gin-vue-admin/global"
	ctlv1 "gin-vue-admin/internal/hmdr_api/controller/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HostController) Get(c *gin.Context) {
	id, err := ctlv1.ParseID(c)
	if err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Any("err", err))
		response.FailWithMessage("参数异常", c)

		return
	}

	if host, err := h.svc.Host().Get(c, id); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rehmdrHost": host}, c)
	}
}
