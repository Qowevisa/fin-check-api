package handlers

import (
	"log"

	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var settingsTypeFilterTransform func(inp *db.SettingsTypeFilter) types.SettingsTypeFilter = func(inp *db.SettingsTypeFilter) types.SettingsTypeFilter {
	return types.SettingsTypeFilter{
		TypeID:     inp.TypeID,
		FilterThis: true,
	}
}

// @Summary Get all settingstypefilters for user
// @Description Get all settingstypefilters for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.SettingsTypeFilter
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /settings/type/all [get]
func SettingsTypeFilterGetAll(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	dbc := db.Connect()
	var entities []*db.SettingsTypeFilter
	if err := dbc.Find(&entities, db.SettingsTypeFilter{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.SettingsTypeFilter
	for _, entity := range entities {
		ret = append(ret, settingsTypeFilterTransform(entity))
	}
	c.JSON(200, ret)
}

// @Summary Get all settingstypefilters for user
// @Description Get all settingstypefilters for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.SettingsTypeFilter
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /settings/type/update [put]
func SettingsTypePutBatch(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	var updates []types.SettingsTypeFilter
	if err := c.ShouldBindJSON(&updates); err != nil {
		log.Printf("err is %v\n", err)
		c.JSON(400, types.ErrorResponse{Message: "Invalid request"})
		return
	}

	dbc := db.Connect()
	var entities []*db.SettingsTypeFilter
	if err := dbc.Find(&entities, db.SettingsTypeFilter{UserID: userID}).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}
	var removeFilters []*db.SettingsTypeFilter
	for _, typeFilter := range entities {
		skipTypeFilter := false
		for _, updateTypeFilter := range updates {
			if typeFilter.TypeID == updateTypeFilter.TypeID {
				if updateTypeFilter.FilterThis == false {
					removeFilters = append(removeFilters, typeFilter)
				}
				skipTypeFilter = true
				break
			}
		}
		if skipTypeFilter {
			continue
		}
	}
	for _, potentiallyNewFilter := range updates {
		wasItHandled := false
		for _, handled := range removeFilters {
			if handled.TypeID == potentiallyNewFilter.TypeID {
				wasItHandled = true
				break
			}
		}
		if wasItHandled {
			continue
		}
		newFilter := &db.SettingsTypeFilter{
			TypeID: potentiallyNewFilter.TypeID,
			UserID: userID,
		}
		if err := dbc.Create(newFilter).Error; err != nil {
			log.Printf("dbc.Create error: %v\n", err)
			continue
		}
	}

	for _, settingsDelete := range removeFilters {
		if err := dbc.Unscoped().Delete(settingsDelete).Error; err != nil {
			log.Printf("dbc.Delete:ERROR: %v\n", err)
			c.JSON(500, types.ErrorResponse{Message: err.Error()})
			return
		}
	}
	c.JSON(200, types.Message{Info: "It was successfull"})
}
