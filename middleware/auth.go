package middleware

import (
	"errors"
	"log"
	"net/http"

	"git.qowevisa.me/Qowevisa/fin-check-api/consts"
	"git.qowevisa.me/Qowevisa/fin-check-api/tokens"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

// Passes UserID with `c.Set("UserID")` as it gets id from token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(consts.COOKIE_SESSION)
		if errors.Is(err, http.ErrNoCookie) {
			c.JSON(401, types.ErrorResponse{Message: "Authorization cookie is required"})
			c.Abort()
			return
		}
		if !tokens.ValidateSessionToken(token) {
			c.JSON(401, types.ErrorResponse{Message: "Invalid authorization cookie"})
			c.Abort()
			return
		}
		session, err := tokens.GetSession(token)
		if err != nil {
			log.Printf("ERROR: tokens.GetSession: %v\n", err)
			c.JSON(500, types.ErrorResponse{Message: "Server error"})
			c.Abort()
			return
		}
		c.Set("UserID", session.UserID)

		c.Next()
	}
}
