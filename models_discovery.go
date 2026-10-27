// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// ScreenerResponse models an object of the `discovery` service.
type ScreenerResponse struct {
	BaseResponse
	// The market whose assets were ranked; absent for subjects.
	AssetClass *AssetClass `json:"asset_class"`
	// The author clusters whose posts were counted, canonicalised; empty when every author except bots was counted.
	Clusters []string `json:"clusters"`
	// UTC start of the ranked 24 hours. The comparison window is the 24 hours before it.
	From time.Time `json:"from"`
	// What was ranked: asset or subject, as requested.
	Kind Kind `json:"kind"`
	// All posts in the ranked population over the previous 24 hours.
	PreviousTotalPosts int64 `json:"previous_total_posts"`
	// Ranked results, best first.
	Results []ScreenerResponseResult `json:"results"`
	// UTC end of the ranked 24 hours: when the ranking was computed.
	To time.Time `json:"to"`
	// All posts in the ranked population over the 24 hours: the denominator of share.
	TotalPosts int64 `json:"total_posts"`
}

// ScreenerResponseResult models an object of the `discovery` service.
type ScreenerResponseResult struct {
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass *AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID *string `json:"asset_id"`
	// Posts classified as bearish.
	BearishPosts int64 `json:"bearish_posts"`
	// Posts classified as bullish.
	BullishPosts int64 `json:"bullish_posts"`
	// What the result is: asset (query it with asset_id and asset_class) or subject (query it with subject_id).
	Kind Kind `json:"kind"`
	// Human-readable name.
	Name string `json:"name"`
	// Posts classified as neutral.
	NeutralPosts int64 `json:"neutral_posts"`
	// This result's posts over the ranked 24 hours.
	Posts int64 `json:"posts"`
	// This result's posts in the immediately preceding window of equal length.
	PreviousPosts int64 `json:"previous_posts"`
	// Raw market-price change over 24 hours, in percentage points.
	PriceChange24hPct *float64 `json:"price_change_24h_pct"`
	// This result's posts as a fraction of all posts in the ranked population over the same 24 hours, from 0 to 1.
	Share float64 `json:"share"`
	// Change in share against the previous window, as a decimal fraction: 0.01 is one percentage point. Absent only when the population had no posts in the previous window.
	ShareChange *float64 `json:"share_change"`
	// The change in share relative to the previous window's share, in percent. Absent when the result had no posts in the previous window.
	ShareChangePct *float64 `json:"share_change_pct"`
	// Canonical identifier for a non-asset person, account, or topic.
	SubjectID *string `json:"subject_id"`
	// Display ticker or market symbol.
	Symbol *string `json:"symbol"`
}
