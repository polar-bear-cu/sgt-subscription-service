package dtos

type CreateSubscriptionRequest struct {
	Name string `json:"name" binding:"required"`
}

type SubscriptionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
