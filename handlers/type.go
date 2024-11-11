package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var typeTransform func(inp *db.Type) types.DbType = func(inp *db.Type) types.DbType {
	return types.DbType{
		ID:      inp.ID,
		Name:    inp.Name,
		Comment: inp.Comment,
		Color:   inp.Color,
	}
}

// @Summary Get dbtype by id
// @Description Get dbtype by id
// @Tags dbtype
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param dbtype path int true "id"
// @Success 200 {object} types.DbType
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /dbtype/:id [get]
func TypeGetId(c *gin.Context) {
	GetHandler(typeTransform)(c)
}

// @Summary Get all types for user
// @Description Get all types for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbType
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /type/all [get]
func TypeGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.Type
	if err := dbc.Find(&entities, db.Type{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbType
	for _, entity := range entities {
		ret = append(ret, typeTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Get dbtype by id
// @Description Get dbtype by id
// @Tags dbtype
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param dbtype body types.DbType true "Type"
// @Success 200 {object} types.Message
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /dbtype/add [post]
func TypeAdd(c *gin.Context) {
	t := &db.Type{}
	CreateHandler(t, func(update types.DbType, dst *db.Type) {
		dst.Name = update.Name
		dst.Comment = update.Comment
		dst.Color = update.Color
	})(c)
}

// @Summary Edit dbtype by id
// @Description Edit dbtype by id
// @Tags dbtype
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param dbtypeID path int true "id"
// @Param dbtype body types.DbType true "Type"
// @Success 200 {object} types.DbType
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /dbtype/edit/:id [put]
func TypePutId(c *gin.Context) {
	UpdateHandler(
		func(updates types.DbType, dst *db.Type) {
			dst.Name = updates.Name
			dst.Comment = updates.Comment
			dst.Color = updates.Color
		},
		typeTransform,
	)(c)
}

// @Summary Delete dbtype by id
// @Description Delete dbtype by id
// @Tags dbtype
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param dbtypeID path int true "id"
// @Success 200 {object} types.DbType
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /dbtype/delete/:id [delete]
func TypeDeleteId(c *gin.Context) {
	DeleteHandler[*db.Type]()(c)
}
