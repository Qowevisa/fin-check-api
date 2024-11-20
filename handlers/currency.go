package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var currencyTransform func(inp *db.Currency) types.DbCurrency = func(inp *db.Currency) types.DbCurrency {
	return types.DbCurrency{
		ID:      inp.ID,
		Name:    inp.Name,
		ISOName: inp.ISOName,
		Symbol:  inp.Symbol,
	}
}

// @Summary Get all currencies for user
// @Description Get all currencies for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbCurrency
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /currency/all [get]
func CurrencyGetAll(c *gin.Context) {
	dbc := db.Connect()
	var entities []*db.Currency
	if err := dbc.Find(&entities).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbCurrency
	for _, entity := range entities {
		ret = append(ret, currencyTransform(entity))
	}
	c.JSON(200, ret)
}
