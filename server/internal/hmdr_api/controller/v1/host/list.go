package host

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
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
