// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

// ListAssetsResponse models an object of the `lookup` service.
type ListAssetsResponse struct {
	BaseResponse
	// Assets matching the supplied discovery filters.
	Assets []ListAssetsResponseAsset `json:"assets"`
	// True when another page is available after the returned records.
	HasMore bool `json:"has_more"`
	// Maximum records requested.
	Limit int `json:"limit"`
	// Number of matching records skipped.
	Offset int `json:"offset"`
}

// ListAssetsResponseAsset models an object of the `lookup` service.
type ListAssetsResponseAsset struct {
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID string `json:"asset_id"`
	// Primary social handle when available.
	Handle *string `json:"handle"`
	// Latest market capitalisation in US dollars. Reported for crypto; absent for stocks, commodities, and indices.
	MarketCapUSD *float64 `json:"market_cap_usd"`
	// Human-readable name.
	Name string `json:"name"`
	// Raw market-price change over 24 hours, in percentage points.
	PriceChange24hPct *float64 `json:"price_change_24h_pct"`
	// Latest or observed market price in US dollars.
	PriceUSD *float64 `json:"price_usd"`
	// Display ticker or market symbol.
	Symbol string `json:"symbol"`
}

// SearchResponse models an object of the `lookup` service.
type SearchResponse struct {
	BaseResponse
	// The free-text string that was resolved, echoed back.
	Query string `json:"query"`
	// Ranked results, best first.
	Results []SearchResponseResult `json:"results"`
}

// SearchResponseResult models an object of the `lookup` service.
type SearchResponseResult struct {
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass *AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID *string `json:"asset_id"`
	// What the result is: asset (query it with asset_id and asset_class) or subject (query it with subject_id).
	Kind Kind `json:"kind"`
	// The name or alias that matched the query, so you can show why a search landed where it did.
	MatchedAlias *string `json:"matched_alias"`
	// Human-readable name.
	Name string `json:"name"`
	// Canonical identifier for a non-asset person, account, or topic.
	SubjectID *string `json:"subject_id"`
	// Display ticker or market symbol.
	Symbol *string `json:"symbol"`
}

// ListSubjectsResponse models an object of the `lookup` service.
type ListSubjectsResponse struct {
	BaseResponse
	// True when another page is available after the returned records.
	HasMore bool `json:"has_more"`
	// Maximum records requested.
	Limit int `json:"limit"`
	// Number of matching records skipped.
	Offset int `json:"offset"`
	// Non-asset subjects matching the discovery filters.
	Subjects []ListSubjectsResponseSubject `json:"subjects"`
}

// ListSubjectsResponseSubject models an object of the `lookup` service.
type ListSubjectsResponseSubject struct {
	// Human-readable name.
	Name string `json:"name"`
	// Canonical identifier for a non-asset person, account, or topic.
	SubjectID string `json:"subject_id"`
}
