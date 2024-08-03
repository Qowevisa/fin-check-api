package middleware

import (
	"git.qowevisa.me/Qowevisa/gonuts/tokens"
	"git.qowevisa.me/Qowevisa/gonuts/types"
	"github.com/gin-gonic/gin"
)

// Passes UserID with `c.Set("UserID")` as it gets id from token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, types.ErrorResponse{Message: "Authorization header is required"})
			c.Abort()
			return
		}

		token := authHeader
		if !tokens.AmIAllowed(token) {
			c.JSON(401, types.ErrorResponse{Message: "Token is invalid"})
			c.Abort()
			return
		}

		if userID, err := tokens.GetID(token); err != nil {
			c.JSON(401, types.ErrorResponse{Message: "Token is invalid ERR4001"})
			c.Abort()
			return
		} else {
			c.Set("UserID", userID)
		}

		c.Next()
	}
}
