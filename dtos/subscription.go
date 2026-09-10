package dtos

type CreateSubscriptionRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateSubscriptionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
