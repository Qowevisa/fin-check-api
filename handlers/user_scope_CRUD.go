package handlers

import (
	"fmt"
	"log"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

// GetHandler returns a generic handler for retrieving an entity by ID.
func GetHandler[T db.UserIdentifiable, R any](transform func(inp T) R) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserID(c)
		if err != nil {
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}

		id, err := ParseID(c)
		if err != nil {
			c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
			return
		}

		dbc := db.Connect()
		// NOTE: DO NOT fuck with this
		// only by pure luck with GORM First function we don't have SEGFAULT
		var entity T
		if err := dbc.First(&entity, id).Error; err != nil {
			c.JSON(400, types.ErrorResponse{Message: err.Error()})
			return
		}

		if entity.GetID() == 0 {
			c.JSON(401, types.ErrorResponse{Message: fmt.Sprintf("Entity with %d id was not found", id)})
			return
		}
		if entity.GetUserID() != userID {
			c.JSON(403, types.ErrorResponse{Message: fmt.Sprintf("This entity is not yours")})
			return
		}

		ret := transform(entity)

		c.JSON(200, ret)
	}
}

// CreateHandler returns a generic handler for creating an entity.
func CreateHandler[T db.UserIdentifiable, R any](entity T, applyChanges func(src R, dst T)) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserID(c)
		if err != nil {
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}

		var updates R
		if err := c.ShouldBindJSON(&updates); err != nil {
			log.Printf("err is %v\n", err)
			c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
			return
		}

		applyChanges(updates, entity)
		entity.SetUserID(userID)
		dbc := db.Connect()
		if err := dbc.Create(entity).Error; err != nil {
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(200, types.Message{Info: fmt.Sprintf("Entity created with ID %d", entity.GetID())})
	}
}

// UpdateHandler returns a generic handler for updating an entity.
func UpdateHandler[T db.UserIdentifiable, R any](
	// Filter used to apply only needed changes from srt to dst before updating dst
	filter func(src R, dst T),
	// Transform used to output only needed fields from inp to response
	transform func(inp T) R) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserID(c)
		if err != nil {
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}

		id, err := ParseID(c)
		if err != nil {
			c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
			return
		}

		dbc := db.Connect()
		// NOTE: DO NOT fuck with this
		// only by pure luck with GORM First function we don't have SEGFAULT
		var entity T
		if err := dbc.First(&entity, id).Error; err != nil {
			c.JSON(400, types.ErrorResponse{Message: err.Error()})
			return
		}

		if entity.GetID() == 0 {
			c.JSON(401, types.ErrorResponse{Message: fmt.Sprintf("Entity with %d id was not found", id)})
			return
		}
		if entity.GetUserID() != userID {
			c.JSON(403, types.ErrorResponse{Message: fmt.Sprintf("This entity is not yours")})
			return
		}

		var updates R
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
			return
		}

		filter(updates, entity)
		if err := dbc.Save(&entity).Error; err != nil {
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}

		ret := transform(entity)
		c.JSON(200, ret)
	}
}

// DeleteHandler returns a generic handler for deleting an entity by ID.
func DeleteHandler[T db.UserIdentifiable]() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserID(c)
		if err != nil {
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}

		id, err := ParseID(c)
		if err != nil {
			c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
			return
		}

		dbc := db.Connect()
		// NOTE: DO NOT fuck with this
		// only by pure luck with GORM First function we don't have SEGFAULT
		var entity T
		if err := dbc.First(&entity, id).Error; err != nil {
			c.JSON(400, types.ErrorResponse{Message: err.Error()})
			return
		}

		if entity.GetID() == 0 {
			c.JSON(401, types.ErrorResponse{Message: fmt.Sprintf("Entity with %d id was not found", id)})
			return
		}
		if entity.GetUserID() != userID {
			c.JSON(403, types.ErrorResponse{Message: fmt.Sprintf("This entity is not yours")})
			return
		}

		if err := dbc.Delete(entity).Error; err != nil {
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(200, types.Message{Info: fmt.Sprintf("Entity with ID %d deleted", entity.GetID())})
	}
}
