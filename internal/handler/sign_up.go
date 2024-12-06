package handler

import (
	"campaign/internal/db/repositories"
	"campaign/proto/generate/servicepb"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// SignUp implements servicepb.CampaignServiceServer.
func (c *campaignHandler) SignUp(ctx context.Context, req *servicepb.SignUpRequest) (*servicepb.SignUpReply, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("internal error")
	}

	err = c.db.InsertUser(ctx, &repositories.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		CampaignID:   req.CampaignId,
	})

	if err != nil {
		return nil, errors.New("duplicate username")
	}
	return &servicepb.SignUpReply{}, nil
}
