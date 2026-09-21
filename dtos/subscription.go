package dtos

import "time"

type SubscriptionRequest struct {
	Name                   string     `json:"name" binding:"required,max=50"`
	Cost                   float64    `json:"cost" binding:"gte=0"`
	Type                   string     `json:"type" binding:"required,oneof=monthly yearly"`
	Category               string     `json:"category" binding:"required,oneof=streaming music productivity technology"`
	NextBillingDate        time.Time  `json:"nextBillingDate" binding:"required"`
	ReminderTimeInAdvanced int64      `json:"reminderTimeInAdvanced" binding:"required,gte=1"`
	FtEndDate              *time.Time `json:"ftEndDate"`
	Status                 string     `json:"status" binding:"omitempty,oneof=active free_trial inactive"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active free_trial inactive"`
}

type SubscriptionResponse struct {
	ID                     string     `json:"id"`
	Name                   string     `json:"name"`
	Cost                   float64    `json:"cost"`
	Type                   string     `json:"type"`
	Category               string     `json:"category"`
	NextBillingDate        time.Time  `json:"nextBillingDate"`
	ReminderTimeInAdvanced int64      `json:"reminderTimeInAdvanced"`
	FtEndDate              *time.Time `json:"ftEndDate"`
	Status                 string     `json:"status"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}
