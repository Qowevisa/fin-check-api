package main

import (
	"time"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	docs "git.qowevisa.me/Qowevisa/gonuts/docs"
	"git.qowevisa.me/Qowevisa/gonuts/handlers"
	"git.qowevisa.me/Qowevisa/gonuts/middleware"
	"git.qowevisa.me/Qowevisa/gonuts/tokens"
	"github.com/gin-contrib/cors"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

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
		categoriesRoutes := api.Group("/category", middleware.AuthMiddleware())
		{
			categoriesRoutes.POST("/add", handlers.CategoryAdd)
			categoriesRoutes.GET("/:id", handlers.CategoryGetId)
			categoriesRoutes.PUT("/edit/:id", handlers.CategoryPutId)
			categoriesRoutes.DELETE("/delete/:id", handlers.CategoryDeleteId)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	go tokens.StartTokens()
	r.Run(":3000")
}
