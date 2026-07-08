package web_server_config

import (
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
)

func updateErrorHandle(c *gin.Context, err error, okmsg, failuremsg string) {
	if err != nil {
		if errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010) {
			global.GVA_LOG.Warn("更新时，指纹校验失败!", zap.Error(err))
			response.FailWithMessage("指纹校验失败, 请重新查询, 刷新配置文件!", c)

			return
		}

		global.GVA_LOG.Error(failuremsg+"!", zap.Any("err", err))
		response.FailWithMessage(failuremsg, c)

		return
	}

	response.OkWithMessage(okmsg, c)
}

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
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().InsertWithClone(c, r.WebServerOptions, r.OriginalFingerprints, r.TargetConfigContextOptions, r.DisableTheTarget),
		"新增成功",
		"新增失败",
	)
}

func (w *WebServerConfigController) InsertWithNew(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().InsertWithNew(c, r.WebServerOptions, r.OriginalFingerprints, r.TargetConfigContextOptions, r.DisableTheTarget),
		"新增成功",
		"新增失败",
	)
}

func (w *WebServerConfigController) UpdateConfig(c *gin.Context) {
	var r metav1.WebServerConfigUpdateOptions
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().UpdateConfig(c, r.WebServerOptions, r.OriginalFingerprints, r.Data),
		"更新成功",
		"更新失败",
	)
}

func (w *WebServerConfigController) ModifyWithClone(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().ModifyWithClone(c, r.WebServerOptions, r.OriginalFingerprints, r.TargetConfigContextOptions),
		"修改成功",
		"修改失败",
	)
}

func (w *WebServerConfigController) ModifyWithNew(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().ModifyWithNew(c, r.WebServerOptions, r.OriginalFingerprints, r.TargetConfigContextOptions),
		"修改成功",
		"修改失败",
	)
}

// @Tags conf
// @Summary 修改指定配置上下文启用状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"修改启用状态成功"}"
// @Router /conf/change-ctx-enabled-state [post]
func (w *WebServerConfigController) ChangeContextEnabledState(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().ChangeContextEnabledState(c, r.WebServerOptions, r.OriginalFingerprints, r.TargetConfigContextOptions),
		"修改启用状态成功",
		"修改启用状态失败",
	)
}

func (w *WebServerConfigController) ModifyContextValue(c *gin.Context) {
	var r metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
	err := c.ShouldBindJSON(&r)
	if err != nil {
		global.GVA_LOG.Error("解析失败!", zap.Any("err", err))
		response.FailWithMessage("解析失败", c)

		return
	}
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().ModifyContextValue(c, r.WebServerOptions, r.OriginalFingerprints, r.TargetConfigContextOptions),
		"修改成功",
		"修改失败",
	)
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
	updateErrorHandle(
		c,
		w.svc.WebServerConfigs().Move(c, r.WebServerOptions, r.OriginalFingerprints, r.TargetConfigContextOptions, r.DisableTheTarget),
		"修改成功",
		"修改失败",
	)
}
