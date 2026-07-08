package agent_info

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	"github.com/tanganyu1114/heimdallr-reborn/model/response"
	"go.uber.org/zap"
)

func (a *AgentInfoController) Get(c *gin.Context) {
	if groups, err := a.svc.AgentInfos().Get(c); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
		return
	} else {
		response.OkWithDetailed(groups, "获取成功", c)
	}
}
