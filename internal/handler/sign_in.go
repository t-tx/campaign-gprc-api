package handler

import (
	"campaign/internal/defined"
	"campaign/internal/pkg"
	"campaign/proto/generate/servicepb"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (c *campaignHandler) SignIn(ctx context.Context, req *servicepb.SignInRequest) (reply *servicepb.SignInReply, err error) {
	user, err := c.db.GetUser(ctx, req.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	var userData = map[string]string{defined.CAMPAIGNID_KEY: user.CampaignID}

	token, err := pkg.GenerateJWT(user.Username, userData, c.jwtSecret)
	if err != nil {
		return nil, errors.New("generate token error")
	}

	return &servicepb.SignInReply{
		Token: token,
	}, nil

}
