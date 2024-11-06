package handlers

import (
	"log"

	"git.qowevisa.me/Qowevisa/fin-check-api/consts"
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/tokens"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

// @Summary Register an user
// @Description Creates user in database as db.User
// @Tags user
// @Accept json
// @Produce json
// @Param user body types.User true "User info"
// @Success 200 {object} types.Account
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /user/register [post]
func UserRegister(c *gin.Context) {
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	dbUser := &db.User{
		Username: user.Username,
		Password: user.Password,
	}
	dbc := db.Connect()
	if err := dbc.Create(dbUser).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	var token1 *tokens.Token
	if token, err := tokens.AddToken(dbUser.ID); err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	} else {
		if token == nil {
			log.Printf("DAFUQ: 002\n")
			c.JSON(500, types.ErrorResponse{Message: "DAFUQ002"})
			return
		}
		token1 = token
	}
	c.SetCookie(consts.COOKIE_SESSION, token1.Val, 3600, "/", "localhost", false, true)
	acc := types.Account{
		ID:       dbUser.ID,
		Token:    token1.Val,
		Username: dbUser.Username,
	}
	c.JSON(200, acc)
}

// @Summary Login for user
// @Description Checks user in database as db.User and gives token
// @Tags user
// @Accept json
// @Produce json
// @Param user body types.User true "User info"
// @Success 200 {object} types.Account
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	dbUser := db.User{
		Username: user.Username,
		Password: user.Password,
	}
	foundUser := db.User{}
	dbc := db.Connect()
	if err := dbc.Find(&foundUser, dbUser).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if foundUser.ID == 0 {
		c.JSON(400, types.ErrorResponse{Message: "Credentials are incorrect"})
		return
	}
	var token1 *tokens.Token
	if token, err := tokens.AddToken(foundUser.ID); err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	} else {
		if token == nil {
			log.Printf("DAFUQ: 002\n")
			c.JSON(500, types.ErrorResponse{Message: "DAFUQ002"})
			return
		}
		token1 = token
	}
	c.SetCookie(consts.COOKIE_SESSION, token1.Val, 3600, "/", "localhost", false, true)
	acc := types.Account{
		ID:       foundUser.ID,
		Token:    token1.Val,
		Username: dbUser.Username,
	}
	c.JSON(200, acc)
}

func isSessionTokenForUserInvalid(userID uint) {}
