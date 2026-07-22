package host

import (
	"github.com/gin-gonic/gin"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	"go.uber.org/zap"
)

func (h *HostController) List(c *gin.Context) {
	var r metav1.ListOptions
	err := c.ShouldBindQuery(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	if list, err := h.svc.Hosts().List(c, r); err != nil {
		global.GVA_LOG.Error("获取失败", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list.Items,
			Total:    list.TotalCount,
			Page:     r.Page,
			PageSize: r.PageSize,
		}, "获取成功", c)
	}
}
