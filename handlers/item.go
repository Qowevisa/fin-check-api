package handlers

import (
	"log"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var itemTransform func(inp *db.Item) types.DbItem = func(inp *db.Item) types.DbItem {
	return types.DbItem{
		ID:             inp.ID,
		CategoryID:     inp.CategoryID,
		CurrentPriceID: inp.CurrentPriceID,
		TypeID:         inp.TypeID,
		Name:           inp.Name,
		Comment:        inp.Comment,
		MetricType:     inp.MetricType,
		MetricValue:    inp.MetricValue,
		Proteins:       inp.Proteins,
		Carbs:          inp.Carbs,
		Fats:           inp.Fats,
		Price:          inp.Price,
	}
}

// @Summary Get item by id
// @Description Get item by id
// @Tags item
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param item path int true "id"
// @Success 200 {object} types.DbItem
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /item/:id [get]
func ItemGetId(c *gin.Context) {
	GetHandler(itemTransform)(c)
}

// @Summary Get all items for user
// @Description Get all items for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbItem
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /item/all [get]
func ItemGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.Item
	if err := dbc.Find(&entities, db.Item{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbItem
	for _, entity := range entities {
		ret = append(ret, itemTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Get all items for user filtered
// @Description Get all items for user based on body criteria
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbItem
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /item/filter [post]
func ItemPostFilter(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	var filterObj types.DbItemSearch
	if err := c.ShouldBindJSON(&filterObj); err != nil {
		log.Printf("err is %v\n", err)
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	dbc := db.Connect()
	var entities []*db.Item
	filter := db.Item{
		UserID:     userID,
		CategoryID: filterObj.CategoryID,
		TypeID:     filterObj.TypeID,
	}

	if err := dbc.Find(&entities, filter).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbItem
	for _, entity := range entities {
		ret = append(ret, itemTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Delete item by id
// @Description Delete item by id
// @Tags item
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param itemID path int true "id"
// @Success 200 {object} types.DbItem
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /item/delete/:id [delete]
func ItemDeleteId(c *gin.Context) {
	DeleteHandler[*db.Item]()(c)
}
