package dtos

import "time"

type CreateSubscriptionRequest struct {
	Name        string    `json:"name" binding:"required"`
	BillingDate time.Time `json:"billingDate" binding:"required"`
}

type CreateSubscriptionResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	BillingDate time.Time `json:"billingDate"`
}
