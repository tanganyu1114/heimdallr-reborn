package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tanganyu1114/heimdallr-reborn/api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/middleware"
)

func InitWorkflowRouter(Router *gin.RouterGroup) {
	WorkflowRouter := Router.Group("workflow").Use(middleware.OperationRecord())
	{
		WorkflowRouter.POST("createWorkFlow", v1.CreateWorkFlow) // 创建工作流
	}
}
