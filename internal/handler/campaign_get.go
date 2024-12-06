package handler

import (
	"campaign/internal/defined"
	"campaign/proto/generate/servicepb"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	a = 1
)

// GetCampaign implements servicepb.CampaignServiceServer.
func (c *campaignHandler) GetCampaign(ctx context.Context, req *servicepb.GetCampaignRequest) (*servicepb.GetCampaignReply, error) {
	campaign, err := c.db.GetCampaign(ctx, req.Id)
	if err != nil {
		return nil, errors.New("campaign not found")
	}

	v, err := c.cache.Get(ctx, defined.GetCampaignRemainingCacheKey(req.Id)).Result()

	remaining, err := strconv.Atoi(v)
	if err != nil {
		remaining = int(campaign.Slot)
	}

	var status string
	now := time.Now()

	switch {
	case now.Before(campaign.ValidFrom):
		status = "upcoming"
	case now.After(campaign.ValidTo):
		status = "expired"
	default:
		status = "active"
	}
	return &servicepb.GetCampaignReply{
		Url:       fmt.Sprintf("https://example.com/campaign/%s", campaign.ID),
		ValidFrom: timestamppb.New(campaign.ValidFrom),
		ValidTo:   timestamppb.New(campaign.ValidTo),
		Slot:      campaign.Slot,
		Status:    status,
		Remaining: int32(remaining),
	}, nil
}
