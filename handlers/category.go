package handlers

import (
	"git.qowevisa.me/Qowevisa/gonuts/db"
	"git.qowevisa.me/Qowevisa/gonuts/types"
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
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /category/:id [get]
func CategoryGetId(c *gin.Context) {
	GetHandler(func(inp *db.Category) types.DbCategory {
		return types.DbCategory{
			ID:       inp.ID,
			Name:     inp.Name,
			ParentID: inp.ParentID,
		}
	})(c)
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
		func(inp *db.Category) types.DbCategory {
			return types.DbCategory{
				ID:       inp.ID,
				Name:     inp.Name,
				ParentID: inp.ParentID,
			}
		},
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
