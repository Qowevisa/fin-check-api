package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var expenseTransform func(inp *db.Expense) types.DbExpense = func(inp *db.Expense) types.DbExpense {
	return types.DbExpense{
		ID:      inp.ID,
		CardID:  inp.CardID,
		TypeID:  inp.TypeID,
		Value:   inp.Value,
		Comment: inp.Comment,
		Date:    inp.Date,
	}
}

// @Summary Get expense by id
// @Description Get expense by id
// @Tags expense
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param expense path int true "id"
// @Success 200 {object} types.DbExpense
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /expense/:id [get]
func ExpenseGetId(c *gin.Context) {
	GetHandler(expenseTransform)(c)
}

// @Summary Get all expenses for user
// @Description Get all expenses for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbExpense
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /expense/all [get]
func ExpenseGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.Expense
	if err := dbc.Find(&entities, db.Expense{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbExpense
	for _, entity := range entities {
		ret = append(ret, expenseTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Add expense by id
// @Description Add expense by id
// @Tags expense
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param expense body types.DbExpense true "Expense"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /expense/add [post]
func ExpenseAdd(c *gin.Context) {
	CreateHandler(&db.Expense{}, func(src types.DbExpense, dst *db.Expense) {
		dst.CardID = src.CardID
		dst.TypeID = src.TypeID
		dst.Value = src.Value
		dst.Comment = src.Comment
		dst.Date = src.Date
	})(c)
}

// @Summary Edit expense by id
// @Description Edit expense by id
// @Tags expense
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param expenseID path int true "id"
// @Param expense body types.DbExpense true "Expense"
// @Success 200 {object} types.DbExpense
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /expense/edit/:id [put]
func ExpensePutId(c *gin.Context) {
	UpdateHandler(
		func(src types.DbExpense, dst *db.Expense) {
			dst.CardID = src.CardID
			dst.TypeID = src.TypeID
			dst.Value = src.Value
			dst.Comment = src.Comment
			dst.Date = src.Date
		},
		expenseTransform,
	)(c)
}

// @Summary Delete expense by id
// @Description Delete expense by id
// @Tags expense
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param expenseID path int true "id"
// @Success 200 {object} types.DbExpense
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /expense/delete/:id [delete]
func ExpenseDeleteId(c *gin.Context) {
	DeleteHandler[*db.Expense]()(c)
}
