package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-subscription-service/dtos"
	"github.com/polar-bear-cu/sgt-subscription-service/models"
	"github.com/polar-bear-cu/sgt-subscription-service/repositories"
	"github.com/polar-bear-cu/sgt-subscription-service/usecases"
)

type SubscriptionController struct {
	uc *usecases.SubscriptionUsecase
}

func NewSubscriptionController(uc *usecases.SubscriptionUsecase) *SubscriptionController {
	return &SubscriptionController{uc: uc}
}

// Create godoc
// @Summary  create a subscription
// @Tags     subscriptions
// @Accept   json
// @Produce  json
// @Param    body  body      dtos.SubscriptionRequest  true  "subscription"
// @Success  201   {object}  dtos.SubscriptionResponse
// @Failure  400   {object}  map[string]string
// @Failure  401   {object}  map[string]string
// @Failure  500   {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/subscriptions [post]
func (ctl *SubscriptionController) Create(c *gin.Context) {
	var req dtos.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := ctl.uc.Create(c.Request.Context(), c.GetString("user_id"), toInput(req))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, toResponse(sub))
}

// List godoc
// @Summary  list my subscriptions (filter, sort, paginate)
// @Tags     subscriptions
// @Produce  json
// @Param    query  query     dtos.ListSubscriptionsQuery  false  "filters, sort and pagination"
// @Success  200    {object}  dtos.ListSubscriptionsResponse
// @Failure  400    {object}  map[string]string
// @Failure  401    {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/subscriptions [get]
func (ctl *SubscriptionController) List(c *gin.Context) {
	var q dtos.ListSubscriptionsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctl.uc.List(c.Request.Context(), c.GetString("user_id"), usecases.ListQuery{
		Name:     q.Name,
		Category: q.Category,
		Status:   q.Status,
		Type:     q.Type,
		SortBy:   q.SortBy,
		Order:    q.Order,
		Page:     q.Page,
		Limit:    q.Limit,
	})
	if err != nil {
		respondErr(c, err)
		return
	}

	items := make([]dtos.SubscriptionResponse, 0, len(res.Items))
	for _, s := range res.Items {
		items = append(items, toResponse(s))
	}
	c.JSON(http.StatusOK, dtos.ListSubscriptionsResponse{
		Items:      items,
		Page:       res.Page,
		Limit:      res.Limit,
		Total:      res.Total,
		TotalPages: res.TotalPages,
	})
}

// GetSummary godoc
// @Summary  count and monthly cost of my active / free-trial subscriptions (yearly / 12)
// @Tags     subscriptions
// @Produce  json
// @Success  200  {object}  dtos.SummaryResponse
// @Failure  401  {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/subscriptions/summary [get]
func (ctl *SubscriptionController) GetSummary(c *gin.Context) {
	sum, err := ctl.uc.Summary(c.Request.Context(), c.GetString("user_id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, dtos.SummaryResponse{
		Count:       sum.Count,
		MonthlyCost: sum.MonthlyCost,
	})
}

// GetByID godoc
// @Summary  get a subscription
// @Tags     subscriptions
// @Produce  json
// @Param    id   path      string  true  "subscription id"
// @Success  200  {object}  dtos.SubscriptionResponse
// @Failure  401  {object}  map[string]string
// @Failure  404  {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/subscriptions/{id} [get]
func (ctl *SubscriptionController) GetByID(c *gin.Context) {
	sub, err := ctl.uc.GetByID(c.Request.Context(), c.GetString("user_id"), c.Param("id"))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(sub))
}

// Update godoc
// @Summary  update a subscription (all fields)
// @Tags     subscriptions
// @Accept   json
// @Produce  json
// @Param    id    path      string                    true  "subscription id"
// @Param    body  body      dtos.SubscriptionRequest  true  "subscription"
// @Success  200   {object}  dtos.SubscriptionResponse
// @Failure  400   {object}  map[string]string
// @Failure  401   {object}  map[string]string
// @Failure  404   {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/subscriptions/{id} [put]
func (ctl *SubscriptionController) Update(c *gin.Context) {
	var req dtos.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := ctl.uc.Update(c.Request.Context(), c.GetString("user_id"), c.Param("id"), toInput(req))
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(sub))
}

// UpdateStatus godoc
// @Summary  change subscription status
// @Tags     subscriptions
// @Accept   json
// @Produce  json
// @Param    id    path      string                   true  "subscription id"
// @Param    body  body      dtos.UpdateStatusRequest  true  "new status"
// @Success  200   {object}  dtos.SubscriptionResponse
// @Failure  400   {object}  map[string]string
// @Failure  401   {object}  map[string]string
// @Failure  404   {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/subscriptions/{id} [patch]
func (ctl *SubscriptionController) UpdateStatus(c *gin.Context) {
	var req dtos.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := ctl.uc.UpdateStatus(c.Request.Context(), c.GetString("user_id"), c.Param("id"), req.Status)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(sub))
}

// Delete godoc
// @Summary  delete a subscription
// @Tags     subscriptions
// @Param    id   path  string  true  "subscription id"
// @Success  204
// @Failure  401  {object}  map[string]string
// @Failure  404  {object}  map[string]string
// @Security BearerAuth
// @Router   /api/v1/subscriptions/{id} [delete]
func (ctl *SubscriptionController) Delete(c *gin.Context) {
	if err := ctl.uc.Delete(c.Request.Context(), c.GetString("user_id"), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Helpers
func respondErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repositories.ErrSubscriptionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, usecases.ErrValidation):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		log.Printf("subscription: internal error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func toInput(p dtos.SubscriptionRequest) usecases.SubscriptionInput {
	return usecases.SubscriptionInput{
		Name:                   p.Name,
		Cost:                   p.Cost,
		Type:                   p.Type,
		Category:               p.Category,
		NextBillingDate:        p.NextBillingDate,
		ReminderTimeInAdvanced: p.ReminderTimeInAdvanced,
		FtEndDate:              p.FtEndDate,
		Status:                 p.Status,
	}
}

func toResponse(s models.Subscription) dtos.SubscriptionResponse {
	return dtos.SubscriptionResponse{
		ID:                     s.ID,
		Name:                   s.Name,
		Cost:                   s.Cost,
		Type:                   s.Type,
		Category:               s.Category,
		NextBillingDate:        s.NextBillingDate,
		ReminderTimeInAdvanced: s.ReminderTimeInAdvanced,
		FtEndDate:              s.FtEndDate,
		Status:                 s.Status,
		CreatedAt:              s.CreatedAt,
		UpdatedAt:              s.UpdatedAt,
	}
}
