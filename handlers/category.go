package handlers

import (
	"fmt"
	"strconv"

	"git.qowevisa.me/Qowevisa/gonuts/db"
	"git.qowevisa.me/Qowevisa/gonuts/types"
	"git.qowevisa.me/Qowevisa/gonuts/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get category by id
// @Description Get category by id
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param category path int true "id"
// @Success 200 {object} types.DbCategory
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/:id [get]
func CategoryGetId(c *gin.Context) {
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

	var dbCategory db.Category
	dbc := db.Connect()
	if err := dbc.Find(&dbCategory, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbCategory.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbCategory.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This category.id is not yours, you sneaky."})
		return
	}

	category := types.DbCategory{
		Name:     dbCategory.Name,
		ParentID: dbCategory.ParentID,
		UserID:   userID,
	}
	c.JSON(200, category)
}

// @Summary Get category by id
// @Description Get category by id
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param category body types.DbCategory true "Category"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/add [post]
func CategoryAdd(c *gin.Context) {
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

	var category types.DbCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	dbCategory := &db.Category{
		Name:     category.Name,
		ParentID: category.ParentID,
		UserID:   userID,
	}
	dbc := db.Connect()
	if err := dbc.Create(&dbCategory).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbCategory.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ004"})
		return
	}
	msg := types.Message{
		Message: fmt.Sprintf("Category with id %d was successfully created!", dbCategory.ID),
	}
	c.JSON(200, msg)
}

// @Summary Edit category by id
// @Description Edit category by id
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param categoryID path int true "id"
// @Param category body types.DbCategory true "Category"
// @Success 200 {object} types.DbCategory
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/edit/:id [put]
func CategoryPutId(c *gin.Context) {
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

	var dbCategory db.Category
	dbc := db.Connect()
	if err := dbc.Find(&dbCategory, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbCategory.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbCategory.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This category.id is not yours, you sneaky."})
		return
	}
	var category types.DbCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	utils.MergeNonZeroFields(category, dbCategory)

	if err := dbc.Save(dbCategory).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	ret := types.DbCategory{
		ID:       dbCategory.ID,
		Name:     dbCategory.Name,
		ParentID: dbCategory.ParentID,
		UserID:   dbCategory.UserID,
	}
	c.JSON(200, ret)
}

// @Summary Delete category by id
// @Description Delete category by id
// @Tags category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param categoryID path int true "id"
// @Success 200 {object} types.DbCategory
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/delete/:id [delete]
func CategoryDeleteId(c *gin.Context) {
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

	var dbCategory db.Category
	dbc := db.Connect()
	if err := dbc.Find(&dbCategory, id).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	if dbCategory.ID == 0 {
		c.JSON(500, types.ErrorResponse{Message: "DAFUQ003"})
		return
	}
	if dbCategory.UserID != userID {
		c.JSON(401, types.ErrorResponse{Message: "This category.id is not yours, you sneaky."})
		return
	}
	if err := dbc.Delete(dbCategory).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	ret := types.DbCategory{
		ID:       dbCategory.ID,
		Name:     dbCategory.Name,
		ParentID: dbCategory.ParentID,
		UserID:   dbCategory.UserID,
	}
	c.JSON(200, ret)
}
