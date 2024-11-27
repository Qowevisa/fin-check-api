package handlers

import (
	"fmt"
	"log"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var expenseTransform func(inp *db.Expense) types.DbExpense = func(inp *db.Expense) types.DbExpense {
	var card types.DbCard
	expenseValueSymbolPostfix := ""
	if inp.Card != nil {
		card = cardTransform(inp.Card)
		if inp.Card.Currency != nil {
			expenseValueSymbolPostfix = fmt.Sprintf(" (%s)", inp.Card.Currency.Symbol)
		}
	} else {
		card = types.DbCard{}
	}
	var showValue string
	showValue = fmt.Sprintf("%d.%02d%s", inp.Value/100, inp.Value%100, expenseValueSymbolPostfix)
	return types.DbExpense{
		ID:        inp.ID,
		CardID:    inp.CardID,
		TypeID:    inp.TypeID,
		Value:     inp.Value,
		Comment:   inp.Comment,
		Date:      inp.Date,
		Card:      card,
		ShowValue: showValue,
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
	if err := dbc.Preload("Card.Currency").Find(&entities, db.Expense{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbExpense
	for _, entity := range entities {
		ret = append(ret, expenseTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Add expense
// @Description Add expense
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

// @Summary Add many expenses
// @Description Add expense by propagating main struct to every child in children field
// @Tags expense
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param expense body types.DbExpenseBulk true "Expense"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /expense/add [post]
func ExpenseBulkCreate(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var u types.DbExpenseBulk
	if err := c.ShouldBindJSON(&u); err != nil {
		log.Printf("err is %v\n", err)
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}
	if u.IsEveryFieldPropagated() {
		c.JSON(400, types.ErrorResponse{Message: "You can't just try to propagate every field for children."})
		return
	}
	he := &db.Helper_ExpenseBulk{
		PropagateCardID:  u.PropagateCardID,
		CardID:           u.CardID,
		PropagateTypeID:  u.PropagateTypeID,
		TypeID:           u.TypeID,
		PropagateValue:   u.PropagateValue,
		Value:            u.Value,
		PropagateComment: u.PropagateComment,
		Comment:          u.Comment,
		PropagateDate:    u.PropagateDate,
		Date:             u.Date,
		UserID:           userID,
	}

	var expenses []*db.Expense
	for _, child := range u.Children {
		c := db.Expense{
			CardID:  child.CardID,
			TypeID:  child.TypeID,
			Value:   child.Value,
			Comment: child.Comment,
			Date:    child.Date,
		}
		expenses = append(expenses, he.CreateExpenseFromChild(c))
	}

	dbc := db.Connect()
	var whatToRollback []*db.Expense
	shouldRollback := false
	defer func() {
		if shouldRollback {
			for _, e := range whatToRollback {
				if err := dbc.Delete(e).Error; err != nil {
					log.Printf("dbc.Delete ERROR: %v\n", err)
				}
			}
		}
	}()
	for _, entity := range expenses {
		if err := dbc.Create(entity).Error; err != nil {
			shouldRollback = true
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}
		whatToRollback = append(whatToRollback, entity)
	}

	c.JSON(200, types.Message{Info: fmt.Sprintf("%d entities were created successfully!", len(expenses))})
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
