package repositories

import (
	"sync"

	"github.com/google/uuid"
	"github.com/polar-bear-cu/sgt-subscription-service/models"
)

type SubscriptionRepository interface {
	Create(s models.Subscription) (models.Subscription, error)
}

type inMemorySubscription struct {
	mu    sync.Mutex
	items []models.Subscription
}

func NewInMemorySubscription() SubscriptionRepository {
	return &inMemorySubscription{}
}

func (r *inMemorySubscription) Create(s models.Subscription) (models.Subscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s.ID = uuid.NewString()
	r.items = append(r.items, s)
	return s, nil
}
