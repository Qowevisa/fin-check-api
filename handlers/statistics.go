package handlers

import (
	"fmt"
	"log"
	"slices"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"git.qowevisa.me/Qowevisa/fin-check-api/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get all statisticss for user
// @Description Get all statisticss for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.StatsTypeCurrencyChart
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /statistics/type [get]
func StatisticsGetAllSpendingsForTypes(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var settingsTypeFilter []*db.SettingsTypeFilter
	if err := dbc.Find(&settingsTypeFilter, db.SettingsTypeFilter{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	var filerTypeIDs []uint
	for _, typeFilter := range settingsTypeFilter {
		filerTypeIDs = append(filerTypeIDs, typeFilter.TypeID)
	}
	var userTypes []*db.Type
	if err := dbc.Find(&userTypes, db.Type{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	var userExpenses []*db.Expense
	if err := dbc.Not(map[string]interface{}{"type_id": filerTypeIDs}).Preload("Card.Currency").Find(&userExpenses, db.Expense{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	currToChart := make(map[uint][]*db.Expense)

	for _, expense := range userExpenses {
		if expense.Card == nil || expense.Card.Currency == nil {
			log.Printf("ERROR: db.Preload DID NOT WORKED OUT!\n")
			c.JSON(500, types.ErrorResponse{Message: "Internal error. E.S.T.1"})
			return
		}
		if val, exists := currToChart[expense.Card.CurrencyID]; !exists {
			currToChart[expense.Card.CurrencyID] = []*db.Expense{}
			currToChart[expense.Card.CurrencyID] = append(currToChart[expense.Card.CurrencyID], expense)
		} else {
			currToChart[expense.Card.CurrencyID] = append(val, expense)
		}
	}
	var ret []types.StatsTypeCurrencyChart
	for _, expenseArray := range currToChart {
		if expenseArray[0].Card == nil || expenseArray[0].Card.Currency == nil {
			log.Printf("ERROR: db.Preload DID NOT WORKED OUT!\n")
			c.JSON(500, types.ErrorResponse{Message: "Internal error. E.S.T.2"})
			return
		}
		currency := expenseArray[0].Card.Currency
		typeToValue := make(map[uint]types.StatsType)

		for _, expense := range expenseArray {
			if val, exists := typeToValue[expense.TypeID]; !exists {
				idx := slices.IndexFunc(userTypes, func(t *db.Type) bool { return t.ID == expense.TypeID })
				typeForChart := userTypes[idx]
				typeToValue[expense.TypeID] = types.StatsType{
					Value: expense.Value,
					Name:  typeForChart.Name,
					Color: typeForChart.Color,
				}
			} else {
				val.Value += expense.Value
				typeToValue[expense.TypeID] = val
			}
		}
		var sum uint64 = 0
		var elements []types.StatsType
		for _, val := range typeToValue {
			elements = append(elements, val)
			sum += val.Value
		}
		slices.SortFunc(elements, func(a, b types.StatsType) int {
			return utils.DescendingSort(a.Value, b.Value)
		})

		ret = append(ret, types.StatsTypeCurrencyChart{
			CurrencyLabel: fmt.Sprintf("%s (%s)", currency.Symbol, currency.ISOName),
			Elements:      elements,
			Total:         sum,
		})
	}
	c.JSON(200, ret)
}
