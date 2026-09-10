package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/polar-bear-cu/sgt-subscription-service/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s models.Subscription) (models.Subscription, error)
	ListByUser(ctx context.Context, userID string) ([]models.Subscription, error)
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
		`INSERT INTO subscriptions (user_id, name) VALUES ($1, $2) RETURNING id`,
		s.UserID, s.Name,
	).Scan(&s.ID)

	return s, err
}

func (r *SubscriptionPostgres) ListByUser(ctx context.Context, userID string) ([]models.Subscription, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name FROM subscriptions WHERE user_id = $1`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []models.Subscription
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name); err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	return out, rows.Err()
}
