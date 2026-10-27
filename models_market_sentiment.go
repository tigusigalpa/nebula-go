// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// CallsResponse models an object of the `market sentiment` service.
type CallsResponse struct {
	BaseResponse
	// Chronologically ordered observations.
	Series []CallsResponseSeriesPoint `json:"series"`
}

// CallsResponseSeriesPoint models an object of the `market sentiment` service.
type CallsResponseSeriesPoint struct {
	// Number of long calls in the bucket.
	LongCalls int64 `json:"long_calls"`
	// Number of short calls in the bucket.
	ShortCalls int64 `json:"short_calls"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp string `json:"timestamp"`
}

// CallReturnsResponse models an object of the `market sentiment` service.
type CallReturnsResponse struct {
	BaseResponse
	// Mean direction-aligned return across resolved calls, in percentage points.
	AverageReturnPct float64 `json:"average_return_pct"`
	// Return-distribution buckets in ascending order.
	Buckets []CallReturnsResponseBucket `json:"buckets"`
	// Resolved calls with a negative direction-aligned return.
	Losses int64 `json:"losses"`
	// Median direction-aligned return across resolved calls, in percentage points.
	MedianReturnPct float64 `json:"median_return_pct"`
	// Number of calls with a resolved outcome.
	ResolvedCalls int64 `json:"resolved_calls"`
	// How many days of resolved calls the statistics include.
	ResolvedDays int `json:"resolved_days"`
	// Winning resolved calls as a percentage from 0 to 100.
	WinRatePct float64 `json:"win_rate_pct"`
	// Resolved calls with a positive direction-aligned return.
	Wins int64 `json:"wins"`
}

// CallReturnsResponseBucket models an object of the `market sentiment` service.
type CallReturnsResponseBucket struct {
	// Resolved calls in this return bucket.
	CallCount int64 `json:"call_count"`
	// Sum of direction-aligned returns with each call capped at ±100%, in percentage points.
	CappedTotalAlignedReturnPct float64 `json:"capped_total_aligned_return_pct"`
	// Human-readable bucket or index label.
	Label string `json:"label"`
	// Inclusive lower return bound in percentage points; null for the first, open-ended bucket.
	LowerBoundPct *float64 `json:"lower_bound_pct"`
	// Sum of direction-aligned returns in this bucket, in percentage points.
	TotalAlignedReturnPct float64 `json:"total_aligned_return_pct"`
	// Exclusive upper return bound in percentage points; null for the last, open-ended bucket.
	UpperBoundPct *float64 `json:"upper_bound_pct"`
}

// ConvictionIndexResponse models an object of the `market sentiment` service.
type ConvictionIndexResponse struct {
	BaseResponse
	// Chronologically ordered observations.
	Series []ConvictionIndexResponseSeriesPoint `json:"series"`
}

// ConvictionIndexResponseSeriesPoint models an object of the `market sentiment` service.
type ConvictionIndexResponseSeriesPoint struct {
	// Nebula index value on a 0–100 scale.
	Index float64 `json:"index"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
}

// FearGreedIndexResponse models an object of the `market sentiment` service.
type FearGreedIndexResponse struct {
	BaseResponse
	// Most recent value for each returned series or the current index reading.
	Current *FearGreedIndexResponseCurrent `json:"current"`
	// Methodology and request-window metadata needed to interpret the response.
	Metadata *FearGreedIndexResponseMetadata `json:"metadata"`
	// Chronologically ordered observations.
	Series []FearGreedIndexResponseSeriesPoint `json:"series"`
}

// FearGreedIndexResponseCurrent models an object of the `market sentiment` service.
type FearGreedIndexResponseCurrent struct {
	// Nebula index value on a 0–100 scale.
	Index float64 `json:"index"`
	// Human-readable bucket or index label.
	Label string `json:"label"`
}

// FearGreedIndexResponseMetadata models an object of the `market sentiment` service.
type FearGreedIndexResponseMetadata struct {
	// Weighted components used to calculate the index.
	Components []FearGreedIndexResponseMetadataComponent `json:"components"`
	// Lookback window represented by the response, in hours.
	Hours int `json:"hours"`
	// Ordered labels used to interpret index ranges.
	Labels []string `json:"labels"`
	// Published scale or interpretation of the metric.
	Scale string `json:"scale"`
}

// FearGreedIndexResponseMetadataComponent models an object of the `market sentiment` service.
type FearGreedIndexResponseMetadataComponent struct {
	// Human-readable explanation of this component.
	Description string `json:"description"`
	// Human-readable name.
	Name string `json:"name"`
	// Component weight as a decimal fraction.
	Weight float64 `json:"weight"`
}

// FearGreedIndexResponseSeriesPoint models an object of the `market sentiment` service.
type FearGreedIndexResponseSeriesPoint struct {
	// Nebula index value on a 0–100 scale.
	Index float64 `json:"index"`
	// Human-readable bucket or index label.
	Label string `json:"label"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
}

// PriceLevelsResponse models an object of the `market sentiment` service.
type PriceLevelsResponse struct {
	BaseResponse
	// Socially mentioned price-level buckets.
	Levels []PriceLevelsResponseLevel `json:"levels"`
}

// PriceLevelsResponseLevel models an object of the `market sentiment` service.
type PriceLevelsResponseLevel struct {
	// Price mentions in this bucket.
	Mentions int `json:"mentions"`
	// Latest or observed market price in US dollars.
	PriceUSD float64 `json:"price_usd"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
}

// SocialVolumeIndexResponse models an object of the `market sentiment` service.
type SocialVolumeIndexResponse struct {
	BaseResponse
	// Chronologically ordered observations.
	Series []SocialVolumeIndexResponseSeriesPoint `json:"series"`
}

// SocialVolumeIndexResponseSeriesPoint models an object of the `market sentiment` service.
type SocialVolumeIndexResponseSeriesPoint struct {
	// The asset's posts on the day itself, bots left out.
	Posts float64 `json:"posts"`
	// The asset's posts over the trailing 7 days as a multiple of its median 7-day post count over the previous 90 days; 1 is usual.
	SocialRelative float64 `json:"social_relative"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
	// The asset's dollar trading volume over the trailing 7 days as a multiple of its median 7-day volume over the previous 90 days; 1 is usual.
	TradingRelative float64 `json:"trading_relative"`
	// From 0 to 100: 50 means trading and conversation are equally elevated against their normal; above 50 trading is ahead of the conversation, below 50 the conversation is ahead of tra
	VolumeIndex float64 `json:"volume_index"`
	// The asset's dollar trading volume on the day itself; 0 on a day without trading, such as a weekend for a stock.
	VolumeUSD float64 `json:"volume_usd"`
}

// SocialVolumeProfileResponse models an object of the `market sentiment` service.
type SocialVolumeProfileResponse struct {
	BaseResponse
	// Price bins ordered from low to high.
	Bins []SocialVolumeProfileResponseBin `json:"bins"`
	// Post count in the busiest price bin.
	MaximumPostsInBin float64 `json:"maximum_posts_in_bin"`
	// Highest profile price in US dollars.
	MaximumPriceUSD float64 `json:"maximum_price_usd"`
	// Methodology and request-window metadata needed to interpret the response.
	Metadata *SocialVolumeProfileResponseMetadata `json:"metadata"`
	// Lowest profile price in US dollars.
	MinimumPriceUSD float64 `json:"minimum_price_usd"`
	// Zero-based index of the busiest bin; absent when the profile is empty.
	PointOfControlBinIndex *int `json:"point_of_control_bin_index"`
	// Post count over the metric's active lookback window.
	TotalPosts float64 `json:"total_posts"`
}

// SocialVolumeProfileResponseBin models an object of the `market sentiment` service.
type SocialVolumeProfileResponseBin struct {
	// How evenly activity is split between directions, from 0 to 1.
	BalanceScore float64 `json:"balance_score"`
	// Posts classified as bearish.
	BearishPosts float64 `json:"bearish_posts"`
	// Bearish share of directional posts in this bin, from 0 to 1.
	BearishShare float64 `json:"bearish_share"`
	// Posts classified as bullish.
	BullishPosts float64 `json:"bullish_posts"`
	// Bullish share of directional posts in this bin, from 0 to 1.
	BullishShare float64 `json:"bullish_share"`
	// How slowly price moved while in the bin, from 0 to 1.
	CompressionScore float64 `json:"compression_score"`
	// Directional post imbalance (bullish − bearish) / (bullish + bearish), from −1 to 1; absent with no directional posts.
	DirectionalImbalance *float64 `json:"directional_imbalance"`
	// True for the bin containing the greatest post activity.
	IsPointOfControl bool `json:"is_point_of_control"`
	// Posts classified as neutral.
	NeutralPosts float64 `json:"neutral_posts"`
	// Exclusive upper edge of the price bin in US dollars.
	PriceHighUSD float64 `json:"price_high_usd"`
	// Inclusive lower edge of the price bin in US dollars.
	PriceLowUSD float64 `json:"price_low_usd"`
	// Midpoint of the price bin in US dollars.
	PriceMidUSD float64 `json:"price_mid_usd"`
	// Composite prominence score for the bin, from 0 to 1.
	SignificanceScore float64 `json:"significance_score"`
	// Post count over the metric's active lookback window.
	TotalPosts float64 `json:"total_posts"`
	// Post count relative to the busiest bin, from 0 to 1.
	TotalShare float64 `json:"total_share"`
	// Separate times the market price entered this bin during the window.
	VisitCount int `json:"visit_count"`
}

// SocialVolumeProfileResponseMetadata models an object of the `market sentiment` service.
type SocialVolumeProfileResponseMetadata struct {
	// UTC time at which the metric was calculated.
	AsOf *time.Time `json:"as_of"`
	// Number of requested price bins.
	BinCount int `json:"bin_count"`
	// Results keyed or grouped by canonical author-cluster identifier.
	Clusters []string `json:"clusters"`
	// Lookback window represented by the response, in hours.
	Hours int `json:"hours"`
	// Sentiment buckets successfully matched to a market price.
	PricedBuckets int `json:"priced_buckets"`
	// Sentiment buckets omitted because no nearby market price was available.
	UnpricedBuckets int `json:"unpriced_buckets"`
}
