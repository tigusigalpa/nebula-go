// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// PredictionsResponse models an object of the `predictions` service.
type PredictionsResponse struct {
	BaseResponse
	// True when another page is available after the returned records.
	HasMore bool `json:"has_more"`
	// Maximum records requested.
	Limit int `json:"limit"`
	// Canonical markets, ordered by best activity rank within each venue. Venue volumes are not summed.
	Markets []PredictionMarketResponse `json:"markets"`
	// Number of matching records skipped.
	Offset int `json:"offset"`
}

// PredictionMarketResponse models an object of the `predictions` service.
type PredictionMarketResponse struct {
	BaseResponse
	// Linked identities for existing analytics endpoints. Assets carry asset_id and asset_class; non-assets carry subject_id. An identity without selectors is not yet available for analy
	Entities []PredictionMarketResponseEntity `json:"entities"`
	// Stable identity within this object type.
	ID string `json:"id"`
	// Venue listings with separately retained rules and quotes.
	Listings []PredictionMarketResponseListing `json:"listings"`
	// Original contract question.
	Question string `json:"question"`
	// Legacy attachment metadata retained for compatibility. Use entities for analytics selectors.
	Subjects []PredictionMarketResponseSubject `json:"subjects"`
}

// PredictionMarketResponseEntity models an object of the `predictions` service.
type PredictionMarketResponseEntity struct {
	// Asset class: equity, commodity, index or crypto.
	AssetClass *AssetClass `json:"asset_class"`
	// Asset identifier accepted by asset analytics endpoints.
	AssetID *string `json:"asset_id"`
	// Canonical subject display name.
	Name string `json:"name"`
	// Non-asset subject identifier accepted by subject analytics endpoints. Mutually exclusive with asset selectors.
	SubjectID *string `json:"subject_id"`
}

// PredictionMarketResponseListing models an object of the `predictions` service.
type PredictionMarketResponseListing struct {
	// Venue trading close in UTC; not necessarily event occurrence or settlement.
	ClosesAt *time.Time `json:"closes_at"`
	// Venue event identifier, or empty if absent.
	EventID string `json:"event_id"`
	// Stable identity within this object type.
	ID string `json:"id"`
	// Venue-reported liquidity in US dollars; null when unavailable.
	LiquidityUSD *float64 `json:"liquidity_usd"`
	// Canonical market identity; equivalent venue listings share it.
	MarketID string `json:"market_id"`
	// UTC time Nebula fetched this quote. Inspect freshness before comparing.
	ObservedAt time.Time `json:"observed_at"`
	// Venue open contracts; null when unavailable.
	OpenInterest *float64 `json:"open_interest"`
	// Venue opening time in UTC; null when unavailable.
	OpensAt *time.Time `json:"opens_at"`
	// Ordered native outcomes with independent price evidence.
	Outcomes []PredictionMarketResponseListingOutcome `json:"outcomes"`
	// Venue: polymarket or kalshi.
	Provider PredictionProvider `json:"provider"`
	// Native venue market identifier.
	ProviderID string `json:"provider_id"`
	// Venue metadata update time; null when unavailable.
	ProviderUpdatedAt *time.Time `json:"provider_updated_at"`
	// Original contract question.
	Question string `json:"question"`
	// Venue resolution source, or empty if absent.
	ResolutionSource string `json:"resolution_source"`
	// Venue expected resolution time in UTC; null when unavailable.
	ResolvesAt *time.Time `json:"resolves_at"`
	// Venue settlement result; empty before a known result.
	Result string `json:"result"`
	// Legacy contract revision token; retained for compatibility.
	Revision string `json:"revision"`
	// Full venue resolution rules.
	Rules string `json:"rules"`
	// Original venue listing URL.
	SourceURL string `json:"source_url"`
	// Listing lifecycle: open, closed, resolved or disputed. For subjects: resolved or pending.
	Status string `json:"status"`
	// Cumulative venue volume in volume_unit; null when unavailable.
	Volume *float64 `json:"volume"`
	// Last 24-hour venue volume in volume_unit; null when unavailable.
	Volume24h *float64 `json:"volume_24h"`
	// usd for Polymarket, contracts for Kalshi. These quantities cannot be added.
	VolumeUnit string `json:"volume_unit"`
}

// PredictionMarketResponseListingOutcome models an object of the `predictions` service.
type PredictionMarketResponseListingOutcome struct {
	// Venue ask as a fraction from 0 to 1; null when unavailable.
	Ask *float64 `json:"ask"`
	// Venue bid as a fraction from 0 to 1; null when unavailable.
	Bid *float64 `json:"bid"`
	// Stable identity within this object type.
	ID string `json:"id"`
	// Venue outcome label.
	Label string `json:"label"`
	// Last traded price from 0 to 1; null when unavailable. Kalshi No is the complement of its reported Yes last trade; Polymarket is token-specific. Observation time is not trade time.
	LastTrade *float64 `json:"last_trade"`
	// Bid-ask midpoint from 0 to 1; null when unavailable.
	Midpoint *float64 `json:"midpoint"`
	// Zero-based native outcome position.
	Position int `json:"position"`
	// Price as a fraction from 0 to 1. Current quotes use valid bid-ask midpoints; history may include provider prices as identified by price_basis. Not subject sentiment.
	Price *float64 `json:"price"`
	// midpoint or unavailable for current quotes; historical series may also use provider for venue history of unspecified quote basis.
	PriceBasis string `json:"price_basis"`
	// Provider-reported reference or settlement price from 0 to 1; null when unavailable. Retained independently of midpoint.
	ProviderPrice *float64 `json:"provider_price"`
}

// PredictionMarketResponseSubject models an object of the `predictions` service.
type PredictionMarketResponseSubject struct {
	// Ordered backend-resolved fallback icon candidates.
	FallbackImageUrls []string `json:"fallback_image_urls"`
	// First backend-resolved subject icon candidate; empty uses fallbacks.
	ImageURL string `json:"image_url"`
	// Canonical subject display name.
	Name string `json:"name"`
	// Question attachment role: entity, topic or place.
	Role string `json:"role"`
	// Listing lifecycle: open, closed, resolved or disputed. For subjects: resolved or pending.
	Status string `json:"status"`
	// Canonical registry subject identifier; preserves existing asset or topic eligibility.
	SubjectID string `json:"subject_id"`
	// Canonical subject types from the existing identity resolver.
	Types []string `json:"types"`
}

// PredictionHistoryResponse models an object of the `predictions` service.
type PredictionHistoryResponse struct {
	BaseResponse
	// Venue-prefixed native market identity.
	ListingID string `json:"listing_id"`
	// Native outcome identity within the listing.
	OutcomeID string `json:"outcome_id"`
	// Deprecated compatibility mirror of series, using time instead of timestamp. Use series for new integrations.
	Points []PredictionHistoryResponsePoint `json:"points"`
	// Historical observations in chronological order, with an explicit price basis; no invented samples.
	Series []PredictionHistoryResponseSeriesPoint `json:"series"`
}

// PredictionHistoryResponsePoint models an object of the `predictions` service.
type PredictionHistoryResponsePoint struct {
	// Price as a fraction from 0 to 1. Current quotes use valid bid-ask midpoints; history may include provider prices as identified by price_basis. Not subject sentiment.
	Price float64 `json:"price"`
	// midpoint or unavailable for current quotes; historical series may also use provider for venue history of unspecified quote basis.
	PriceBasis string `json:"price_basis"`
	// UTC observation time.
	Time time.Time `json:"time"`
}

// PredictionHistoryResponseSeriesPoint models an object of the `predictions` service.
type PredictionHistoryResponseSeriesPoint struct {
	// Price as a fraction from 0 to 1. Current quotes use valid bid-ask midpoints; history may include provider prices as identified by price_basis. Not subject sentiment.
	Price float64 `json:"price"`
	// midpoint or unavailable for current quotes; historical series may also use provider for venue history of unspecified quote basis.
	PriceBasis string `json:"price_basis"`
	// Observation time in UTC RFC 3339.
	Timestamp time.Time `json:"timestamp"`
}
