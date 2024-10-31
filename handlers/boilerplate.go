package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ERROR_UTIL_USERIDER1 = errors.New("internal error 001: UserID not found")
	ERROR_UTIL_USERIDER2 = errors.New("internal error 002: UserID conversion error")
)

func GetUserID(c *gin.Context) (uint, error) {
	userIDAny, exists := c.Get("UserID")
	if !exists {
		return 0, ERROR_UTIL_USERIDER1
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		return 0, ERROR_UTIL_USERIDER2
	}
	return userID, nil
}

//

// Parses id as uint
func ParseID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	idVal, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid request: %w", err)
	}
	return uint(idVal), nil
}
