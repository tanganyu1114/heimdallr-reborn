package group

import (
	"gin-vue-admin/global"
	ctlv1 "gin-vue-admin/internal/hmdr_api/controller/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (g *GroupController) Delete(c *gin.Context) {
	id, err := ctlv1.ParseID(c)
	if err != nil {
		global.GVA_LOG.Error("参数异常!", zap.Any("err", err))
		response.FailWithMessage("参数异常", c)

		return
	}

	if err = g.svc.Group().Delete(c, id); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)

		return
	}

	response.OkWithMessage("删除成功", c)
}
