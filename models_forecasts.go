// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// ForecastAssetsResponse models an object of the `forecasts` service.
type ForecastAssetsResponse struct {
	BaseResponse
	// Assets with a current forecast: stocks, commodities, and indices first, then crypto; most-discussed first within each class.
	Assets []ForecastAssetsResponseAsset `json:"assets"`
	// True when another page is available after the returned records.
	HasMore bool `json:"has_more"`
	// Maximum records requested.
	Limit int `json:"limit"`
	// Number of matching records skipped.
	Offset int `json:"offset"`
	// How many assets match the request across every page.
	Total int `json:"total"`
}

// ForecastAssetsResponseAsset models an object of the `forecasts` service.
type ForecastAssetsResponseAsset struct {
	// UTC start of the asset's current forecast: the end of the last complete 15 minutes it was computed from.
	AsOf time.Time `json:"as_of"`
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID string `json:"asset_id"`
	// Human-readable name.
	Name string `json:"name"`
	// Display ticker or market symbol.
	Symbol string `json:"symbol"`
}

// MindshareForecastResponse models an object of the `forecasts` service.
type MindshareForecastResponse struct {
	BaseResponse
	// UTC start of the forecast: the end of the last complete 15 minutes it was computed from. Null when the asset has never been scored.
	AsOf *time.Time `json:"as_of"`
	// The chance the share is higher than current_share in 72 hours, from 0 to 1.
	ChanceHigher72h *float64 `json:"chance_higher_72h"`
	// The asset's share at as_of, from 0 to 1: where every path starts, and the value the chances are measured against.
	CurrentShare *float64 `json:"current_share"`
	// How far ahead the forecast runs, in hours: always 72.
	HorizonHours int `json:"horizon_hours"`
	// The 90th percentile of the share in 72 hours: a 90% chance of ending below it.
	LikelyHighShare *float64 `json:"likely_high_share"`
	// The 10th percentile of the share in 72 hours: a 90% chance of ending above it.
	LikelyLowShare *float64 `json:"likely_low_share"`
	// Why the asset cannot be forecast: too_little_history, stale_data, or missing_price. Null when status is ok.
	Reason *ForecastReason `json:"reason"`
	// True when this request re-ran the forecast (refresh=true and the re-run finished in time); false when it is the latest scheduled five-minute forecast.
	Refreshed bool `json:"refreshed"`
	// 24 possible paths for the share, each aligned with series. Draw them to show the spread of outcomes.
	SamplePaths [][]float64 `json:"sample_paths"`
	// The forecast distribution every 4 hours from as_of: 18 steps, oldest first.
	Series []MindshareForecastResponseSeriesPoint `json:"series"`
	// ok, or not_forecastable when the asset cannot be forecast (the summary fields are then null and the series empty; reason says why).
	Status ForecastStatus `json:"status"`
	// The chance of a sudden surge in attention on the asset within 72 hours, from 0 to 1. Null when there is no estimate.
	TakeoffChance72h *float64 `json:"takeoff_chance_72h"`
}

// MindshareForecastResponseSeriesPoint models an object of the `forecasts` service.
type MindshareForecastResponseSeriesPoint struct {
	// The expected share at this step.
	MeanShare float64 `json:"mean_share"`
	// The 10th percentile of the share at this step, from 0 to 1.
	P10Share float64 `json:"p10_share"`
	// The 25th percentile of the share at this step, from 0 to 1.
	P25Share float64 `json:"p25_share"`
	// The 50th percentile of the share at this step, from 0 to 1.
	P50Share float64 `json:"p50_share"`
	// The 75th percentile of the share at this step, from 0 to 1.
	P75Share float64 `json:"p75_share"`
	// The 90th percentile of the share at this step, from 0 to 1.
	P90Share float64 `json:"p90_share"`
	// 21 quantiles of the share at this step: the 0th, 5th, 10th, ... 100th percentiles, lowest first. Interpolate between them for the chance of landing in any range.
	QuantileShares []float64 `json:"quantile_shares"`
	// UTC time the step forecasts: as_of plus 4, 8, ... 72 hours.
	Timestamp time.Time `json:"timestamp"`
}

// SentimentForecastResponse models an object of the `forecasts` service.
type SentimentForecastResponse struct {
	BaseResponse
	// UTC start of the forecast: the end of the last complete 15 minutes it was computed from. Null when the asset has never been scored.
	AsOf *time.Time `json:"as_of"`
	// The chance the sentiment signal is higher than current_sentiment_signal in 24 hours, from 0 to 1.
	ChanceHigher24h *float64 `json:"chance_higher_24h"`
	// True when chance_higher_24h is under 0.3 or over 0.7: the forecast leans clearly one way. Nearer 0.5, read it as no clear direction.
	ChanceHigher24hConfident *bool `json:"chance_higher_24h_confident"`
	// The chance the sentiment signal is higher than current_sentiment_signal in 72 hours, from 0 to 1.
	ChanceHigher72h *float64 `json:"chance_higher_72h"`
	// True when chance_higher_72h is under 0.3 or over 0.7: the forecast leans clearly one way. Nearer 0.5, read it as no clear direction.
	ChanceHigher72hConfident *bool `json:"chance_higher_72h_confident"`
	// The sentiment signal at as_of, on the /sentiment scale: bullish minus bearish posts over the trailing 24 hours, recent posts weighted more. Every path starts here, and the chances
	CurrentSentimentSignal *float64 `json:"current_sentiment_signal"`
	// How far ahead the forecast runs, in hours: always 72.
	HorizonHours int `json:"horizon_hours"`
	// Why the asset cannot be forecast: too_little_history, stale_data, or missing_price. Null when status is ok.
	Reason *ForecastReason `json:"reason"`
	// True when this request re-ran the forecast (refresh=true and the re-run finished in time); false when it is the latest scheduled forecast.
	Refreshed bool `json:"refreshed"`
	// 24 possible paths for the sentiment signal, each aligned with series. Draw them to show the spread of outcomes.
	SamplePaths [][]float64 `json:"sample_paths"`
	// The forecast every 4 hours from as_of: 18 steps, oldest first.
	Series []SentimentForecastResponseSeriesPoint `json:"series"`
	// ok, or not_forecastable when the asset cannot be forecast (the summary fields are then null and the series empty; reason says why).
	Status ForecastStatus `json:"status"`
	// 1 while sentiment is turning bullish (the last few hours lean clearly more bullish than the last day), -1 while it is turning bearish, 0 otherwise.
	Turning *int `json:"turning"`
}

// SentimentForecastResponseSeriesPoint models an object of the `forecasts` service.
type SentimentForecastResponseSeriesPoint struct {
	// The bearish posts expected in this 4-hour step alone, bots excluded: like /sentiment's interval posts, not its trailing 24-hour counts.
	BearishPosts *SentimentForecastResponseSeriesPointBearishPosts `json:"bearish_posts"`
	// The bullish posts expected in this 4-hour step alone, bots excluded: like /sentiment's interval posts, not its trailing 24-hour counts.
	BullishPosts *SentimentForecastResponseSeriesPointBearishPosts `json:"bullish_posts"`
	// The neutral posts expected in this 4-hour step alone, bots excluded: like /sentiment's interval posts, not its trailing 24-hour counts.
	NeutralPosts *SentimentForecastResponseSeriesPointBearishPosts `json:"neutral_posts"`
	// The sentiment signal at the end of the step, on the /sentiment scale: where the line is likely to be.
	SentimentSignal *SentimentForecastResponseSeriesPointSentimentSignal `json:"sentiment_signal"`
	// UTC time the step ends: as_of plus 4, 8, ... 72 hours.
	Timestamp time.Time `json:"timestamp"`
	// The posts (bullish, bearish, and neutral together) expected in this 4-hour step alone, bots excluded: like /sentiment's interval posts, not its trailing 24-hour counts. Its percent
	TotalPosts *SentimentForecastResponseSeriesPointBearishPosts `json:"total_posts"`
}

// SentimentForecastResponseSeriesPointBearishPosts models an object of the `forecasts` service.
type SentimentForecastResponseSeriesPointBearishPosts struct {
	// The expected number of bearish posts in the step.
	Mean float64 `json:"mean"`
	// The 10th percentile of the bearish posts in the step.
	P10 float64 `json:"p10"`
	// The 25th percentile of the bearish posts in the step.
	P25 float64 `json:"p25"`
	// The 50th percentile of the bearish posts in the step.
	P50 float64 `json:"p50"`
	// The 75th percentile of the bearish posts in the step.
	P75 float64 `json:"p75"`
	// The 90th percentile of the bearish posts in the step.
	P90 float64 `json:"p90"`
}

// SentimentForecastResponseSeriesPointSentimentSignal models an object of the `forecasts` service.
type SentimentForecastResponseSeriesPointSentimentSignal struct {
	// The expected sentiment signal at the end of the step.
	Mean float64 `json:"mean"`
	// The 10th percentile of the sentiment signal at the end of the step.
	P10 float64 `json:"p10"`
	// The 50th percentile of the sentiment signal at the end of the step.
	P50 float64 `json:"p50"`
	// The 90th percentile of the sentiment signal at the end of the step.
	P90 float64 `json:"p90"`
	// 21 quantiles of the sentiment signal at the end of the step: the 0th, 5th, 10th, ... 100th percentiles, lowest first. Interpolate between them for the chance of landing in any rang
	Quantiles []float64 `json:"quantiles"`
}
