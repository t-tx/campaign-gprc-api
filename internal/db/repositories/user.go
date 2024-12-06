package repositories

import (
	"context"
	"fmt"
	"time"
)

// User represents a user record in the database
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	CampaignID   string    `json:"campaign_id" db:"campaign_id"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// InsertUser adds a new user to the database
func (r *Repository) InsertUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, password_hash, campaign_id) 
			  VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, user.Username, user.PasswordHash, user.CampaignID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

// GetUser retrieves a single user by ID
func (r *Repository) GetUser(ctx context.Context, username string) (*User, error) {

	query := `SELECT username, password_hash, campaign_id,  created_at FROM users WHERE username=$1`
	row := r.db.QueryRow(ctx, query, username)

	user := &User{}
	err := row.Scan(&user.Username, &user.PasswordHash, &user.CampaignID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// new feed
// ListUsers retrieves all users from the database
