package middleware

import (
	"errors"
	"net/http"

	"git.qowevisa.me/Qowevisa/fin-check-api/consts"
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
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
		var session *db.Session
		if validated, tmpSession := tokens.ValidateAndGetSessionToken(token); !validated {
			c.JSON(401, types.ErrorResponse{Message: "Invalid authorization cookie"})
			c.Abort()
			return
		} else {
			session = tmpSession
		}
		c.Set("UserID", session.UserID)

		c.Next()
	}
}
