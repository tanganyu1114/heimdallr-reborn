package host

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *HostController) DeleteCollection(c *gin.Context) {
	var r metav1.IDsOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	if err = h.svc.Hosts().DeleteCollection(c, r); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)

		return
	}

	response.OkWithMessage("批量删除成功", c)
}
