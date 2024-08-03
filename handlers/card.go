package handlers

import (
	"fmt"
	"strconv"

	"git.qowevisa.me/Qowevisa/gonuts/db"
	"git.qowevisa.me/Qowevisa/gonuts/types"
	"git.qowevisa.me/Qowevisa/gonuts/utils"
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
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /card/:id [get]
func CardGetId(c *gin.Context) {
	userIDAny, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 001"})
		return
	}

	var userID uint
	if userIDVal, ok := userIDAny.(uint); !ok {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 002"})
		return
	} else {
		userID = userIDVal
	}

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
	if dbCard.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This card.id is not yours, you sneaky."})
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
	userIDAny, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 001"})
		return
	}

	var userID uint
	if userIDVal, ok := userIDAny.(uint); !ok {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 002"})
		return
	} else {
		userID = userIDVal
	}

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
		UserID:         userID,
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

// @Summary Edit card by id
// @Description Edit card by id
// @Tags card
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param cardID path int true "id"
// @Param card body types.DbCard true "Card"
// @Success 200 {object} types.DbCard
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /card/edit/:id [put]
func CardPutId(c *gin.Context) {
	userIDAny, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 001"})
		return
	}

	var userID uint
	if userIDVal, ok := userIDAny.(uint); !ok {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 002"})
		return
	} else {
		userID = userIDVal
	}

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
	if dbCard.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This card.id is not yours, you sneaky."})
		return
	}
	var card types.DbCard
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	utils.MergeNonZeroFields(card, dbCard)

	if err := dbc.Save(dbCard).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	ret := types.DbCard{
		ID:             dbCard.ID,
		Name:           dbCard.Name,
		Value:          dbCard.Value,
		HaveCreditLine: dbCard.HaveCreditLine,
		CreditLine:     dbCard.CreditLine,
	}
	c.JSON(200, ret)
}

// @Summary Delete card by id
// @Description Delete card by id
// @Tags card
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param cardID path int true "id"
// @Success 200 {object} types.DbCard
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /card/delete/:id [delete]
func CardDeleteId(c *gin.Context) {
	userIDAny, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 001"})
		return
	}

	var userID uint
	if userIDVal, ok := userIDAny.(uint); !ok {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 002"})
		return
	} else {
		userID = userIDVal
	}

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
	if dbCard.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This card.id is not yours, you sneaky."})
		return
	}
	if err := dbc.Delete(dbCard).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	ret := types.DbCard{
		ID:             dbCard.ID,
		Name:           dbCard.Name,
		Value:          dbCard.Value,
		HaveCreditLine: dbCard.HaveCreditLine,
		CreditLine:     dbCard.CreditLine,
	}
	c.JSON(200, ret)
}
