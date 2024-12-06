package handler

import (
	"campaign/internal/pkg"
	"campaign/proto/generate/servicepb"
	"context"
	"errors"
)

func (c *campaignHandler) GetUserSubscription(ctx context.Context, req *servicepb.GetUserSubscriptionRequest) (*servicepb.GetUserSubscriptionReply, error) {
	auth, isLogined := pkg.GetAuthContext(ctx)
	if !isLogined {
		return nil, errors.New("unauthorized")
	}
	output, err := c.db.GetUserSubscription(ctx, auth.Username)
	if err != nil {
		return nil, err
	}
	return &servicepb.GetUserSubscriptionReply{
		SubscriptionId: output.SubscriptionId,
		Status:         output.Status,
	}, nil
}
