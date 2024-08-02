package main

import (
	"fmt"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"git.qowevisa.me/Qowevisa/gonuts/db"
	docs "git.qowevisa.me/Qowevisa/gonuts/docs"
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
	dbc := db.Connect()
	if dbc != nil {
		fmt.Printf("yay\n")
	}
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"

	// Routes defined in the routes package
	routes := r.Group("/api")
	{
		routes.GET("/home", getHome)
		routes.GET("/user/:name", getUser)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":3000")
}

// @Summary Says hello
// @Description Get an account by ID
// @Tags hello hello2
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} Account
// @Failure 400 {object} ErrorResponse
// @Router /home [get]
func getHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the API!",
	})
}

// @Summary Says hello1
// @Description Get an account by ID1
// @Tags hello hello2
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} Account
// @Failure 400 {object} ErrorResponse
// @Router /accounts/{id} [get]
func getUser(c *gin.Context) {
	name := c.Param("name")
	c.JSON(200, gin.H{
		"message": "Hello, " + name + "!",
	})
}
