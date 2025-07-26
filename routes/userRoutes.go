package routes

import (
	"github.com/Divyshekhar/7th-sem-project-be/controllers"
	"github.com/Divyshekhar/7th-sem-project-be/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(ctx *gin.Engine) {
	userGroup := ctx.Group("/user")
	{
		userGroup.POST("/signup", controllers.CreateUser)
		userGroup.POST("/login", controllers.Login)
		userGroup.POST("/logout", controllers.Logout)
		userGroup.DELETE("/delete", controllers.Delete)
		userGroup.PATCH("/update", middleware.RequireAuth(), controllers.Update)
	}
}
