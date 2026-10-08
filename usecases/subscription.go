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
	loc  *time.Location
}

func NewSubscription(repo repositories.SubscriptionRepository, loc *time.Location) *SubscriptionUsecase {
	return &SubscriptionUsecase{repo: repo, loc: loc}
}

func (u *SubscriptionUsecase) Create(ctx context.Context, userID string, in SubscriptionInput) (models.Subscription, error) {
	if err := validateInput(&in); err != nil {
		return models.Subscription{}, err
	}
	s := toModel(userID, in)
	s.BillingDay = in.NextBillingDate.In(u.loc).Day()
	return u.repo.Create(ctx, s)
}

func (u *SubscriptionUsecase) GetByID(ctx context.Context, userID, id string) (models.Subscription, error) {
	return u.repo.FindByID(ctx, id, userID)
}

func (u *SubscriptionUsecase) Update(ctx context.Context, userID, id string, in SubscriptionInput) (models.Subscription, error) {
	if err := validateInput(&in); err != nil {
		return models.Subscription{}, err
	}
	current, err := u.repo.FindByID(ctx, id, userID)
	if err != nil {
		return models.Subscription{}, err
	}

	s := toModel(userID, in)
	s.ID = id
	// Keep the anchor when the date is unchanged: a date clamped to Feb 28 must still return to the 31st.
	s.BillingDay = current.BillingDay
	if !in.NextBillingDate.Equal(current.NextBillingDate) {
		s.BillingDay = in.NextBillingDate.In(u.loc).Day()
	}
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

const (
	defaultPageLimit = 10
	maxPageLimit     = 100
)

type ListQuery struct {
	Name     string
	Category string
	Status   string
	Type     string
	SortBy   string
	Order    string
	Page     int
	Limit    int
}

type ListResult struct {
	Items      []models.Subscription
	Page       int
	Limit      int
	Total      int
	TotalPages int
}

func (u *SubscriptionUsecase) List(ctx context.Context, userID string, q ListQuery) (ListResult, error) {
	page := max(q.Page, 1)
	limit := q.Limit
	if limit <= 0 {
		limit = defaultPageLimit
	}
	limit = min(limit, maxPageLimit)

	items, total, err := u.repo.List(ctx, userID, repositories.ListParams{
		Name:     q.Name,
		Category: q.Category,
		Status:   q.Status,
		Type:     q.Type,
		SortBy:   q.SortBy,
		Order:    q.Order,
		Limit:    limit,
		Offset:   (page - 1) * limit,
	})
	if err != nil {
		return ListResult{}, err
	}

	return ListResult{
		Items:      items,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + limit - 1) / limit,
	}, nil
}

func (u *SubscriptionUsecase) Summary(ctx context.Context, userID string) (models.SubscriptionSummary, error) {
	return u.repo.Summary(ctx, userID)
}

func (u *SubscriptionUsecase) ListByUser(ctx context.Context, userID string) ([]models.Subscription, error) {
	return u.repo.ListByUser(ctx, userID)
}

func (u *SubscriptionUsecase) ListDueReminders(ctx context.Context, date string) ([]models.DueReminder, error) {
	return u.repo.ListDueReminders(ctx, date)
}

// AdvanceBillingDates turns ended free trials into active records and go to the first cycle date
func (u *SubscriptionUsecase) AdvanceBillingDates(ctx context.Context, date string) (advanced, converted int64, err error) {
	today, err := time.ParseInLocation(time.DateOnly, date, u.loc)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: date must be YYYY-MM-DD", ErrValidation)
	}

	converted, err = u.repo.ConvertEndedTrials(ctx, date)
	if err != nil {
		return 0, 0, err
	}

	past, err := u.repo.ListPastBilling(ctx, date)
	if err != nil {
		return 0, converted, err
	}

	var errs error
	for _, s := range past {
		next := s.NextBillingDate
		for next.Before(today) {
			next = nextCycle(next, s.Type, s.BillingDay, u.loc)
		}
		ok, err := u.repo.AdvanceBillingDate(ctx, s.ID, s.NextBillingDate, next)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("subscription %s: %w", s.ID, err))
			continue
		}
		if ok {
			advanced++
		}
	}
	return advanced, converted, errs
}

// nextCycle rebuilds the date from billingDay
func nextCycle(cur time.Time, typ string, billingDay int, loc *time.Location) time.Time {
	l := cur.In(loc)
	months := time.Month(1)
	if typ == models.TypeYearly {
		months = 12
	}
	first := time.Date(l.Year(), l.Month()+months, 1, l.Hour(), l.Minute(), l.Second(), 0, loc)
	last := first.AddDate(0, 1, -1).Day()
	return time.Date(first.Year(), first.Month(), min(billingDay, last), l.Hour(), l.Minute(), l.Second(), 0, loc)
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
