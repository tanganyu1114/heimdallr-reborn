package group

import (
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (g *GroupController) Delete(c *gin.Context) {
	var r v1.Group
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Any("err", err))
		response.FailWithMessage("参数异常", c)

		return
	}

	if err = g.svc.Groups().Delete(c, r.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)

		return
	}

	response.OkWithMessage("删除成功", c)
}
