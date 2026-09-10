package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/polar-bear-cu/sgt-subscription-service/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s models.Subscription) (models.Subscription, error)
}

type SubscriptionPostgres struct {
	db *pgxpool.Pool
}

func NewSubscriptionPostgres(db *pgxpool.Pool) SubscriptionRepository {
	return &SubscriptionPostgres{db: db}
}

func (r *SubscriptionPostgres) Create(ctx context.Context, s models.Subscription) (models.Subscription, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO subscriptions (user_id, name) VALUES ($1, $2) RETURNING id`,
		s.UserID, s.Name,
	).Scan(&s.ID)
	return s, err
}
