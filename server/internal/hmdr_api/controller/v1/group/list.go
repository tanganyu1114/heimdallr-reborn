package group

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (g *GroupController) List(c *gin.Context) {
	var r metav1.ListOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	if list, err := g.svc.Groups().List(c, r); err != nil {
		global.GVA_LOG.Error("获取失败", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(list, "获取成功", c)
	}
}
