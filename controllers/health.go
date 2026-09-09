package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-subscription-service/dtos"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, dtos.HealthResponse{Status: "ok"})
}
