package handler

import (
	"campaign/internal/defined"
	"campaign/internal/pkg"
	"campaign/proto/generate/servicepb"
	"context"
	"errors"
)

func (c *campaignHandler) UserInfo(ctx context.Context, req *servicepb.UserInfoRequest) (reply *servicepb.UserInfoReply, err error) {
	auth, isLogined := pkg.GetAuthContext(ctx)
	if !isLogined {
		return nil, errors.New("unauthorized")
	}

	return &servicepb.UserInfoReply{
		Username:   auth.Username,
		CampaignId: auth.Payload[defined.CAMPAIGNID_KEY],
	}, nil
}
