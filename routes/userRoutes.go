package routes

import (
	"github.com/Divyshekhar/7th-sem-project-be/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(ctx *gin.Engine) {
	userGroup := ctx.Group("/user")
	{
		userGroup.POST("/signup", controllers.CreateUser)
		userGroup.POST("/login", controllers.Login)
		userGroup.POST("/logout", controllers.Logout)
		userGroup.DELETE("/delete", controllers.Delete)
	}
}
