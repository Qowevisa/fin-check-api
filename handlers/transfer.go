package handlers

import (
	"fmt"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var transferTransform func(inp *db.Transfer) types.DbTransfer = func(inp *db.Transfer) types.DbTransfer {
	var fromCard types.DbCard
	var toCard types.DbCard
	if inp.FromCard != nil {
		fromCard = cardTransform(inp.FromCard)
	}
	if inp.ToCard != nil {
		toCard = cardTransform(inp.ToCard)
	}
	haveDiffCurrs := false
	if inp.FromCard != nil && inp.FromCard.Currency != nil && inp.ToCard != nil && inp.ToCard.Currency != nil {
		haveDiffCurrs = inp.FromCard.CurrencyID != inp.ToCard.CurrencyID
	}
	var showValue string
	if haveDiffCurrs {
		showValue = fmt.Sprintf("%d.%02d%s -> %d.%02d%s",
			inp.FromValue/100,
			inp.FromValue%100,
			inp.FromCard.Currency.Symbol,
			inp.ToValue/100,
			inp.ToValue%100,
			inp.ToCard.Currency.Symbol,
		)
	} else {
		showValue = fmt.Sprintf("%d.%02d", inp.Value/100, inp.Value%100)
	}
	return types.DbTransfer{
		ID:         inp.ID,
		FromCardID: inp.FromCardID,
		ToCardID:   inp.ToCardID,
		Value:      inp.Value,
		FromValue:  inp.FromValue,
		ToValue:    inp.ToValue,
		Date:       inp.Date,
		//
		ShowValue:               showValue,
		HaveDifferentCurrencies: haveDiffCurrs,
		FromCard:                fromCard,
		ToCard:                  toCard,
	}
}

// @Summary Get transfer by id
// @Description Get transfer by id
// @Tags transfer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param transfer path int true "id"
// @Success 200 {object} types.DbTransfer
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /transfer/:id [get]
func TransferGetId(c *gin.Context) {
	GetHandler(transferTransform)(c)
}

// @Summary Get all transfers for user
// @Description Get all transfers for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbTransfer
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /transfer/all [get]
func TransferGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.Transfer
	if err := dbc.Preload("FromCard.Currency").Preload("ToCard.Currency").Find(&entities, db.Transfer{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbTransfer
	for _, entity := range entities {
		ret = append(ret, transferTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Add transfer
// @Description Add transfer
// @Tags transfer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param transfer body types.DbTransfer true "Transfer"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /transfer/add [post]
func TransferAdd(c *gin.Context) {
	CreateHandler(&db.Transfer{}, func(src types.DbTransfer, dst *db.Transfer) {
		dst.FromCardID = src.FromCardID
		dst.ToCardID = src.ToCardID
		dst.Value = src.Value
		dst.Date = src.Date
		dst.FromValue = src.FromValue
		dst.ToValue = src.ToValue
	})(c)
}

// @Summary Edit transfer by id
// @Description Edit transfer by id
// @Tags transfer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param transferID path int true "id"
// @Param transfer body types.DbTransfer true "Transfer"
// @Success 200 {object} types.DbTransfer
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /transfer/edit/:id [put]
func TransferPutId(c *gin.Context) {
	UpdateHandler(
		func(src types.DbTransfer, dst *db.Transfer) {
			dst.FromCardID = src.FromCardID
			dst.ToCardID = src.ToCardID
			dst.Value = src.Value
			dst.Date = src.Date
			dst.FromValue = src.FromValue
			dst.ToValue = src.ToValue
		},
		transferTransform,
	)(c)
}

// @Summary Delete transfer by id
// @Description Delete transfer by id
// @Tags transfer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param transferID path int true "id"
// @Success 200 {object} types.DbTransfer
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /transfer/delete/:id [delete]
func TransferDeleteId(c *gin.Context) {
	DeleteHandler[*db.Transfer]()(c)
}
