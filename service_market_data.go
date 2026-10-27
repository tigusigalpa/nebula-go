// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// PriceHistory - Price history. GET /price-history
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-price-history/
func (s *MarketDataService) PriceHistory(ctx context.Context, req *PriceHistoryRequest) (*PriceHistoryResponse, error) {
	return get[PriceHistoryResponse](ctx, s.client, "/price-history", req, "get_price_history")
}
