package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var cardTransform func(inp *db.Card) types.DbCard = func(inp *db.Card) types.DbCard {
	return types.DbCard{
		ID:             inp.ID,
		Name:           inp.Name,
		Balance:        inp.Balance,
		HaveCreditLine: inp.HaveCreditLine,
		CreditLine:     inp.CreditLine,
	}
}

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
	GetHandler(cardTransform)(c)
}

// @Summary Get all cards for user
// @Description Get all cards for user
// @Tags card
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbCard
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /card/all [get]
func CardGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.Card
	if err := dbc.Find(&entities, db.Card{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbCard
	for _, entity := range entities {
		ret = append(ret, cardTransform(entity))
	}
	c.JSON(200, ret)
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
	card := &db.Card{}
	CreateHandler(card, func(src types.DbCard, dst *db.Card) {
		dst.Name = src.Name
		dst.Balance = src.Balance
		dst.HaveCreditLine = src.HaveCreditLine
		dst.CreditLine = src.CreditLine
	})(c)
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
	UpdateHandler(
		// Filter used to apply only needed changes from srt to dst before updating dst
		func(src types.DbCard, dst *db.Card) {
			dst.Name = src.Name
			dst.Balance = src.Balance
			dst.CreditLine = src.CreditLine
			dst.HaveCreditLine = src.HaveCreditLine
		},
		cardTransform)(c)
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
	DeleteHandler[*db.Card]()(c)
}
