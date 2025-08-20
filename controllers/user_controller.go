package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Divyshekhar/7th-sem-project-be/initializers"
	"github.com/Divyshekhar/7th-sem-project-be/models"
	"github.com/Divyshekhar/7th-sem-project-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(ctx *gin.Context) {
	var body struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid output"})
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to has password"})
		return
	}
	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  string(hashedPass),
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
	ctx.SetSameSite(http.SameSiteLaxMode)
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
	if err := ctx.ShouldBindJSON(&body); err != nil {
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
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("jwt_token", tokenStr, 3600*72, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User logged in",
		"user":    user,
	})
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user logged out",
	})
}

// remove in prod
func Delete(ctx *gin.Context) {
	initializers.Db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.User{})
}

func Update(ctx *gin.Context) {
	user, ok := utils.CheckUser(ctx)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token or invalid token"})
		return
	}
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	allowed := map[string]bool{
		"FirstName": true,
		"LastName":  true,
		"email":     true,
		"password":  true,
	}
	safeUpdates := map[string]interface{}{}
	for key, value := range updates {
		if allowed[key] {
			safeUpdates[key] = value
		}
	}
	if pw, ok := safeUpdates["password"].(string); ok {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating hash while updating account"})
			return
		}
		safeUpdates["password"] = hashedPass
	}

	if len(safeUpdates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	if err := initializers.Db.Model(&user).Updates(safeUpdates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user details"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User data updated successfully",
		"patient": user,
	})
}

type UpdatedUser struct {
	Id        uint
	FirstName string
	LastName  string
	Email     string
}

func GetInfo(c *gin.Context) {
	user, ok := utils.CheckUser(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "could not find user"})
	}
	response := UpdatedUser{
		Id: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
	}
		c.JSON(http.StatusOK, response)
}
