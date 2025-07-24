package utils

import (
	"net/http"

	"github.com/Divyshekhar/7th-sem-project-be/initializers"
	"github.com/Divyshekhar/7th-sem-project-be/models"
	"github.com/gin-gonic/gin"
)

func CheckUser(ctx *gin.Context) (*models.User, bool){
	userId, exists := ctx.Get("user_id")
	if !exists{
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authentication found"})
		return nil, false
	}
	var user models.User
	if err := initializers.Db.First(&user, userId).Error; err != nil{
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no user for the id found"})
		return nil, false
	}
	return &user, true
	

}
