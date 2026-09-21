package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/polar-bear-cu/sgt-subscription-service/models"
	"github.com/polar-bear-cu/sgt-subscription-service/repositories"
)

var ErrValidation = errors.New("validation failed")

type SubscriptionInput struct {
	Name                   string
	Cost                   float64
	Type                   string
	Category               string
	NextBillingDate        time.Time
	ReminderTimeInAdvanced int64
	FtEndDate              *time.Time
	Status                 string
}

type SubscriptionUsecase struct {
	repo repositories.SubscriptionRepository
}

func NewSubscription(repo repositories.SubscriptionRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{repo: repo}
}

func (u *SubscriptionUsecase) Create(ctx context.Context, userID string, in SubscriptionInput) (models.Subscription, error) {
	if err := validateInput(&in); err != nil {
		return models.Subscription{}, err
	}
	return u.repo.Create(ctx, toModel(userID, in))
}

func (u *SubscriptionUsecase) GetByID(ctx context.Context, userID, id string) (models.Subscription, error) {
	return u.repo.FindByID(ctx, id, userID)
}

func (u *SubscriptionUsecase) Update(ctx context.Context, userID, id string, in SubscriptionInput) (models.Subscription, error) {
	if err := validateInput(&in); err != nil {
		return models.Subscription{}, err
	}
	s := toModel(userID, in)
	s.ID = id
	return u.repo.Update(ctx, s)
}

func (u *SubscriptionUsecase) UpdateStatus(ctx context.Context, userID, id, status string) (models.Subscription, error) {
	if status == models.StatusFreeTrial {
		current, err := u.repo.FindByID(ctx, id, userID)
		if err != nil {
			return models.Subscription{}, err
		}
		if current.FtEndDate == nil {
			return models.Subscription{}, fmt.Errorf("%w: ftEndDate is required before switching to free_trial", ErrValidation)
		}
	}
	return u.repo.UpdateStatus(ctx, id, userID, status)
}

func (u *SubscriptionUsecase) Delete(ctx context.Context, userID, id string) error {
	return u.repo.Delete(ctx, id, userID)
}

func (u *SubscriptionUsecase) ListByUser(ctx context.Context, userID string) ([]models.Subscription, error) {
	return u.repo.ListByUser(ctx, userID)
}

func (u *SubscriptionUsecase) GetUpcomingForBilling(ctx context.Context, within time.Duration) ([]models.Subscription, error) {
	return u.repo.ListUpcomingForBilling(ctx, within)
}

func validateInput(in *SubscriptionInput) error {
	if in.Status == "" {
		in.Status = models.StatusActive
	}
	if in.Status == models.StatusFreeTrial && in.FtEndDate == nil {
		return fmt.Errorf("%w: ftEndDate is required when status is free_trial", ErrValidation)
	}
	return nil
}

func toModel(userID string, in SubscriptionInput) models.Subscription {
	return models.Subscription{
		UserID:                 userID,
		Name:                   in.Name,
		Cost:                   in.Cost,
		Type:                   in.Type,
		Category:               in.Category,
		NextBillingDate:        in.NextBillingDate,
		ReminderTimeInAdvanced: in.ReminderTimeInAdvanced,
		FtEndDate:              in.FtEndDate,
		Status:                 in.Status,
	}
}
