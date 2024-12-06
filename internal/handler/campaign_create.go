package handler

import (
	"campaign/internal/db/repositories"
	"campaign/internal/defined"
	"campaign/internal/pkg"
	"campaign/proto/generate/servicepb"
	"context"
	"errors"
)

// CreateCampaign implements servicepb.CampaignServiceServer.
func (c *campaignHandler) CreateCampaign(ctx context.Context, req *servicepb.CreateCampaignRequest) (*servicepb.CreateCampaignReply, error) {
	newId, err := pkg.GenerateUniqueID(35)
	if err != nil {
		return nil, errors.New("internal error")
	}
	from := req.ValidFrom.AsTime()
	to := req.ValidTo.AsTime()
	_, err = c.cache.Set(ctx, defined.GetCampaignRemainingCacheKey(newId), req.Slot, to.Sub(from)).Result()
	if err != nil {
		return nil, errors.New("internal cache error")
	}

	err = c.db.InsertCampaign(ctx, &repositories.Campaign{
		ID:        newId,
		ValidFrom: req.ValidFrom.AsTime(),
		ValidTo:   req.ValidTo.AsTime(),
		Slot:      req.Slot,
	})
	if err != nil {
		return nil, errors.New("internal db error")
	}

	return &servicepb.CreateCampaignReply{
		Id: newId,
	}, nil
}
