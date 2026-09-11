package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/polar-bear-cu/sgt-subscription-service/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s models.Subscription) (models.Subscription, error)
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

// Methods
func (r *SubscriptionPostgres) Create(ctx context.Context, s models.Subscription) (models.Subscription, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO subscriptions (user_id, name, billing_date) VALUES ($1, $2, $3) RETURNING id`,
		s.UserID, s.Name, s.BillingDate,
	).Scan(&s.ID)

	return s, err
}

func (r *SubscriptionPostgres) ListByUser(ctx context.Context, userID string) ([]models.Subscription, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, billing_date FROM subscriptions WHERE user_id = $1`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []models.Subscription
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.BillingDate); err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	return out, rows.Err()
}

func (r *SubscriptionPostgres) ListUpcomingForBilling(ctx context.Context, within time.Duration) ([]models.Subscription, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, billing_date FROM subscriptions
		 WHERE billing_date BETWEEN now() AND $1`,
		time.Now().Add(within),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Subscription
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.BillingDate); err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	return out, rows.Err()
}
