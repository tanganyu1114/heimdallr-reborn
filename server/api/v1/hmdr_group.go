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

// @Tags HmdrGroup
// @Summary 创建HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrGroup true "创建HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hmdrGroup/createHmdrGroup [post]
func CreateHmdrGroup(c *gin.Context) {
	var hmdrGroup model.HmdrGroup
	_ = c.ShouldBindJSON(&hmdrGroup)
	if err := service.CreateHmdrGroup(hmdrGroup); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags HmdrGroup
// @Summary 删除HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrGroup true "删除HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /hmdrGroup/deleteHmdrGroup [delete]
func DeleteHmdrGroup(c *gin.Context) {
	var hmdrGroup model.HmdrGroup
	_ = c.ShouldBindJSON(&hmdrGroup)
	if err := service.DeleteHmdrGroup(hmdrGroup); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags HmdrGroup
// @Summary 批量删除HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /hmdrGroup/deleteHmdrGroupByIds [delete]
func DeleteHmdrGroupByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteHmdrGroupByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags HmdrGroup
// @Summary 更新HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrGroup true "更新HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /hmdrGroup/updateHmdrGroup [put]
func UpdateHmdrGroup(c *gin.Context) {
	var hmdrGroup model.HmdrGroup
	_ = c.ShouldBindJSON(&hmdrGroup)
	if err := service.UpdateHmdrGroup(&hmdrGroup); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags HmdrGroup
// @Summary 用id查询HmdrGroup
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.HmdrGroup true "用id查询HmdrGroup"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /hmdrGroup/findHmdrGroup [get]
func FindHmdrGroup(c *gin.Context) {
	var hmdrGroup model.HmdrGroup
	_ = c.ShouldBindQuery(&hmdrGroup)
	if err, rehmdrGroup := service.GetHmdrGroup(hmdrGroup.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rehmdrGroup": rehmdrGroup}, c)
	}
}

// @Tags HmdrGroup
// @Summary 分页获取HmdrGroup列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.HmdrGroupSearch true "分页获取HmdrGroup列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /hmdrGroup/getHmdrGroupList [get]
func GetHmdrGroupList(c *gin.Context) {
	var pageInfo request.HmdrGroupSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetHmdrGroupInfoList(pageInfo); err != nil {
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
