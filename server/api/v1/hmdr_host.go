package v1

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
	"gin-vue-admin/model/response"
	"gin-vue-admin/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags HmdrHost
// @Summary 创建HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrHost true "创建HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hmdrHost/createHmdrHost [post]
func CreateHmdrHost(c *gin.Context) {
	var hmdrHost model.HmdrHost
	_ = c.ShouldBindJSON(&hmdrHost)
	if err := service.CreateHmdrHost(hmdrHost); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags HmdrHost
// @Summary 删除HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrHost true "删除HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /hmdrHost/deleteHmdrHost [delete]
func DeleteHmdrHost(c *gin.Context) {
	var hmdrHost model.HmdrHost
	_ = c.ShouldBindJSON(&hmdrHost)
	if err := service.DeleteHmdrHost(hmdrHost); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags HmdrHost
// @Summary 批量删除HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /hmdrHost/deleteHmdrHostByIds [delete]
func DeleteHmdrHostByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteHmdrHostByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags HmdrHost
// @Summary 更新HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrHost true "更新HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /hmdrHost/updateHmdrHost [put]
func UpdateHmdrHost(c *gin.Context) {
	var hmdrHost model.HmdrHost
	_ = c.ShouldBindJSON(&hmdrHost)
	if err := service.UpdateHmdrHost(&hmdrHost); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags HmdrHost
// @Summary 用id查询HmdrHost
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrHost true "用id查询HmdrHost"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /hmdrHost/findHmdrHost [get]
func FindHmdrHost(c *gin.Context) {
	var hmdrHost model.HmdrHost
	_ = c.ShouldBindQuery(&hmdrHost)
	if err, rehmdrHost := service.GetHmdrHost(hmdrHost.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rehmdrHost": rehmdrHost}, c)
	}
}

// @Tags HmdrHost
// @Summary 分页获取HmdrHost列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.HmdrHostSearch true "分页获取HmdrHost列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hmdrHost/getHmdrHostList [get]
func GetHmdrHostList(c *gin.Context) {
	var pageInfo request.HmdrHostSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetHmdrHostInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
