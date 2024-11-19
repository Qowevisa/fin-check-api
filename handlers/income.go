package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var incomeTransform func(inp *db.Income) types.DbIncome = func(inp *db.Income) types.DbIncome {
	return types.DbIncome{
		ID:      inp.ID,
		CardID:  inp.CardID,
		Comment: inp.Comment,
		Value:   inp.Value,
		Date:    inp.Date,
	}
}

// @Summary Get income by id
// @Description Get income by id
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param income path int true "id"
// @Success 200 {object} types.DbIncome
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /income/:id [get]
func IncomeGetId(c *gin.Context) {
	GetHandler(incomeTransform)(c)
}

// @Summary Get all incomes for user
// @Description Get all incomes for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbIncome
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /income/all [get]
func IncomeGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.Income
	if err := dbc.Find(&entities, db.Income{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbIncome
	for _, entity := range entities {
		ret = append(ret, incomeTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Add income
// @Description Add income
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param income body types.DbIncome true "Income"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /income/add [post]
func IncomeAdd(c *gin.Context) {
	CreateHandler(
		&db.Income{},
		func(src types.DbIncome, dst *db.Income) {
			dst.CardID = src.CardID
			dst.Value = src.Value
			dst.Comment = src.Comment
			dst.Date = src.Date
		},
	)(c)
}

// @Summary Edit income by id
// @Description Edit income by id
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param incomeID path int true "id"
// @Param income body types.DbIncome true "Income"
// @Success 200 {object} types.DbIncome
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /income/edit/:id [put]
func IncomePutId(c *gin.Context) {
	UpdateHandler(
		func(src types.DbIncome, dst *db.Income) {
			dst.CardID = src.CardID
			dst.Value = src.Value
			dst.Comment = src.Comment
			dst.Date = src.Date
		},
		incomeTransform,
	)(c)
}

// @Summary Delete income by id
// @Description Delete income by id
// @Tags income
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param incomeID path int true "id"
// @Success 200 {object} types.DbIncome
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /income/delete/:id [delete]
func IncomeDeleteId(c *gin.Context) {
	DeleteHandler[*db.Income]()(c)
}
