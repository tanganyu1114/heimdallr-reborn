package web_server_config

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags conf
// @Summary 插入需被克隆的配置上下文
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"新增成功"}"
// @Router /conf/insert-clone-ctx [post]
func (w *WebServerConfigController) InsertWithClone(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	err = w.svc.WebServerConfigs().InsertWithClone(c, r.WebServerOptions, r.TargetConfigContextOptions)
	if err != nil {
		global.GVA_LOG.Error("新增失败!")
		response.FailWithMessage("新增失败", c)

		return
	}

	response.OkWithMessage("新增成功", c)
}

func (w *WebServerConfigController) InsertWithNew(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	err = w.svc.WebServerConfigs().InsertWithNew(c, r.WebServerOptions, r.TargetConfigContextOptions)
	if err != nil {
		global.GVA_LOG.Error("新增失败!")
		response.FailWithMessage("新增失败", c)

		return
	}

	response.OkWithMessage("新增成功", c)
}

func (w *WebServerConfigController) ModifyWithClone(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	err = w.svc.WebServerConfigs().ModifyWithClone(c, r.WebServerOptions, r.TargetConfigContextOptions)
	if err != nil {
		global.GVA_LOG.Error("修改失败!")
		response.FailWithMessage("修改失败", c)

		return
	}

	response.OkWithMessage("修改成功", c)
}

func (w *WebServerConfigController) ModifyWithNew(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	err = w.svc.WebServerConfigs().ModifyWithNew(c, r.WebServerOptions, r.TargetConfigContextOptions)
	if err != nil {
		global.GVA_LOG.Error("修改失败!")
		response.FailWithMessage("修改失败", c)

		return
	}

	response.OkWithMessage("修改成功", c)
}

func (w *WebServerConfigController) ModifyContextValue(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	err = w.svc.WebServerConfigs().ModifyContextValue(c, r.WebServerOptions, r.TargetConfigContextOptions)
	if err != nil {
		global.GVA_LOG.Error("修改失败!")
		response.FailWithMessage("修改失败", c)

		return
	}

	response.OkWithMessage("修改成功", c)
}

// @Tags conf
// @Summary 移动指定配置上下文
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"修改成功"}"
// @Router /conf/move-ctx [post]
func (w *WebServerConfigController) Move(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	err = w.svc.WebServerConfigs().Move(c, r.WebServerOptions, r.TargetConfigContextOptions)
	if err != nil {
		global.GVA_LOG.Error("修改失败!")
		response.FailWithMessage("修改失败", c)

		return
	}

	response.OkWithMessage("修改成功", c)
}
