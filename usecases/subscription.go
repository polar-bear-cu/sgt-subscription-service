package usecases

import (
	"github.com/polar-bear-cu/sgt-subscription-service/models"
	"github.com/polar-bear-cu/sgt-subscription-service/repositories"
)

type SubscriptionUsecase struct {
	repo repositories.SubscriptionRepository
}

func NewSubscription(repo repositories.SubscriptionRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{repo: repo}
}

func (u *SubscriptionUsecase) Create(name string) (models.Subscription, error) {
	return u.repo.Create(models.Subscription{Name: name})
}
