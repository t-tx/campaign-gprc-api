package repositories

import (
	"context"
	"fmt"
	"time"
)

// Campaign represents a campaign record in the database
type Campaign struct {
	ID        string    `json:"id" db:"id"`
	ValidFrom time.Time `json:"valid_from" db:"valid_from"`
	ValidTo   time.Time `json:"valid_to" db:"valid_to"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Slot      int32     `json:"slot" db:"slot"`
}

// InsertCampaign adds a new campaign to the database
func (r *Repository) InsertCampaign(ctx context.Context, campaign *Campaign) error {
	query := `INSERT INTO campaigns (id, valid_from, valid_to, slot)
			  VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(ctx, query, campaign.ID, campaign.ValidFrom, campaign.ValidTo, campaign.Slot)
	if err != nil {
		return fmt.Errorf("failed to insert campaign: %w", err)
	}
	return nil
}

// GetCampaign retrieves a single campaign by ID
func (r *Repository) GetCampaign(ctx context.Context, id string) (*Campaign, error) {
	query := `SELECT id, valid_from, valid_to, slot, created_at FROM campaigns WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)

	campaign := &Campaign{}
	err := row.Scan(&campaign.ID, &campaign.ValidFrom, &campaign.ValidTo, &campaign.Slot, &campaign.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign: %w", err)
	}
	return campaign, nil
}

// ListCampaigns retrieves all campaigns from the database
func (r *Repository) ListCampaigns(ctx context.Context) ([]*Campaign, error) {
	query := `SELECT id, url, created_at FROM campaigns`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []*Campaign
	for rows.Next() {
		var campaign Campaign
		if err := rows.Scan(&campaign.ID, &campaign.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		campaigns = append(campaigns, &campaign)
	}
	return campaigns, nil
}
