// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// PriceHistoryResponse models an object of the `market data` service.
type PriceHistoryResponse struct {
	BaseResponse
	// Chronologically ordered observations.
	Series []PriceHistoryResponseSeriesPoint `json:"series"`
}

// PriceHistoryResponseSeriesPoint models an object of the `market data` service.
type PriceHistoryResponseSeriesPoint struct {
	// Observed high price in US dollars.
	HighUSD *float64 `json:"high_usd"`
	// Observed low price in US dollars.
	LowUSD *float64 `json:"low_usd"`
	// Latest or observed market price in US dollars.
	PriceUSD *float64 `json:"price_usd"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
}
