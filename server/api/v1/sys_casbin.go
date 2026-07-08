package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	"github.com/tanganyu1114/heimdallr-reborn/server/service"
	"github.com/tanganyu1114/heimdallr-reborn/server/utils"
	"go.uber.org/zap"
)

// @Tags Casbin
// @Summary 更新角色api权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "权限id, 权限模型列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /casbin/UpdateCasbin [post]
func UpdateCasbin(c *gin.Context) {
	var cmr request.CasbinInReceive
	_ = c.ShouldBindJSON(&cmr)
	if err := utils.Verify(cmr, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.UpdateCasbin(cmr.AuthorityId, cmr.CasbinInfos); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Casbin
// @Summary 获取权限列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CasbinInReceive true "权限id, 权限模型列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /casbin/getPolicyPathByAuthorityId [post]
func GetPolicyPathByAuthorityId(c *gin.Context) {
	var casbin request.CasbinInReceive
	_ = c.ShouldBindJSON(&casbin)
	if err := utils.Verify(casbin, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	paths := service.GetPolicyPathByAuthorityId(casbin.AuthorityId)
	response.OkWithDetailed(response.PolicyPathResponse{Paths: paths}, "获取成功", c)
}
