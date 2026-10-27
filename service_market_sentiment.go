// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// Calls - Long/short calls. GET /calls
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-calls/
func (s *MarketSentimentService) Calls(ctx context.Context, req *CallsRequest) (*CallsResponse, error) {
	return get[CallsResponse](ctx, s.client, "/calls", req, "get_calls")
}

// CallReturns - Calls outcome distribution. GET /calls/returns
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-calls-returns/
func (s *MarketSentimentService) CallReturns(ctx context.Context, req *CallReturnsRequest) (*CallReturnsResponse, error) {
	return get[CallReturnsResponse](ctx, s.client, "/calls/returns", req, "get_call_returns")
}

// ConvictionIndex - Conviction index. GET /conviction-index
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-conviction-index/
func (s *MarketSentimentService) ConvictionIndex(ctx context.Context, req *ConvictionIndexRequest) (*ConvictionIndexResponse, error) {
	return get[ConvictionIndexResponse](ctx, s.client, "/conviction-index", req, "get_conviction_index")
}

// FearGreedIndex - Fear and Greed index. GET /fear-greed-index
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-fear-greed-index/
func (s *MarketSentimentService) FearGreedIndex(ctx context.Context, req *FearGreedIndexRequest) (*FearGreedIndexResponse, error) {
	return get[FearGreedIndexResponse](ctx, s.client, "/fear-greed-index", req, "get_fear_greed_index")
}

// PriceLevels - Socially mentioned price levels. GET /price-levels
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-price-levels/
func (s *MarketSentimentService) PriceLevels(ctx context.Context, req *PriceLevelsRequest) (*PriceLevelsResponse, error) {
	return get[PriceLevelsResponse](ctx, s.client, "/price-levels", req, "get_price_levels")
}

// SocialVolumeIndex - Social vs trading volume. GET /social-volume-index
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-social-volume-index/
func (s *MarketSentimentService) SocialVolumeIndex(ctx context.Context, req *SocialVolumeIndexRequest) (*SocialVolumeIndexResponse, error) {
	return get[SocialVolumeIndexResponse](ctx, s.client, "/social-volume-index", req, "get_social_volume_index")
}

// SocialVolumeProfile - Social volume profile. GET /social-volume-profile
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-social-volume-profile/
func (s *MarketSentimentService) SocialVolumeProfile(ctx context.Context, req *SocialVolumeProfileRequest) (*SocialVolumeProfileResponse, error) {
	return get[SocialVolumeProfileResponse](ctx, s.client, "/social-volume-profile", req, "get_social_volume_profile")
}
