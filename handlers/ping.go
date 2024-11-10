package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

// @Summary Ping
// @Description Pong
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} types.Message
// @Router /ping [get]
func PingGet(c *gin.Context) {
	c.JSON(200, types.Message{Info: "Pong!"})
}
