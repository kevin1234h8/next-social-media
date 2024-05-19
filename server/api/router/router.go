package router

import (
	"social/project/api/handler"
	"social/project/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeUserRouter(router *gin.Engine) {
	userGroup := router.Group("/api/user")
	{
		userGroup.GET("/", middleware.AuthMiddleware(), handler.GetUserInfo)
		userGroup.POST("/register", handler.Register)
		userGroup.POST("/login", handler.Login)
	}
}
