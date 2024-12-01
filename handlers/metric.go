package handlers

import (
	"git.qowevisa.me/Qowevisa/fin-check-api/db"
	"git.qowevisa.me/Qowevisa/fin-check-api/types"
	"github.com/gin-gonic/gin"
)

var metricTransform func(inp *db.Metric) types.DbMetric = func(inp *db.Metric) types.DbMetric {
	return types.DbMetric{
		Type:  inp.Type,
		Name:  inp.Name,
		Short: inp.Short,
	}
}

// @Summary Get all metrics for user
// @Description Get all metrics for user
// @Tags type
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} []types.DbMetric
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Security ApiKeyAuth
// @Router /metric/all [get]
func MetricGetAll(c *gin.Context) {
	dbc := db.Connect()
	var entities []*db.Metric
	if err := dbc.Find(&entities).Error; err != nil {
		c.JSON(500, types.ErrorResponse{Message: err.Error()})
		return
	}

	var ret []types.DbMetric
	for _, entity := range entities {
		ret = append(ret, metricTransform(entity))
	}
	c.JSON(200, ret)
}
