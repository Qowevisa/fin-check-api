package main

import (
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	docs "git.qowevisa.me/Qowevisa/gonuts/docs"
	"git.qowevisa.me/Qowevisa/gonuts/handlers"
	"git.qowevisa.me/Qowevisa/gonuts/middleware"
	"git.qowevisa.me/Qowevisa/gonuts/tokens"
	"github.com/gin-gonic/gin"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api
func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"

	// Routes defined in the routes package
	api := r.Group("/api")
	{
		userRoutes := api.Group("/user")
		{
			userRoutes.POST("/register", handlers.UserRegister)
			userRoutes.POST("/login", handlers.UserLogin)
		}
		cardsRoutes := api.Group("/card", middleware.AuthMiddleware())
		{
			cardsRoutes.POST("/add", handlers.CardAdd)
			cardsRoutes.GET("/:id", handlers.CardGetId)
			cardsRoutes.PUT("/edit/:id", handlers.CardPutId)
			cardsRoutes.DELETE("/delete/:id", handlers.CardDeleteId)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go tokens.StartTokens()
	r.Run(":3000")
}
