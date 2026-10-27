// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// EmotionsResponse models an object of the `social` service.
type EmotionsResponse struct {
	BaseResponse
	// Methodology and request-window metadata needed to interpret the response.
	Metadata *EmotionsResponseMetadata `json:"metadata"`
	// Emotion observations keyed by the detected emotion label.
	Series map[string][]EmotionsResponseSeriesEntry `json:"series"`
}

// EmotionsResponseMetadata models an object of the `social` service.
type EmotionsResponseMetadata struct {
	// How clients should interpret buckets omitted or emitted without data.
	EmptyBucketPolicy string `json:"empty_bucket_policy"`
	// Number of buckets with no qualifying data.
	EmptyBuckets int `json:"empty_buckets"`
	// Time bucket size used by the returned series.
	Granularity Granularity `json:"granularity"`
	// Lookback window represented by the response, in hours.
	Hours int `json:"hours"`
	// Most recent UTC bucket with no qualifying data.
	LatestEmptyBucket *time.Time `json:"latest_empty_bucket"`
}

// EmotionsResponseSeries models an object of the `social` service.
type EmotionsResponseSeries struct {
}

// EmotionsResponseSeriesEntry models an object of the `social` service.
type EmotionsResponseSeriesEntry struct {
	// Raw number of posts carrying this emotion in the bucket.
	PostCount int64 `json:"post_count"`
	// This emotion's share of all emotion tags in the bucket, from 0 to 1.
	Proportion float64 `json:"proportion"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
}

// MindshareResponse models an object of the `social` service.
type MindshareResponse struct {
	BaseResponse
	// Chronologically ordered observations.
	Series []MindshareResponseSeriesPoint `json:"series"`
}

// MindshareResponseSeriesPoint models an object of the `social` service.
type MindshareResponseSeriesPoint struct {
	// Share of eligible social attention as a decimal fraction from 0 to 1.
	Share float64 `json:"share"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
}

// RelatedSubjectsResponse models an object of the `social` service.
type RelatedSubjectsResponse struct {
	BaseResponse
	// Distinct authors represented by the matching posts.
	Authors int64 `json:"authors"`
	// Mean sentiment score across all anchor posts in the requested window.
	AverageSentimentScore float64 `json:"average_sentiment_score"`
	// Inclusive UTC start of the measured window.
	From time.Time `json:"from"`
	// Lookback window represented by the response, in hours.
	Hours int `json:"hours"`
	// Matching post count over the requested window.
	Posts int64 `json:"posts"`
	// Non-asset subjects returned by the discovery operation.
	Subjects []RelatedSubjectsResponseSubject `json:"subjects"`
	// Exclusive UTC end of the measured window.
	To time.Time `json:"to"`
}

// RelatedSubjectsResponseSubject models an object of the `social` service.
type RelatedSubjectsResponseSubject struct {
	// Distinct authors represented by the matching posts.
	Authors int64 `json:"authors"`
	// Posts classified as bearish.
	BearishPosts int64 `json:"bearish_posts"`
	// Posts classified as bullish.
	BullishPosts int64 `json:"bullish_posts"`
	// UTC time of the first matching post in the requested window.
	FirstPostAt *time.Time `json:"first_post_at"`
	// UTC time of the latest matching post in the requested window.
	LastPostAt *time.Time `json:"last_post_at"`
	// Human-readable name.
	Name string `json:"name"`
	// Posts classified as neutral.
	NeutralPosts int64 `json:"neutral_posts"`
	// Matching post count over the requested window.
	Posts int64 `json:"posts"`
	// Fractional change from posts_prior to posts_recent; positive values indicate growth.
	PostsDelta float64 `json:"posts_delta"`
	// Matching posts in the first half of the requested window.
	PostsPrior int64 `json:"posts_prior"`
	// Matching posts in the second half of the requested window.
	PostsRecent int64 `json:"posts_recent"`
	// Difference between sentiment_score for the shared conversation and the anchor's complete-window average.
	SentimentDelta float64 `json:"sentiment_delta"`
	// Signed sentiment score; positive is bullish and negative is bearish.
	SentimentScore float64 `json:"sentiment_score"`
	// Anchor posts that also mention this subject as a decimal fraction of all anchor posts in the window.
	Share float64 `json:"share"`
	// Canonical identifier for a non-asset person, account, or topic.
	SubjectID string `json:"subject_id"`
}

// SentimentResponse models an object of the `social` service.
type SentimentResponse struct {
	BaseResponse
	// The same measure for other author groups: bots, and top_callers when that population is available. The non-bot population is series itself.
	Cohorts map[string][]SentimentResponseCohortsEntry `json:"cohorts"`
	// Chronologically ordered observations.
	Series []SentimentResponseCohortsEntry `json:"series"`
}

// SentimentResponseCohorts models an object of the `social` service.
type SentimentResponseCohorts struct {
}

// SentimentResponseCohortsEntry models an object of the `social` service.
type SentimentResponseCohortsEntry struct {
	// Bearish posts over the trailing 24 hours ending at this point; interval_bearish_posts is the bucket's.
	BearishPosts float64 `json:"bearish_posts"`
	// Bullish posts over the trailing 24 hours ending at this point; interval_bullish_posts is the bucket's.
	BullishPosts float64 `json:"bullish_posts"`
	// Bearish posts in this emitted time bucket.
	IntervalBearishPosts *float64 `json:"interval_bearish_posts"`
	// Bullish posts in this emitted time bucket.
	IntervalBullishPosts *float64 `json:"interval_bullish_posts"`
	// Neutral posts in this emitted time bucket.
	IntervalNeutralPosts *float64 `json:"interval_neutral_posts"`
	// Posts in this emitted time bucket, before lookback accumulation.
	IntervalPosts *float64 `json:"interval_posts"`
	// Neutral posts over the trailing 24 hours ending at this point; interval_neutral_posts is the bucket's.
	NeutralPosts float64 `json:"neutral_posts"`
	// Decay-weighted net bullish posts: bullish minus bearish. Positive is bullish, negative is bearish, and magnitude grows with activity.
	SentimentSignal float64 `json:"sentiment_signal"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
	// Posts over the trailing 24 hours ending at this point, not this bucket alone; interval_posts is the bucket's.
	TotalPosts float64 `json:"total_posts"`
}

// SentimentTimeframesResponse models an object of the `social` service.
type SentimentTimeframesResponse struct {
	BaseResponse
	// Latest sentiment reading for each available timeframe horizon.
	Current map[string]float64 `json:"current"`
	// Methodology and request-window metadata needed to interpret the response.
	Metadata *SentimentTimeframesResponseMetadata `json:"metadata"`
	// Sentiment series keyed by short_term, medium_term, and long_term horizon.
	Series map[string][]SentimentTimeframesResponseSeriesEntry `json:"series"`
}

// SentimentTimeframesResponseCurrent models an object of the `social` service.
type SentimentTimeframesResponseCurrent struct {
}

// SentimentTimeframesResponseMetadata models an object of the `social` service.
type SentimentTimeframesResponseMetadata struct {
	// UTC time at which the metric was calculated.
	AsOf *time.Time `json:"as_of"`
	// Short name for the calculation method.
	Calculation *string `json:"calculation"`
	// Human-readable formula for the calculation.
	Formula *string `json:"formula"`
	// Despite the name, the length of the linear decay window, in hours: a post counts fully when new, half at the midpoint, and not at all once it is this old.
	HalfLifeHours *float64 `json:"half_life_hours"`
	// Historical lookback used by the calculation, in days.
	LookbackDays *int `json:"lookback_days"`
	// Published scale or interpretation of the metric.
	Scale *string `json:"scale"`
}

// SentimentTimeframesResponseSeries models an object of the `social` service.
type SentimentTimeframesResponseSeries struct {
}

// SentimentTimeframesResponseSeriesEntry models an object of the `social` service.
type SentimentTimeframesResponseSeriesEntry struct {
	// Posts in this emitted time bucket, before lookback accumulation.
	IntervalPosts *float64 `json:"interval_posts"`
	// Decay-weighted net bullish posts: bullish minus bearish. Positive is bullish, negative is bearish, and magnitude grows with activity.
	SentimentSignal float64 `json:"sentiment_signal"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
	// Post count over the metric's active lookback window.
	TotalPosts *float64 `json:"total_posts"`
}

// SocialMomentumResponse models an object of the `social` service.
type SocialMomentumResponse struct {
	BaseResponse
	// Social-momentum series keyed by canonical author-cluster identifier.
	Clusters map[string][]SocialMomentumResponseClustersEntry `json:"clusters"`
}

// SocialMomentumResponseClusters models an object of the `social` service.
type SocialMomentumResponseClusters struct {
}

// SocialMomentumResponseClustersEntry models an object of the `social` service.
type SocialMomentumResponseClustersEntry struct {
	// Arithmetic mean of the posts' sentiment scores, on the −100 (very bearish) to +100 (very bullish) scale.
	AverageSentimentScore float64 `json:"average_sentiment_score"`
	// Posts classified as bearish.
	BearishPosts int64 `json:"bearish_posts"`
	// Posts classified as bullish.
	BullishPosts int64 `json:"bullish_posts"`
	// Posts classified as neutral.
	NeutralPosts int64 `json:"neutral_posts"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
	// Post count over the metric's active lookback window.
	TotalPosts int64 `json:"total_posts"`
	// Distinct posting accounts in the bucket; at coarser-than-hourly granularity this is an upper bound summed from hourly counts.
	UniquePosters int64 `json:"unique_posters"`
}
