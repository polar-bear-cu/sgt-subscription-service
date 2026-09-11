package models

import "time"

type Subscription struct {
	ID          string
	UserID      string
	Name        string
	BillingDate time.Time
}
