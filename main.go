package main

import (
	"github.com/Divyshekhar/7th-sem-project-be/initializers"
	"github.com/Divyshekhar/7th-sem-project-be/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDb()
	// initializers.Migrate()
	// initializers.SeedSubjects()
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://divyshekhar.vercel.app"},
		AllowMethods: []string{"POST", "GET", "OPTIONS", "DELETE", "PUT"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is up and running"})
	})

	routes.RegisterUserRoutes(router)
	routes.RegisterQuestionRoutes(router)

	router.Run(":8080")
}
