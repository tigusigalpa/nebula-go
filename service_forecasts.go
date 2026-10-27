// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// ListAssets - Forecastable assets. GET /forecasts/assets
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-forecasts-assets/
func (s *ForecastsService) ListAssets(ctx context.Context, req *ForecastAssetsRequest) (*ForecastAssetsResponse, error) {
	return get[ForecastAssetsResponse](ctx, s.client, "/forecasts/assets", req, "list_forecast_assets")
}

// Mindshare - Mindshare forecast. GET /forecasts/mindshare
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-forecasts-mindshare/
func (s *ForecastsService) Mindshare(ctx context.Context, req *MindshareForecastRequest) (*MindshareForecastResponse, error) {
	return get[MindshareForecastResponse](ctx, s.client, "/forecasts/mindshare", req, "get_mindshare_forecast")
}

// Sentiment - Sentiment forecast. GET /forecasts/sentiment
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-forecasts-sentiment/
func (s *ForecastsService) Sentiment(ctx context.Context, req *SentimentForecastRequest) (*SentimentForecastResponse, error) {
	return get[SentimentForecastResponse](ctx, s.client, "/forecasts/sentiment", req, "get_sentiment_forecast")
}
