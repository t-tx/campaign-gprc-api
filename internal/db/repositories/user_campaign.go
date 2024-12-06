package repositories

import (
	"context"
	"fmt"
	"time"
)

// UserCampaign represents a record in the UserCampaigns table.
type UserCampaign struct {
	ID         int       `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Status     string    `json:"status" db:"status"`
	CampaignID int       `json:"campaign_id" db:"campaign_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// InsertUserCampaign adds a new user-campaign association to the database.
func (r *Repository) InsertUserCampaign(ctx context.Context, userCampaign *UserCampaign) error {
	query := `INSERT INTO user_campaigns (username, campaign_id, created_at) 
			  VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(ctx, query, userCampaign.Username, userCampaign.CampaignID, userCampaign.CreatedAt).Scan(&userCampaign.ID)
	if err != nil {
		return fmt.Errorf("failed to insert user campaign: %w", err)
	}
	return nil
}

// UpdateUserCampaign updates an existing user-campaign association in the database.
func (r *Repository) UpdateUserCampaign(ctx context.Context, userCampaign *UserCampaign) error {
	query := `UPDATE user_campaigns 
			  SET username=$1, campaign_id=$2 
			  WHERE id=$3`
	_, err := r.db.Exec(ctx, query, userCampaign.Username, userCampaign.CampaignID, userCampaign.ID)
	if err != nil {
		return fmt.Errorf("failed to update user campaign: %w", err)
	}
	return nil
}

// GetUserCampaign retrieves a single user-campaign association by ID.
func (r *Repository) GetUserCampaign(ctx context.Context, username string) (*UserCampaign, error) {
	query := `SELECT id, username, campaign_id, created_at 
			  FROM user_campaigns 
			  WHERE username=$1`
	row := r.db.QueryRow(ctx, query, username)

	userCampaign := &UserCampaign{}
	err := row.Scan(&userCampaign.ID, &userCampaign.Username, &userCampaign.CampaignID, &userCampaign.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user campaign: %w", err)
	}
	return userCampaign, nil
}

// ListUserCampaigns retrieves all user-campaign associations from the database.
func (r *Repository) ListUserCampaigns(ctx context.Context) ([]UserCampaign, error) {
	query := `SELECT id, username, campaign_id, created_at 
			  FROM user_campaigns`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list user campaigns: %w", err)
	}
	defer rows.Close()

	var userCampaigns []UserCampaign
	for rows.Next() {
		var userCampaign UserCampaign
		if err := rows.Scan(&userCampaign.ID, &userCampaign.Username, &userCampaign.CampaignID, &userCampaign.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user campaign: %w", err)
		}
		userCampaigns = append(userCampaigns, userCampaign)
	}
	return userCampaigns, nil
}
