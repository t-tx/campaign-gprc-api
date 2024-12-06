package handler

import (
	"campaign/internal/defined"
	"campaign/internal/pkg"
	"campaign/proto/generate/servicepb"
	"context"
	"strconv"

	"github.com/shopspring/decimal"
)

func (c *campaignHandler) ListAvailableSubscription(ctx context.Context, req *servicepb.ListAvailableSubscriptionRequest) (*servicepb.ListAvailableSubscriptionReply, error) {
	auth, isLogined := pkg.GetAuthContext(ctx)
	if !isLogined {
		return c.listAvailableSubscription_guest(ctx, req)
	}
	return c.listAvailableSubscription_member(ctx, auth, req)

}

func (c *campaignHandler) listAvailableSubscription_member(ctx context.Context, authPayload *pkg.AuthContext, req *servicepb.ListAvailableSubscriptionRequest) (*servicepb.ListAvailableSubscriptionReply, error) {
	campaignId := authPayload.Payload[defined.CAMPAIGNID_KEY]
	v, err := c.cache.Get(ctx, defined.GetCampaignRemainingCacheKey(campaignId)).Result()
	var remaining int
	if err == nil {
		remaining, _ = strconv.Atoi(v)
	}

	subscriptions, err := c.db.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	var availableSubscription []*servicepb.ListAvailableSubscriptionReplyData
	for _, s := range subscriptions {
		var discount = zero
		if s.ID == defined.SILVER_SUBSCRIPTION && remaining > 0 {
			discount = discountPrice(s.Price, 30)
		}
		availableSubscription = append(availableSubscription, &servicepb.ListAvailableSubscriptionReplyData{
			Id:       s.ID,
			Name:     s.Name,
			Price:    s.Price.String(),
			Discount: discount.String(),
		})
	}
	return &servicepb.ListAvailableSubscriptionReply{
		Data: availableSubscription,
	}, nil
}

var (
	zero = decimal.NewFromInt(0)
)

func (c *campaignHandler) listAvailableSubscription_guest(ctx context.Context, req *servicepb.ListAvailableSubscriptionRequest) (*servicepb.ListAvailableSubscriptionReply, error) {
	subscriptions, err := c.db.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	var availableSubscription []*servicepb.ListAvailableSubscriptionReplyData
	for _, s := range subscriptions {
		availableSubscription = append(availableSubscription, &servicepb.ListAvailableSubscriptionReplyData{
			Id:    s.ID,
			Name:  s.Name,
			Price: s.Price.String(),
		})
	}
	return &servicepb.ListAvailableSubscriptionReply{
		Data: availableSubscription,
	}, nil
}

func discountPrice(price decimal.Decimal, discount int) decimal.Decimal {
	return price.Mul(decimal.NewFromInt(int64(discount))).Div(decimal.NewFromInt(100))
}
