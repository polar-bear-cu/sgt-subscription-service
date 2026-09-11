package usecases

import (
	"context"
	"time"

	"github.com/polar-bear-cu/sgt-subscription-service/models"
	"github.com/polar-bear-cu/sgt-subscription-service/repositories"
)

type SubscriptionUsecase struct {
	repo repositories.SubscriptionRepository
}

func NewSubscription(repo repositories.SubscriptionRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{repo: repo}
}

func (u *SubscriptionUsecase) Create(ctx context.Context, userID, name string, billingDate time.Time) (models.Subscription, error) {
	return u.repo.Create(ctx, models.Subscription{UserID: userID, Name: name, BillingDate: billingDate})
}

func (u *SubscriptionUsecase) ListByUser(ctx context.Context, userID string) ([]models.Subscription, error) {
	return u.repo.ListByUser(ctx, userID)
}

func (u *SubscriptionUsecase) GetUpcomingForBilling(ctx context.Context, within time.Duration) ([]models.Subscription, error) {
	return u.repo.ListUpcomingForBilling(ctx, within)
}
