package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/server/api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/middleware"
)

func InitSysDictionaryRouter(Router *gin.RouterGroup) {
	SysDictionaryRouter := Router.Group("sysDictionary").Use(middleware.OperationRecord())
	{
		SysDictionaryRouter.POST("createSysDictionary", v1.CreateSysDictionary)   // 新建SysDictionary
		SysDictionaryRouter.DELETE("deleteSysDictionary", v1.DeleteSysDictionary) // 删除SysDictionary
		SysDictionaryRouter.PUT("updateSysDictionary", v1.UpdateSysDictionary)    // 更新SysDictionary
		SysDictionaryRouter.GET("findSysDictionary", v1.FindSysDictionary)        // 根据ID获取SysDictionary
		SysDictionaryRouter.GET("getSysDictionaryList", v1.GetSysDictionaryList)  // 获取SysDictionary列表
	}
}
