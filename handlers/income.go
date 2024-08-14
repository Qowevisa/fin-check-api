package handlers

import (
	"fmt"
	"strconv"

	"git.qowevisa.me/Qowevisa/gonuts/db"
	"git.qowevisa.me/Qowevisa/gonuts/types"
	"git.qowevisa.me/Qowevisa/gonuts/utils"
	"github.com/gin-gonic/gin"
)

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
	userIDAny, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 001"})
		return
	}

	var userID uint
	if userIDVal, ok := userIDAny.(uint); !ok {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 002"})
		return
	} else {
		userID = userIDVal
	}

	idStr := c.Param("id")
	var id uint
	if idVal, err := strconv.ParseUint(idStr, 10, 32); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	} else {
		id = uint(idVal)
	}

	var dbIncome db.Income
	dbc := db.Connect()
	if err := dbc.Find(&dbIncome, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbIncome.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbIncome.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This income.id is not yours, you sneaky."})
		return
	}

	ret := types.DbIncome{
		ID:      dbIncome.ID,
		CardID:  dbIncome.CardID,
		Comment: dbIncome.Comment,
		Value:   dbIncome.Value,
		Date:    dbIncome.Date,
		UserID:  dbIncome.UserID,
	}
	c.JSON(200, ret)
}

// @Summary Get income by id
// @Description Get income by id
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
	userIDAny, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 001"})
		return
	}

	var userID uint
	if userIDVal, ok := userIDAny.(uint); !ok {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 002"})
		return
	} else {
		userID = userIDVal
	}

	var income types.DbIncome
	if err := c.ShouldBindJSON(&income); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}
	if income.UserID != 0 && userID != income.UserID {
		c.JSON(403, types.ErrorResponse{Message: "UserID in body is different than yours!"})
	}
	if income.UserID == 0 {
		income.UserID = userID
	}

	dbIncome := &db.Income{
		CardID:  income.CardID,
		Comment: income.Comment,
		Value:   income.Value,
		Date:    income.Date,
		UserID:  income.UserID,
	}
	dbc := db.Connect()
	if err := dbc.Create(&dbIncome).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbIncome.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ004"})
		return
	}
	msg := types.Message{
		Message: fmt.Sprintf("Income with id %d was successfully created!", dbIncome.ID),
	}
	c.JSON(200, msg)
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
	userIDAny, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 001"})
		return
	}

	var userID uint
	if userIDVal, ok := userIDAny.(uint); !ok {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 002"})
		return
	} else {
		userID = userIDVal
	}

	idStr := c.Param("id")
	var id uint
	if idVal, err := strconv.ParseUint(idStr, 10, 32); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	} else {
		id = uint(idVal)
	}

	var dbIncome db.Income
	dbc := db.Connect()
	if err := dbc.Find(&dbIncome, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbIncome.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbIncome.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This income.id is not yours, you sneaky."})
		return
	}
	var income types.DbIncome
	if err := c.ShouldBindJSON(&income); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	utils.MergeNonZeroFields(income, dbIncome)

	if err := dbc.Save(dbIncome).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	ret := types.DbIncome{
		ID:      dbIncome.ID,
		CardID:  dbIncome.CardID,
		Comment: dbIncome.Comment,
		Value:   dbIncome.Value,
		Date:    dbIncome.Date,
		UserID:  dbIncome.UserID,
	}
	c.JSON(200, ret)
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
	userIDAny, exists := c.Get("UserID")
	if !exists {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 001"})
		return
	}

	var userID uint
	if userIDVal, ok := userIDAny.(uint); !ok {
		c.JSON(500, types.ErrorResponse{Message: "Internal error 002"})
		return
	} else {
		userID = userIDVal
	}

	idStr := c.Param("id")
	var id uint
	if idVal, err := strconv.ParseUint(idStr, 10, 32); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	} else {
		id = uint(idVal)
	}

	var dbIncome db.Income
	dbc := db.Connect()
	if err := dbc.Find(&dbIncome, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbIncome.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbIncome.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This income.id is not yours, you sneaky."})
		return
	}
	if err := dbc.Delete(dbIncome).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	ret := types.DbIncome{
		ID:      dbIncome.ID,
		CardID:  dbIncome.CardID,
		Comment: dbIncome.Comment,
		Value:   dbIncome.Value,
		Date:    dbIncome.Date,
		UserID:  dbIncome.UserID,
	}
	c.JSON(200, ret)
}
