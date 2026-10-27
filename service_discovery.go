// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// Screener - Screener. GET /screener
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-screener/
func (s *DiscoveryService) Screener(ctx context.Context, req *ScreenerRequest) (*ScreenerResponse, error) {
	return get[ScreenerResponse](ctx, s.client, "/screener", req, "screener")
}

// TopMomentum - Top momentum. GET /top-momentum
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-top-momentum/
func (s *DiscoveryService) TopMomentum(ctx context.Context, req *TopMomentumRequest) (*ScreenerResponse, error) {
	return get[ScreenerResponse](ctx, s.client, "/top-momentum", req, "get_top_momentum")
}

// TopTrending - Top trending. GET /top-trending
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-top-trending/
func (s *DiscoveryService) TopTrending(ctx context.Context, req *TopTrendingRequest) (*ScreenerResponse, error) {
	return get[ScreenerResponse](ctx, s.client, "/top-trending", req, "get_top_trending")
}
