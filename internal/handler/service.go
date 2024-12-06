package handler

import (
	"campaign/internal/db/repositories"
	"campaign/proto/generate/servicepb"
	"context"

	"github.com/go-redis/redis/v8"
)

type campaignHandler struct {
	servicepb.UnimplementedCampaignServiceServer
	db        *repositories.Repository
	cache     *redis.Client
	jwtSecret []byte
}

func New(db *repositories.Repository, cache *redis.Client, jwtSecret []byte) *campaignHandler {
	s := &campaignHandler{
		db:        db,
		cache:     cache,
		jwtSecret: jwtSecret,
	}
	go s.simulateQueueWorker(context.Background())
	return s
}
