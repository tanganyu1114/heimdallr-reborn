package host

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
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
