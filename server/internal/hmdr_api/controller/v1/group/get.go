package group

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	"github.com/tanganyu1114/heimdallr-reborn/model/response"
	"go.uber.org/zap"
)

func (g *GroupController) Get(c *gin.Context) {
	var r v1.Group
	err := c.ShouldBindQuery(&r)
	if err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Any("err", err))
		response.FailWithMessage("参数异常", c)

		return
	}

	if group, err := g.svc.Groups().Get(c, r.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rehmdrGroup": group}, c)
	}

}
