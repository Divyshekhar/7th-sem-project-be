package main

import (
	intializers "github.com/Divyshekhar/7th-sem-project-be/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	intializers.LoadEnv()
	intializers.ConnectDb()
}

func main() {
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is up and running"})
	})

	router.Run(":8080")
}
