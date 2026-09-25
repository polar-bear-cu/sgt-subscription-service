package models

import "time"

const (
	TypeMonthly = "monthly"
	TypeYearly  = "yearly"
)

const (
	StatusActive    = "active"
	StatusFreeTrial = "free_trial"
	StatusInactive  = "inactive"
)

const (
	CategoryStreaming    = "streaming"
	CategoryMusic        = "music"
	CategoryProductivity = "productivity"
	CategoryTechnology   = "technology"
)

type Subscription struct {
	ID                     string
	UserID                 string
	Name                   string
	Cost                   float64
	Type                   string
	Category               string
	NextBillingDate        time.Time
	ReminderTimeInAdvanced int64
	FtEndDate              *time.Time
	Status                 string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type SubscriptionSummary struct {
	Count       int
	MonthlyCost float64
}
