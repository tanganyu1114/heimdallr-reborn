package group

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	"github.com/tanganyu1114/heimdallr-reborn/model/response"
	"go.uber.org/zap"
)

func (g *GroupController) Update(c *gin.Context) {
	var r v1.Group
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}

	if err = g.svc.Groups().Update(c, r); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)

		return
	}

	response.OkWithMessage("更新成功", c)
}
