package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-subscription-service/dtos"
	"github.com/polar-bear-cu/sgt-subscription-service/usecases"
)

type SubscriptionController struct {
	uc *usecases.SubscriptionUsecase
}

func NewSubscriptionController(uc *usecases.SubscriptionUsecase) *SubscriptionController {
	return &SubscriptionController{uc: uc}
}

func (ctl *SubscriptionController) CreateSubscription(c *gin.Context) {
	var req dtos.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	sub, err := ctl.uc.Create(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dtos.CreateSubscriptionResponse{ID: sub.ID, Name: sub.Name})
}
