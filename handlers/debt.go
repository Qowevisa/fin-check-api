package handlers

import (
	"fmt"
	"strconv"

	"git.qowevisa.me/Qowevisa/gonuts/db"
	"git.qowevisa.me/Qowevisa/gonuts/types"
	"git.qowevisa.me/Qowevisa/gonuts/utils"
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
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /debt/:id [get]
func DebtGetId(c *gin.Context) {
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

	var dbDebt db.Debt
	dbc := db.Connect()
	if err := dbc.Find(&dbDebt, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbDebt.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbDebt.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This debt.id is not yours, you sneaky."})
		return
	}

	ret := types.DbDebt{
		ID:       dbDebt.ID,
		CardID:   dbDebt.CardID,
		Comment:  dbDebt.Comment,
		Value:    dbDebt.Value,
		IOwe:     dbDebt.IOwe,
		Date:     dbDebt.Date,
		DateEnd:  dbDebt.DateEnd,
		Finished: dbDebt.Finished,
		UserID:   dbDebt.UserID,
	}
	c.JSON(200, ret)
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
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /debt/add [post]
func DebtAdd(c *gin.Context) {
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

	var debt types.DbDebt
	if err := c.ShouldBindJSON(&debt); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}
	if debt.UserID != 0 && userID != debt.UserID {
		c.JSON(403, types.ErrorResponse{Message: "UserID in body is different than yours!"})
	}
	if debt.UserID == 0 {
		debt.UserID = userID
	}

	dbDebt := &db.Debt{
		CardID:   debt.CardID,
		Comment:  debt.Comment,
		Value:    debt.Value,
		IOwe:     debt.IOwe,
		Date:     debt.Date,
		DateEnd:  debt.DateEnd,
		Finished: debt.Finished,
		UserID:   debt.UserID,
	}
	dbc := db.Connect()
	if err := dbc.Create(&dbDebt).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbDebt.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ004"})
		return
	}
	msg := types.Message{
		Message: fmt.Sprintf("Debt with id %d was successfully created!", dbDebt.ID),
	}
	c.JSON(200, msg)
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
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /debt/edit/:id [put]
func DebtPutId(c *gin.Context) {
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

	var dbDebt db.Debt
	dbc := db.Connect()
	if err := dbc.Find(&dbDebt, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbDebt.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbDebt.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This debt.id is not yours, you sneaky."})
		return
	}
	var debt types.DbDebt
	if err := c.ShouldBindJSON(&debt); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	utils.MergeNonZeroFields(debt, dbDebt)

	if err := dbc.Save(dbDebt).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	ret := types.DbDebt{
		ID:       dbDebt.ID,
		CardID:   dbDebt.CardID,
		Comment:  dbDebt.Comment,
		Value:    dbDebt.Value,
		IOwe:     dbDebt.IOwe,
		Date:     dbDebt.Date,
		DateEnd:  dbDebt.DateEnd,
		Finished: dbDebt.Finished,
		UserID:   dbDebt.UserID,
	}
	c.JSON(200, ret)
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
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /debt/delete/:id [delete]
func DebtDeleteId(c *gin.Context) {
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

	var dbDebt db.Debt
	dbc := db.Connect()
	if err := dbc.Find(&dbDebt, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbDebt.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbDebt.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This debt.id is not yours, you sneaky."})
		return
	}
	if err := dbc.Delete(dbDebt).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	ret := types.DbDebt{
		ID:       dbDebt.ID,
		CardID:   dbDebt.CardID,
		Comment:  dbDebt.Comment,
		Value:    dbDebt.Value,
		IOwe:     dbDebt.IOwe,
		Date:     dbDebt.Date,
		DateEnd:  dbDebt.DateEnd,
		Finished: dbDebt.Finished,
		UserID:   dbDebt.UserID,
	}
	c.JSON(200, ret)
}
