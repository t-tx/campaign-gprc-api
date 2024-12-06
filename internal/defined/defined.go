package defined

const (
	CAMPAIGNID_KEY = "campaign_id"

	SILVER_SUBSCRIPTION = "silver"
	ACTIVE_STATUS       = "active"
	PROCESSING_STATUS   = "processing"
	ORDER_ID_LENGTH     = 20
)

func GetCampaignRemainingCacheKey(campaignId string) string {
	return "campaign_remaining:" + campaignId
}

func GetCampaignEndedCacheKey(campaignId string) string {
	return "campaign_ended:" + campaignId
}

func GetCampaignUsedCacheKey(campaignId string) string {
	return "campaign_used:" + campaignId
}
