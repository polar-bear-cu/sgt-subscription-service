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

// CreateSubscription godoc
// @Summary  create a subscription
// @Tags     subscriptions
// @Accept   json
// @Produce  json
// @Param    body  body      dtos.CreateSubscriptionRequest  true  "subscription"
// @Success  201   {object}  dtos.CreateSubscriptionResponse
// @Failure  400   {object}  map[string]string
// @Failure  500   {object}  map[string]string
// @Router   /api/v1/subscriptions [post]
func (ctl *SubscriptionController) CreateSubscription(c *gin.Context) {
	var req dtos.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	sub, err := ctl.uc.Create(c.Request.Context(), userID, req.Name, req.BillingDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dtos.CreateSubscriptionResponse{
		ID:          sub.ID,
		Name:        sub.Name,
		BillingDate: sub.BillingDate,
	})
}
