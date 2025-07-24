package controllers

import (
	"net/http"

	"github.com/Divyshekhar/7th-sem-project-be/utils"
	"github.com/gin-gonic/gin"
)

func GetQuestions(ctx *gin.Context) {
	subject := ctx.Param("subject")
	if subject == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no id specified"})
		return
	}
	// check for the user in the database
	user, ok := utils.CheckUser(ctx)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no user found"})
		return
	}
	output, err := utils.Generate(*user, subject)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error generating the llm response"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "response generated",
		"response": output,
	})

}
