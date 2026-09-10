package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-subscription-service/dtos"
)

// GetHealth godoc
// @Summary  health check
// @Tags     ops
// @Produce  json
// @Success  200 {object} dtos.HealthResponse
// @Router   /health [get]
func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, dtos.HealthResponse{
		Status:      "ok",
		Timestamp:   time.Now(),
		ServiceName: "subscription",
	})
}
