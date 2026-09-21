package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/polar-bear-cu/sgt-subscription-service/models"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionRepository interface {
	Create(ctx context.Context, s models.Subscription) (models.Subscription, error)
	FindByID(ctx context.Context, id, userID string) (models.Subscription, error)
	Update(ctx context.Context, s models.Subscription) (models.Subscription, error)
	UpdateStatus(ctx context.Context, id, userID, status string) (models.Subscription, error)
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, userID string) ([]models.Subscription, error)
	ListUpcomingForBilling(ctx context.Context, within time.Duration) ([]models.Subscription, error)
}

type SubscriptionPostgres struct {
	db *pgxpool.Pool
}

// Constructor
func NewSubscriptionPostgres(db *pgxpool.Pool) SubscriptionRepository {
	return &SubscriptionPostgres{db: db}
}

const subscriptionColumns = `id, user_id, name, cost, type, category, next_billing_date,
	reminder_time_in_advanced, ft_end_date, status, created_at, updated_at`

// Methods
func (r *SubscriptionPostgres) Create(ctx context.Context, s models.Subscription) (models.Subscription, error) {
	row := r.db.QueryRow(ctx,
		`INSERT INTO subscriptions
			(user_id, name, cost, type, category, next_billing_date, reminder_time_in_advanced, ft_end_date, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING `+subscriptionColumns,
		s.UserID, s.Name, s.Cost, s.Type, s.Category, s.NextBillingDate,
		s.ReminderTimeInAdvanced, s.FtEndDate, s.Status,
	)
	return scanSubscription(row)
}

func (r *SubscriptionPostgres) FindByID(ctx context.Context, id, userID string) (models.Subscription, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+subscriptionColumns+` FROM subscriptions WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	return scanSubscription(row)
}

func (r *SubscriptionPostgres) Update(ctx context.Context, s models.Subscription) (models.Subscription, error) {
	row := r.db.QueryRow(ctx,
		`UPDATE subscriptions
		 SET name = $3, cost = $4, type = $5, category = $6, next_billing_date = $7,
		     reminder_time_in_advanced = $8, ft_end_date = $9, status = $10, updated_at = now()
		 WHERE id = $1 AND user_id = $2
		 RETURNING `+subscriptionColumns,
		s.ID, s.UserID, s.Name, s.Cost, s.Type, s.Category, s.NextBillingDate,
		s.ReminderTimeInAdvanced, s.FtEndDate, s.Status,
	)
	return scanSubscription(row)
}

func (r *SubscriptionPostgres) UpdateStatus(ctx context.Context, id, userID, status string) (models.Subscription, error) {
	row := r.db.QueryRow(ctx,
		`UPDATE subscriptions SET status = $3, updated_at = now()
		 WHERE id = $1 AND user_id = $2
		 RETURNING `+subscriptionColumns,
		id, userID, status,
	)
	return scanSubscription(row)
}

func (r *SubscriptionPostgres) Delete(ctx context.Context, id, userID string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSubscriptionNotFound
	}
	return nil
}

func (r *SubscriptionPostgres) ListByUser(ctx context.Context, userID string) ([]models.Subscription, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+subscriptionColumns+` FROM subscriptions WHERE user_id = $1 ORDER BY next_billing_date`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSubscriptions(rows)
}

func (r *SubscriptionPostgres) ListUpcomingForBilling(ctx context.Context, within time.Duration) ([]models.Subscription, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+subscriptionColumns+` FROM subscriptions
		 WHERE status IN ($1, $2)
		   AND next_billing_date BETWEEN now() AND $3
		 ORDER BY next_billing_date`,
		models.StatusActive, models.StatusFreeTrial, time.Now().Add(within),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSubscriptions(rows)
}

// Helpers
func scanSubscription(row pgx.Row) (models.Subscription, error) {
	var s models.Subscription
	err := row.Scan(
		&s.ID, &s.UserID, &s.Name, &s.Cost, &s.Type, &s.Category, &s.NextBillingDate,
		&s.ReminderTimeInAdvanced, &s.FtEndDate, &s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Subscription{}, ErrSubscriptionNotFound
	}
	return s, err
}

func scanSubscriptions(rows pgx.Rows) ([]models.Subscription, error) {
	out := make([]models.Subscription, 0)
	for rows.Next() {
		s, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
