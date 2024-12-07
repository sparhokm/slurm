package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sparhokm/slurm/subscription/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	conn *pgxpool.Pool
}

func NewSubscriptionRepository(conn *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{conn: conn}
}

func (r *SubscriptionRepository) Create(ctx context.Context, s *models.Subscription) error {
	_, err := r.conn.Exec(
		ctx,
		"INSERT INTO subscriptions (user_id, prefix, files_owner_id, created_at) VALUES ($1, $2, $3, $4)",
		s.UserID, s.Prefix, s.FilesOwnerID, time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepository) FindUserSubscribers(
	ctx context.Context,
	filepath string,
) ([]models.SubscriptionShort, error) {
	rows, err := r.conn.Query(
		ctx,
		`SELECT DISTINCT user_id, files_owner_id FROM subscriptions WHERE $1 LIKE prefix || '%'`,
		filepath,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.SubscriptionShort
	for rows.Next() {
		sub := models.SubscriptionShort{}
		err := rows.Scan(&sub.UserID, &sub.FilesOwnerID)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		subs = append(subs, sub)
	}

	return subs, nil
}
