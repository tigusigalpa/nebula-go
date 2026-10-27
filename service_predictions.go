// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// ListMarkets - Prediction markets. GET /predictions
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-predictions/
func (s *PredictionsService) ListMarkets(ctx context.Context, req *PredictionsRequest) (*PredictionsResponse, error) {
	return get[PredictionsResponse](ctx, s.client, "/predictions", req, "list_prediction_markets")
}

// History - Prediction outcome history. GET /predictions/history
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-prediction-history/
func (s *PredictionsService) History(ctx context.Context, req *PredictionHistoryRequest) (*PredictionHistoryResponse, error) {
	return get[PredictionHistoryResponse](ctx, s.client, "/predictions/history", req, "get_prediction_history")
}

// Market - Prediction market detail. GET /predictions/market
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-prediction-market/
func (s *PredictionsService) Market(ctx context.Context, req *PredictionMarketRequest) (*PredictionMarketResponse, error) {
	return get[PredictionMarketResponse](ctx, s.client, "/predictions/market", req, "get_prediction_market")
}
