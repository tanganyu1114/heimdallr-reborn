package group

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	"go.uber.org/zap"
)

func (g *GroupController) DeleteCollections(c *gin.Context) {
	var r metav1.IDsOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	if err = g.svc.Groups().DeleteCollections(c, r); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)

		return
	}

	response.OkWithMessage("批量删除成功", c)
}
