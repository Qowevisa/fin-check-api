package handlers

import (
	"fmt"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var categoryTransform func(*db.Category) types.DbCategory = func(inp *db.Category) types.DbCategory {
	nameWithParent := inp.Name
	if inp.Parent != nil {
		nameWithParent = fmt.Sprintf("%s -> %s", inp.Parent.Name, inp.Name)
	}
	return types.DbCategory{
		ID:             inp.ID,
		Name:           inp.Name,
		ParentID:       inp.ParentID,
		NameWithParent: nameWithParent,
	}
}

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
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/:id [get]
func CategoryGetId(c *gin.Context) {
	GetHandler(categoryTransform)(c)
}

// @Summary Get all categories for user
// @Description Get all categories for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbCategory
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/all [get]
func CategoryGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.Category
	if err := dbc.Preload("Parent").Find(&entities, db.Category{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbCategory
	for _, entity := range entities {
		ret = append(ret, categoryTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Add category
// @Description Add category
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
	CreateHandler(&db.Category{}, func(src types.DbCategory, dst *db.Category) {
		dst.Name = src.Name
		dst.ParentID = src.ParentID
	})(c)
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
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/edit/:id [put]
func CategoryPutId(c *gin.Context) {
	UpdateHandler(
		func(src types.DbCategory, dst *db.Category) {
			dst.Name = src.Name
			dst.ParentID = src.ParentID
		},
		categoryTransform,
	)(c)
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
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/delete/:id [delete]
func CategoryDeleteId(c *gin.Context) {
	DeleteHandler[*db.Category]()(c)
}
