package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var transferTransform func(inp *db.Transfer) types.DbTransfer = func(inp *db.Transfer) types.DbTransfer {
	return types.DbTransfer{
		ID:         inp.ID,
		FromCardID: inp.FromCardID,
		ToCardID:   inp.ToCardID,
		Value:      inp.Value,
		Date:       inp.Date,
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
	if err := dbc.Find(&entities, db.Transfer{UserID: userID}).Error; err != nil {
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
