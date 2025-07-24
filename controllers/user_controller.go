package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Divyshekhar/7th-sem-project-be/initializers"
	"github.com/Divyshekhar/7th-sem-project-be/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(ctx *gin.Context) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid output"})
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to has password"})
		return
	}
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPass),
	}
	tx := initializers.Db.Create(&user)
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating jwt string"})
		return
	}
	ctx.SetCookie("jwt_token", tokenStr, 3600*72, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(ctx *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid inputs"})
		return
	}
	var user models.User
	result := initializers.Db.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error signing jwt"})
		return
	}
	ctx.SetCookie("jwt_token", tokenStr, 3600*72, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User logged in",
		"user":    user,
	})
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
}
// remove in prod
func Delete(ctx *gin.Context) {
	initializers.Db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
}
