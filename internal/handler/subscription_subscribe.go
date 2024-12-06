package handler

import (
	"campaign/internal/db/repositories"
	"campaign/internal/defined"
	"campaign/internal/pkg"
	"campaign/proto/generate/servicepb"
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// SubscribeSubscription implements servicepb.CampaignServiceServer.
func (c *campaignHandler) SubscribeSubscription(ctx context.Context, req *servicepb.SubscribeSubscriptionRequest) (*servicepb.SubscribeSubscriptionReply, error) {
	auth, isLogined := pkg.GetAuthContext(ctx)
	if !isLogined {
		return nil, errors.New("unauthorized")
	}

	if req.SubscriptionId != defined.SILVER_SUBSCRIPTION || req.CampaignId == "" {
		return &servicepb.SubscribeSubscriptionReply{}, nil
	}

	campaignId, hasCampaignVoucher := auth.Payload[defined.CAMPAIGNID_KEY]
	if !hasCampaignVoucher {
		return &servicepb.SubscribeSubscriptionReply{}, nil
	}

	//check if campaign is ended
	_, err := c.cache.Get(ctx, defined.GetCampaignEndedCacheKey(campaignId)).Result()
	if err != redis.Nil {
		if err != nil {
			return nil, errors.New("internal error, try again later")
		}
		return nil, errors.New("campaign voucher is empty")
	}

	//decrement voucher remaining
	curRemaining, err := c.cache.Decr(ctx, defined.GetCampaignRemainingCacheKey(campaignId)).Result()
	if err != nil {
		return nil, errors.New("internal error, try again later")
	}
	if curRemaining < 0 {
		// set campaign ended
		c.cache.Set(ctx, defined.GetCampaignEndedCacheKey(campaignId), true, time.Hour*24)

		return nil, errors.New("campaign voucher is empty")
	}

	if curRemaining == 0 {
		c.cache.Set(ctx, defined.GetCampaignEndedCacheKey(campaignId), true, time.Hour*24)
	}
	orderId, _ := pkg.GenerateUniqueID(defined.ORDER_ID_LENGTH)

	simulateSendToQueue(orderPayload{
		Id:             orderId,
		Username:       auth.Username,
		SubscriptionId: req.SubscriptionId,
		CampaignId:     auth.Payload[defined.CAMPAIGNID_KEY],
	})

	return &servicepb.SubscribeSubscriptionReply{
		Status: defined.PROCESSING_STATUS,
	}, nil
}

type orderPayload struct {
	Username       string
	SubscriptionId string
	CampaignId     string
	Id             string
}

var onOrderUsernamesCh = make(chan orderPayload, 1000)

func simulateSendToQueue(order orderPayload) {
	onOrderUsernamesCh <- order
}

func (c *campaignHandler) simulateQueueWorker(ctx context.Context) {
	go func() {
		for order := range onOrderUsernamesCh {
			err := c.db.InsertUserSubscription(ctx, &repositories.UserSubscription{
				Id:             order.Id,
				Username:       order.Username,
				Status:         defined.ACTIVE_STATUS,
				SubscriptionId: order.SubscriptionId,
				CampaignId:     order.CampaignId,
			})
			if err != nil {
				simulateRetryToQueue(order)
			}
			//ack
		}
	}()
}

func simulateRetryToQueue(order orderPayload) {
	// onOrderUsernamesCh <- order
}
