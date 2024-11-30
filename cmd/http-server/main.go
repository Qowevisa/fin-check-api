package main

import (
	"fmt"
	"os"
	"time"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	docs "git.qowevisa.me/Qowevisa/fin-check-api/docs"
	"git.qowevisa.me/Qowevisa/fin-check-api/handlers"
	"git.qowevisa.me/Qowevisa/fin-check-api/middleware"
	"git.qowevisa.me/Qowevisa/fin-check-api/tokens"
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

// @host gonapi.qowevisa.click
// @BasePath /api
func main() {
	if err := db.Init(); err != nil {
		fmt.Printf("ERROR: db.Init: %v\n", err)
		os.Exit(1)
	}
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes defined in the routes package
	api := r.Group("/api")
	{
		api.GET("/ping", handlers.PingGet)
		api.GET("/authping", middleware.AuthMiddleware(), handlers.PingGet)
		userRoutes := api.Group("/user")
		{
			userRoutes.POST("/register", handlers.UserRegister)
			userRoutes.POST("/login", handlers.UserLogin)
		}
		cardsRoutes := api.Group("/card", middleware.AuthMiddleware())
		{
			cardsRoutes.POST("/add", handlers.CardAdd)
			cardsRoutes.GET("/:id", handlers.CardGetId)
			cardsRoutes.GET("/all", handlers.CardGetAll)
			cardsRoutes.PUT("/edit/:id", handlers.CardPutId)
			cardsRoutes.DELETE("/delete/:id", handlers.CardDeleteId)
		}
		categoriesRoutes := api.Group("/category", middleware.AuthMiddleware())
		{
			categoriesRoutes.POST("/add", handlers.CategoryAdd)
			categoriesRoutes.GET("/:id", handlers.CategoryGetId)
			categoriesRoutes.GET("/all", handlers.CategoryGetAll)
			categoriesRoutes.PUT("/edit/:id", handlers.CategoryPutId)
			categoriesRoutes.DELETE("/delete/:id", handlers.CategoryDeleteId)
		}
		debtRoutes := api.Group("/debt", middleware.AuthMiddleware())
		{
			debtRoutes.POST("/add", handlers.DebtAdd)
			debtRoutes.GET("/:id", handlers.DebtGetId)
			debtRoutes.PUT("/edit/:id", handlers.DebtPutId)
			debtRoutes.DELETE("/delete/:id", handlers.DebtDeleteId)
		}
		incomeRoutes := api.Group("/income", middleware.AuthMiddleware())
		{
			incomeRoutes.POST("/add", handlers.IncomeAdd)
			incomeRoutes.GET("/:id", handlers.IncomeGetId)
			incomeRoutes.GET("/all", handlers.IncomeGetAll)
			incomeRoutes.PUT("/edit/:id", handlers.IncomePutId)
			incomeRoutes.DELETE("/delete/:id", handlers.IncomeDeleteId)
		}
		typesRoutes := api.Group("/type", middleware.AuthMiddleware())
		{
			typesRoutes.POST("/add", handlers.TypeAdd)
			typesRoutes.GET("/:id", handlers.TypeGetId)
			typesRoutes.GET("/all", handlers.TypeGetAll)
			typesRoutes.PUT("/edit/:id", handlers.TypePutId)
			typesRoutes.DELETE("/delete/:id", handlers.TypeDeleteId)
		}
		expensesRoutes := api.Group("/expense", middleware.AuthMiddleware())
		{
			expensesRoutes.POST("/add", handlers.ExpenseAdd)
			expensesRoutes.POST("/bulk_create", handlers.ExpenseBulkCreate)
			expensesRoutes.GET("/:id", handlers.ExpenseGetId)
			expensesRoutes.GET("/all", handlers.ExpenseGetAll)
			expensesRoutes.PUT("/edit/:id", handlers.ExpensePutId)
			expensesRoutes.DELETE("/delete/:id", handlers.ExpenseDeleteId)
		}
		transfersRoutes := api.Group("/transfer", middleware.AuthMiddleware())
		{
			transfersRoutes.POST("/add", handlers.TransferAdd)
			transfersRoutes.GET("/:id", handlers.TransferGetId)
			transfersRoutes.GET("/all", handlers.TransferGetAll)
			transfersRoutes.PUT("/edit/:id", handlers.TransferPutId)
			transfersRoutes.DELETE("/delete/:id", handlers.TransferDeleteId)
		}
		itemRoutes := api.Group("/item", middleware.AuthMiddleware())
		{
			itemRoutes.GET("/:id", handlers.ItemGetId)
			itemRoutes.GET("/all", handlers.ItemGetAll)
			itemRoutes.POST("/filter", handlers.ItemPostFilter)
			itemRoutes.DELETE("/delete/:id", handlers.ItemDeleteId)
		}
		metricRoutes := api.Group("/metric", middleware.AuthMiddleware())
		{
			metricRoutes.GET("/all", handlers.MetricGetAll)
		}
		paymentRoutes := api.Group("/payment", middleware.AuthMiddleware())
		{
			paymentRoutes.POST("/add", handlers.PaymentAdd)
			paymentRoutes.GET("/all", handlers.PaymentGetAll)
		}
		currencyRoutes := api.Group("/currency", middleware.AuthMiddleware())
		{
			currencyRoutes.GET("/all", handlers.CurrencyGetAll)
		}
		statisticRoute := api.Group("/statistics", middleware.AuthMiddleware())
		{
			statisticRoute.GET("/type", handlers.StatisticsGetAllSpendingsForTypes)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	tokens.Init()
	r.Run("127.0.0.1:3000")
}
