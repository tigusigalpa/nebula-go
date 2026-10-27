// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// CalendarEventsResponse models an object of the `intelligence` service.
type CalendarEventsResponse struct {
	BaseResponse
	// Calendar events matching the supplied filters.
	CalendarEvents []CalendarEventsResponseCalendarEvent `json:"calendar_events"`
}

// CalendarEventsResponseCalendarEvent models an object of the `intelligence` service.
type CalendarEventsResponseCalendarEvent struct {
	// Stable public identifier for the calendar event.
	CalendarEventID string `json:"calendar_event_id"`
	// Confidence in the event, from 1 (low) to 10 (high).
	ConfidenceScore int `json:"confidence_score"`
	// UTC time when the record was created.
	CreatedAt time.Time `json:"created_at"`
	// Market or subject domains attached to the record.
	Domains []string `json:"domains"`
	// UTC event end time when known.
	EndTime *time.Time `json:"end_time"`
	// Assets and non-asset subjects linked to this record.
	Entities []CalendarEventsResponseCalendarEventEntity `json:"entities"`
	// Event-type-specific structured metadata. Keys vary by event type.
	EventMetadata map[string]any `json:"event_metadata"`
	// More specific event classification when available.
	EventSubtype *string `json:"event_subtype"`
	// High-level event classification.
	EventType string `json:"event_type"`
	// Importance from 1 (low) to 3 (high).
	ImportanceScore int `json:"importance_score"`
	// Where the event physically takes place, when a place is known. Absent for online and on-chain events.
	Location *string `json:"location"`
	// Estimated long-term market impact, from 1 (low) to 3 (high).
	LongTermImpactScore int `json:"long_term_impact_score"`
	// Estimated macroeconomic impact, from 1 (low) to 3 (high).
	MacroImpactScore int `json:"macro_impact_score"`
	// Estimated direct price impact, from 1 (low) to 3 (high).
	PriceImpactScore int `json:"price_impact_score"`
	// Estimated short-term market impact, from 1 (low) to 3 (high).
	ShortTermImpactScore int `json:"short_term_impact_score"`
	// Number of sources supporting the event.
	SourceCount int `json:"source_count"`
	// Source records supporting the parent response.
	Sources []CalendarEventsResponseCalendarEventSource `json:"sources"`
	// UTC time when the intelligence item or event starts.
	StartTime time.Time `json:"start_time"`
	// Current lifecycle status of the event: scheduled, tentative, cancelled, or postponed.
	Status string `json:"status"`
	// Evidence-grounded summary of the intelligence item.
	Summary *string `json:"summary"`
	// Searchable labels attached to the record.
	Tags []string `json:"tags"`
	// Estimated tail-risk impact, from 1 (low) to 3 (high).
	TailRiskImpactScore int `json:"tail_risk_impact_score"`
	// Precision of the published event time: minute, hour, day, week, month, year, or unknown. A day-precision event starts at 00:00 UTC that day.
	TimePrecision string `json:"time_precision"`
	// Human-readable event title.
	Title string `json:"title"`
	// Number of supporting sources marked trusted.
	TrustedSourceCount int `json:"trusted_source_count"`
	// UTC time when the record was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// CalendarEventsResponseCalendarEventEntity models an object of the `intelligence` service.
type CalendarEventsResponseCalendarEventEntity struct {
	// The asset's class, sent with asset_id.
	AssetClass *AssetClass `json:"asset_class"`
	// For a stock, commodity, index, or coin, the asset_id the asset endpoints take; absent for people, accounts, and topics, and for the rare asset Nebula cannot place.
	AssetID *string `json:"asset_id"`
	// Entity-link confidence when available.
	Confidence *float64 `json:"confidence"`
	// Nebula's id for the linked entity, such as eq:NVDA, gecko:1, or text:ai. To query the entity, use asset_id and asset_class, or subject_id, from the same object.
	EntityID string `json:"entity_id"`
	// True when this is the record's primary linked entity.
	IsPrimary bool `json:"is_primary"`
	// Human-readable name.
	Name string `json:"name"`
	// Relationship of the linked entity to the calendar event.
	Role string `json:"role"`
	// For a person, account, or topic, the subject_id the subject-aware endpoints take; absent for assets.
	SubjectID *string `json:"subject_id"`
}

// CalendarEventsResponseCalendarEventEventMetadata models an object of the `intelligence` service.
type CalendarEventsResponseCalendarEventEventMetadata struct {
}

// CalendarEventsResponseCalendarEventSource models an object of the `intelligence` service.
type CalendarEventsResponseCalendarEventSource struct {
	// Identifier assigned by the external source.
	ExternalID string `json:"external_id"`
	// Upstream data provider when available.
	Provider *PredictionProvider `json:"provider"`
	// UTC publication time.
	PublishedAt time.Time `json:"published_at"`
	// Original publisher when available.
	Publisher *string `json:"publisher"`
	// Source platform or feed.
	Source string `json:"source"`
	// Source subtype.
	Subtype string `json:"subtype"`
	// Whether the source is classified as trusted.
	Trusted bool `json:"trusted"`
	// Public source URL when available.
	URL *string `json:"url"`
}

// InsightsResponse models an object of the `intelligence` service.
type InsightsResponse struct {
	BaseResponse
	// Curated intelligence items matching the supplied filters.
	Insights []InsightsResponseInsight `json:"insights"`
}

// InsightsResponseInsight models an object of the `intelligence` service.
type InsightsResponseInsight struct {
	// Primary intelligence category when assigned.
	Category *string `json:"category"`
	// Market or subject domains attached to the record.
	Domains []string `json:"domains"`
	// Assets and non-asset subjects linked to this record.
	Entities []InsightsResponseInsightEntity `json:"entities"`
	// Concise intelligence headline.
	Headline string `json:"headline"`
	// Importance from 1 (low) to 3 (high).
	ImportanceScore int `json:"importance_score"`
	// Stable identifier for the intelligence item.
	InsightID int64 `json:"insight_id"`
	// Signed sentiment score; positive is bullish and negative is bearish.
	SentimentScore int `json:"sentiment_score"`
	// UTC time when the intelligence item or event starts.
	StartTime *time.Time `json:"start_time"`
	// Evidence-grounded summary of the intelligence item.
	Summary string `json:"summary"`
	// Searchable labels attached to the record.
	Tags []string `json:"tags"`
	// UTC time when the record was last updated.
	UpdatedAt *time.Time `json:"updated_at"`
}

// InsightsResponseInsightEntity models an object of the `intelligence` service.
type InsightsResponseInsightEntity struct {
	// The asset's class, sent with asset_id.
	AssetClass *AssetClass `json:"asset_class"`
	// For a stock, commodity, index, or coin, the asset_id the asset endpoints take; absent for people, accounts, and topics, and for the rare asset Nebula cannot place.
	AssetID *string `json:"asset_id"`
	// Nebula's id for the linked entity, such as eq:NVDA, gecko:1, or text:ai. To query the entity, use asset_id and asset_class, or subject_id, from the same object.
	EntityID string `json:"entity_id"`
	// Knowledge-graph classifications for the linked entity.
	EntityTypes []string `json:"entity_types"`
	// True when the record directly concerns this linked entity.
	IsDirect bool `json:"is_direct"`
	// True when this is the record's primary linked entity.
	IsPrimary bool `json:"is_primary"`
	// Human-readable name.
	Name string `json:"name"`
	// Relationship of the linked entity to the intelligence item.
	Role string `json:"role"`
	// For a person, account, or topic, the subject_id the subject-aware endpoints take; absent for assets.
	SubjectID *string `json:"subject_id"`
}

// InsightAttentionResponse models an object of the `intelligence` service.
type InsightAttentionResponse struct {
	BaseResponse
	// Conversation and impact readout for the requested intelligence item.
	Attention *InsightAttentionResponseAttention `json:"attention"`
	// Stable identifier for the intelligence item.
	InsightID int64 `json:"insight_id"`
	// Price reactions for priceable entities linked to the intelligence item.
	PriceReactions []InsightAttentionResponsePriceReaction `json:"price_reactions"`
}

// InsightAttentionResponseAttention models an object of the `intelligence` service.
type InsightAttentionResponseAttention struct {
	// Subject id used to anchor the attention scan.
	AnchorSubjectID string `json:"anchor_subject_id"`
	// Distinct authors represented by the matching posts.
	Authors int64 `json:"authors"`
	// Arithmetic mean of the posts' sentiment scores, on the −100 (very bearish) to +100 (very bullish) scale.
	AverageSentimentScore float64 `json:"average_sentiment_score"`
	// Posts classified as bearish.
	BearishPosts int64 `json:"bearish_posts"`
	// Width of each bucket in hours.
	BucketHours int `json:"bucket_hours"`
	// Attention buckets in chronological order.
	Buckets []InsightAttentionResponseAttentionBucket `json:"buckets"`
	// Posts classified as bullish.
	BullishPosts int64 `json:"bullish_posts"`
	// Author-cluster mix for the matching posts, ranked by authors.
	ClusterMix []InsightAttentionResponseAttentionClusterMixItem `json:"cluster_mix"`
	// Distinct author clusters represented by the matching posts.
	Clusters int64 `json:"clusters"`
	// Inclusive UTC start of the measured window.
	From time.Time `json:"from"`
	// Structured impact assessment for the intelligence item.
	Impact *InsightAttentionResponseAttentionImpact `json:"impact"`
	// UTC time of the latest matching post in the requested window.
	LastPostAt *time.Time `json:"last_post_at"`
	// Posts classified as neutral.
	NeutralPosts int64 `json:"neutral_posts"`
	// Subject ids where matching posts had to mention at least one; empty means no optional side was required.
	OptionalSubjects []string `json:"optional_subjects"`
	// Matching post count over the requested window.
	Posts int64 `json:"posts"`
	// Subject ids that all matching posts had to mention.
	RequiredSubjects []string `json:"required_subjects"`
	// Subject ids actually used for the returned attention window after any fallback.
	SubjectsUsed []string `json:"subjects_used"`
	// Exclusive UTC end of the measured window.
	To time.Time `json:"to"`
	// Window placement around the insight: before or around.
	WindowKind string `json:"window_kind"`
}

// InsightAttentionResponseAttentionBucket models an object of the `intelligence` service.
type InsightAttentionResponseAttentionBucket struct {
	// Distinct authors represented by the matching posts.
	Authors int64 `json:"authors"`
	// Arithmetic mean of the posts' sentiment scores, on the −100 (very bearish) to +100 (very bullish) scale.
	AverageSentimentScore float64 `json:"average_sentiment_score"`
	// Posts classified as bearish.
	BearishPosts int64 `json:"bearish_posts"`
	// Posts classified as bullish.
	BullishPosts int64 `json:"bullish_posts"`
	// Posts classified as neutral.
	NeutralPosts int64 `json:"neutral_posts"`
	// Matching post count over the requested window.
	Posts int64 `json:"posts"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
}

// InsightAttentionResponseAttentionClusterMixItem models an object of the `intelligence` service.
type InsightAttentionResponseAttentionClusterMixItem struct {
	// Distinct authors represented by the matching posts.
	Authors int64 `json:"authors"`
	// Canonical author-cluster identifier.
	ClusterID string `json:"cluster_id"`
	// Matching post count over the requested window.
	Posts *int64 `json:"posts"`
	// Share of eligible social attention as a decimal fraction from 0 to 1.
	Share float64 `json:"share"`
}

// InsightAttentionResponseAttentionImpact models an object of the `intelligence` service.
type InsightAttentionResponseAttentionImpact struct {
	// Maximum possible score for each impact vector.
	MaxScore int `json:"max_score"`
	// Signed sentiment score; positive is bullish and negative is bearish.
	SentimentScore int `json:"sentiment_score"`
	// Impact vectors behind the item-level impact assessment.
	Vectors []InsightAttentionResponseAttentionImpactVector `json:"vectors"`
}

// InsightAttentionResponseAttentionImpactVector models an object of the `intelligence` service.
type InsightAttentionResponseAttentionImpactVector struct {
	// Stable machine-readable key for this impact vector.
	Key string `json:"key"`
	// Human-readable bucket or index label.
	Label string `json:"label"`
	// Score for this impact vector.
	Score int `json:"score"`
}

// InsightAttentionResponsePriceReaction models an object of the `intelligence` service.
type InsightAttentionResponsePriceReaction struct {
	// UTC time used as the baseline for the price reaction.
	InsightAt time.Time `json:"insight_at"`
	// Stable identifier for the intelligence item.
	InsightID int64 `json:"insight_id"`
	// Latest price change versus the insight baseline, present for recent insights.
	LivePct *float64 `json:"live_pct"`
	// Fixed-horizon price moves once each horizon has elapsed.
	Moves []InsightAttentionResponsePriceReactionMove `json:"moves"`
	// Chronologically ordered observations.
	Series []InsightAttentionResponsePriceReactionSeriesPoint `json:"series"`
	// Canonical identifier for a non-asset person, account, or topic.
	SubjectID string `json:"subject_id"`
}

// InsightAttentionResponsePriceReactionMove models an object of the `intelligence` service.
type InsightAttentionResponsePriceReactionMove struct {
	// Percentage-point price change versus the baseline at the insight time.
	ChangePct float64 `json:"change_pct"`
	// Lookback window represented by the response, in hours.
	Hours int `json:"hours"`
}

// InsightAttentionResponsePriceReactionSeriesPoint models an object of the `intelligence` service.
type InsightAttentionResponsePriceReactionSeriesPoint struct {
	// Percentage-point price change versus the baseline at the insight time.
	ChangePct float64 `json:"change_pct"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
}
