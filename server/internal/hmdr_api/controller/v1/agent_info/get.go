package agent_info

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (a *AgentInfoController) Get(c *gin.Context) {
	if groups, err := a.svc.AgentInfo().Get(c); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	} else {
		response.OkWithDetailed(groups, "获取成功", c)
	}
}
