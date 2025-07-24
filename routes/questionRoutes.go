package routes

import (
	"github.com/Divyshekhar/7th-sem-project-be/controllers"
	"github.com/Divyshekhar/7th-sem-project-be/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterQuestionRoutes(ctx *gin.Engine){
	questionGroup := ctx.Group("/question")
	{
		questionGroup.POST("/:subject", middleware.RequireAuth(), controllers.GetQuestions)	
	}
}