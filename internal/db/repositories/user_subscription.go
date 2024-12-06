package repositories

import (
	"context"
	"fmt"
	"time"
)

// UserSubscription represents a record in the UserSubscriptions table.
type UserSubscription struct {
	Id             string    `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Status         string    `json:"status" db:"status"`
	CampaignId     string    `json:"campaign_id" db:"campaign_id"`
	SubscriptionId string    `json:"subscription_id" db:"subscription_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// InsertUserSubscription adds a new user-subscription association to the database.
func (r *Repository) InsertUserSubscription(ctx context.Context, us *UserSubscription) error {
	query := `INSERT INTO UserSubscriptions (id, username, subscription_id, status, campaign_id) 
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, us.Id, us.Username, us.SubscriptionId, us.Status, us.CampaignId)
	if err != nil {
		return fmt.Errorf("failed to insert user subscription: %w", err)
	}
	return nil
}

// UpdateUserSubscription updates an existing user-subscription association in the database.
func (r *Repository) UpdateUserSubscription(ctx context.Context, userSubscription *UserSubscription) error {
	query := `UPDATE UserSubscriptions 
			  SET username=$1, subscription_id=$2 
			  WHERE id=$3`
	_, err := r.db.Exec(ctx, query, userSubscription.Username, userSubscription.SubscriptionId, userSubscription.Id)
	if err != nil {
		return fmt.Errorf("failed to update user subscription: %w", err)
	}
	return nil
}

// GetUserSubscription retrieves a single user-subscription association by ID.
func (r *Repository) GetUserSubscription(ctx context.Context, username string) (*UserSubscription, error) {
	query := `SELECT id, username, subscription_id, campaign_id, status, created_at
			  FROM UserSubscriptions
			  WHERE username=$1`
	row := r.db.QueryRow(ctx, query, username)

	userSubscription := &UserSubscription{}
	err := row.Scan(&userSubscription.Id, &userSubscription.Username, &userSubscription.SubscriptionId, &userSubscription.CampaignId, &userSubscription.Status, &userSubscription.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user subscription: %w", err)
	}
	return userSubscription, nil
}
