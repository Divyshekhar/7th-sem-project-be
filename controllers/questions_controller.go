package controllers

import (
	"net/http"
	"strconv"

	"github.com/Divyshekhar/7th-sem-project-be/initializers"
	"github.com/Divyshekhar/7th-sem-project-be/models"
	"github.com/Divyshekhar/7th-sem-project-be/utils"
	"github.com/gin-gonic/gin"
)

func GenerateQuestions(ctx *gin.Context) {
	subject := ctx.Param("subject")
	if subject == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no subject specified"})
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
func GetQuestion(ctx *gin.Context) {
	subject := ctx.Param("subject")
	if subject == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no subject specified"})
		return
	}

	user, ok := utils.CheckUser(ctx)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalid or not found"})
		return
	}

	var subjectRecod models.Subject
	err := initializers.Db.Where("name = ?", subject).First(&subjectRecod).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "subject not found"})
		return
	}

	var userSubject models.UserSubject
	err = initializers.Db.Where("user_id = ? AND subject_id = ?", user.ID, subjectRecod.ID).First(&userSubject).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user is not enrolled in the subject"})
		return
	}

	// Pagination logic
	pageStr := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * 10

	var questions []models.Question
	err = initializers.Db.
		Where("user_subject_id = ?", userSubject.ID).
		Offset(offset).
		Limit(10).
		Find(&questions).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching questions"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"questions": questions,
		"page":      page,
	})
}
