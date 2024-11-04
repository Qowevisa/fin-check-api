package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

// @Summary Get debt by id
// @Description Get debt by id
// @Tags debt
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param debt path int true "id"
// @Success 200 {object} types.DbDebt
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /debt/:id [get]
func DebtGetId(c *gin.Context) {
	GetHandler(func(inp *db.Debt) types.DbDebt {
		return types.DbDebt{
			ID:       inp.ID,
			CardID:   inp.CardID,
			Comment:  inp.Comment,
			Value:    inp.Value,
			IOwe:     inp.IOwe,
			Date:     inp.Date,
			DateEnd:  inp.DateEnd,
			Finished: inp.Finished,
		}
	})(c)
}

// @Summary Get debt by id
// @Description Get debt by id
// @Tags debt
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param debt body types.DbDebt true "Debt"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /debt/add [post]
func DebtAdd(c *gin.Context) {
	debt := &db.Debt{}
	CreateHandler(debt,
		func(src types.DbDebt, dst *db.Debt) {
			dst.CardID = src.CardID
			dst.Comment = src.Comment
			dst.Value = src.Value
			dst.IOwe = src.IOwe
			dst.Date = src.Date
			dst.DateEnd = src.DateEnd
			dst.Finished = src.Finished
		},
	)(c)
}

// @Summary Edit debt by id
// @Description Edit debt by id
// @Tags debt
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param debtID path int true "id"
// @Param debt body types.DbDebt true "Debt"
// @Success 200 {object} types.DbDebt
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /debt/edit/:id [put]
func DebtPutId(c *gin.Context) {
	UpdateHandler(
		// Filter used to apply only needed changes from srt to dst before updating dst
		func(src types.DbDebt, dst *db.Debt) {
			dst.CardID = src.CardID
			dst.Comment = src.Comment
			dst.Value = src.Value
			dst.IOwe = src.IOwe
			dst.Date = src.Date
			dst.DateEnd = src.DateEnd
			dst.Finished = src.Finished
		},
		func(inp *db.Debt) types.DbDebt {
			return types.DbDebt{
				ID:       inp.ID,
				CardID:   inp.CardID,
				Comment:  inp.Comment,
				Value:    inp.Value,
				IOwe:     inp.IOwe,
				Date:     inp.Date,
				DateEnd:  inp.DateEnd,
				Finished: inp.Finished,
			}
		})(c)
}

// @Summary Delete debt by id
// @Description Delete debt by id
// @Tags debt
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param debtID path int true "id"
// @Success 200 {object} types.DbDebt
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /debt/delete/:id [delete]
func DebtDeleteId(c *gin.Context) {
	DeleteHandler[*db.Debt]()(c)
}
