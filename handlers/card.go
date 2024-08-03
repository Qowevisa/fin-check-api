package handlers

import (
	"fmt"
	"strconv"

	"git.qowevisa.me/Qowevisa/gonuts/db"
	"git.qowevisa.me/Qowevisa/gonuts/types"
	"github.com/gin-gonic/gin"
)

// @Summary Get card by id
// @Description Get card by id
// @Tags card
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param card path int true "id"
// @Success 200 {object} types.DbCard
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /card/:id [get]
func CardGetId(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if idVal, err := strconv.ParseUint(idStr, 10, 32); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	} else {
		id = uint(idVal)
	}

	var dbCard db.Card
	dbc := db.Connect()
	if err := dbc.Find(&dbCard, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbCard.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	card := types.DbCard{
		Name:           dbCard.Name,
		Value:          dbCard.Value,
		HaveCreditLine: dbCard.HaveCreditLine,
		CreditLine:     dbCard.CreditLine,
	}
	c.JSON(200, card)
}

// @Summary Get card by id
// @Description Get card by id
// @Tags card
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param card body types.DbCard true "Card"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /card/add [post]
func CardAdd(c *gin.Context) {
	var card types.DbCard
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	dbCard := &db.Card{
		Name:           card.Name,
		Value:          card.Value,
		HaveCreditLine: card.HaveCreditLine,
		CreditLine:     card.CreditLine,
	}
	dbc := db.Connect()
	if err := dbc.Create(&dbCard).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbCard.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ004"})
		return
	}
	msg := types.Message{
		Message: fmt.Sprintf("Card with id %d was successfully created!", dbCard.ID),
	}
	c.JSON(200, msg)
}
