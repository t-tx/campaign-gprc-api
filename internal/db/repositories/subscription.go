package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Subscription represents a record in the Subscriptions table.
type Subscription struct {
	ID        string          `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Price     decimal.Decimal `json:"price" db:"price"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

// ListSubscriptions retrieves all subscriptions from the database.
func (r *Repository) ListSubscriptions(ctx context.Context) ([]*Subscription, error) {
	query := `SELECT id, name, price, created_at 
			  FROM subscriptions`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*Subscription
	for rows.Next() {
		var subscription Subscription
		if err := rows.Scan(&subscription.ID, &subscription.Name, &subscription.Price, &subscription.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}
	return subscriptions, nil
}

func (r *Repository) GetSubscription(ctx context.Context, id string) (*Subscription, error) {
	query := `SELECT id, name, price, created_at 
			  FROM subscriptions WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var subscription Subscription
	if err := row.Scan(&subscription.ID, &subscription.Name, &subscription.Price, &subscription.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to scan subscription: %w", err)
	}

	return &subscription, nil
}

func (r *Repository) ListSubscriptionsAsMap(ctx context.Context) (map[string]*Subscription, error) {
	query := `SELECT id, name, price, created_at 
			  FROM subscriptions`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions = make(map[string]*Subscription)
	for rows.Next() {
		var subscription Subscription
		if err := rows.Scan(&subscription.ID, &subscription.Name, &subscription.Price, &subscription.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions[subscription.ID] = &subscription
	}
	return subscriptions, nil
}
