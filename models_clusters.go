// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// ClustersResponse models an object of the `clusters` service.
type ClustersResponse struct {
	BaseResponse
	// Author-cluster summaries for the selected asset class.
	Clusters []ClustersResponseCluster `json:"clusters"`
	// Lookback window represented by the response, in hours.
	Hours int `json:"hours"`
}

// ClustersResponseCluster models an object of the `clusters` service.
type ClustersResponseCluster struct {
	// Percent change in posts against the previous window of equal length (76 posts after 84 is -9.52); absent when no comparison is possible.
	ActivityChangePct *float64 `json:"activity_change_pct"`
	// Posts classified as bearish.
	BearishPosts float64 `json:"bearish_posts"`
	// Posts classified as bullish.
	BullishPosts float64 `json:"bullish_posts"`
	// Canonical author-cluster identifier.
	ClusterID string `json:"cluster_id"`
	// Posts classified as neutral.
	NeutralPosts float64 `json:"neutral_posts"`
	// Post count in the immediately preceding window of equal length.
	PreviousTotalPosts float64 `json:"previous_total_posts"`
	// The latest value of the hourly sentiment series: the most recent hour with posts.
	RecentSentimentSignal float64 `json:"recent_sentiment_signal"`
	// Chronologically ordered observations.
	Series []ClustersResponseClusterSeriesPoint `json:"series"`
	// Post count over the metric's active lookback window.
	TotalPosts float64 `json:"total_posts"`
	// Net bullish sentiment signal over the complete requested window.
	WindowSentimentSignal float64 `json:"window_sentiment_signal"`
}

// ClustersResponseClusterSeriesPoint models an object of the `clusters` service.
type ClustersResponseClusterSeriesPoint struct {
	// Posts classified as bearish.
	BearishPosts float64 `json:"bearish_posts"`
	// Posts classified as bullish.
	BullishPosts float64 `json:"bullish_posts"`
	// Decay-weighted net bullish posts: bullish minus bearish. Positive is bullish, negative is bearish, and magnitude grows with activity.
	SentimentSignal float64 `json:"sentiment_signal"`
	// UTC RFC 3339 timestamp for the observation or bucket start.
	Timestamp time.Time `json:"timestamp"`
	// Post count over the metric's active lookback window.
	TotalPosts float64 `json:"total_posts"`
}

// ClustersComparisonResponse models an object of the `clusters` service.
type ClustersComparisonResponse struct {
	BaseResponse
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// Every cluster with posts in the window, most posts first.
	Clusters []ClustersComparisonResponseCluster `json:"clusters"`
	// Window length in complete UTC days, ending at the latest complete day.
	Days int `json:"days"`
}

// ClustersComparisonResponseCluster models an object of the `clusters` service.
type ClustersComparisonResponseCluster struct {
	// Percent change in posts against the previous window of equal length (76 posts after 84 is -9.52); absent when no comparison is possible.
	ActivityChangePct *float64 `json:"activity_change_pct"`
	// The cluster's bearish posts over the window.
	BearishPosts int64 `json:"bearish_posts"`
	// The cluster's bullish posts over the window.
	BullishPosts int64 `json:"bullish_posts"`
	// The cluster's bullish ÷ (bullish + bearish) posts over the window, from 0 to 1. Null with no bullish or bearish posts.
	ClusterBullishShare *float64 `json:"cluster_bullish_share"`
	// Canonical author-cluster identifier.
	ClusterID string `json:"cluster_id"`
	// The same share for the market the cluster reads, over the same days and asset classes (the equity market for a posting community): the baseline cluster_bullish_share compares with.
	MarketBullishShare *float64 `json:"market_bullish_share"`
	// The cluster's posts over the window.
	Posts int64 `json:"posts"`
	// Every complete UTC day of the window, oldest first, days without posts included.
	Series []ClustersComparisonResponseClusterSeriesPoint `json:"series"`
	// The asset the cluster turned to more than the market did over the last 7 complete days; null when there is none.
	TrendingAsset *ClustersComparisonResponseClusterTrendingAsset `json:"trending_asset"`
}

// ClustersComparisonResponseClusterSeriesPoint models an object of the `clusters` service.
type ClustersComparisonResponseClusterSeriesPoint struct {
	// The cluster's bearish posts on the day.
	BearishPosts int64 `json:"bearish_posts"`
	// The cluster's bullish posts on the day.
	BullishPosts int64 `json:"bullish_posts"`
	// The cluster's posts on the day.
	Posts int64 `json:"posts"`
	// UTC start of the day: 00:00.
	Timestamp time.Time `json:"timestamp"`
}

// ClustersComparisonResponseClusterTrendingAsset models an object of the `clusters` service.
type ClustersComparisonResponseClusterTrendingAsset struct {
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID string `json:"asset_id"`
	// Human-readable name.
	Name string `json:"name"`
	// Display ticker or market symbol.
	Symbol string `json:"symbol"`
}

// DivisiveAssetsResponse models an object of the `clusters` service.
type DivisiveAssetsResponse struct {
	BaseResponse
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// The divisive assets, most evenly split first, up to limit.
	Assets []DivisiveAssetsResponseAsset `json:"assets"`
	// Window length in complete UTC days, ending at the latest complete day.
	Days int `json:"days"`
}

// DivisiveAssetsResponseAsset models an object of the `clusters` service.
type DivisiveAssetsResponseAsset struct {
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID string `json:"asset_id"`
	// Clusters leaning bearish on the asset by 0.1 or more.
	BearishClusterCount int `json:"bearish_cluster_count"`
	// Clusters leaning bullish on the asset by 0.1 or more.
	BullishClusterCount int `json:"bullish_cluster_count"`
	// Every cluster with a stance on the asset (8 or more bullish or bearish posts), most bearish first.
	ClusterStances []DivisiveAssetsResponseAssetClusterStance `json:"cluster_stances"`
	// The market's net stance on the same asset over the same window, on the same −1 to +1 scale, from every author except bots.
	MarketStance float64 `json:"market_stance"`
	// The smaller side's share of the clusters leaning clearly bullish or bearish (0.1 or more either way), from 0 (all one way) to 0.5 (an even split). Assets are ranked by it.
	MinorityShare float64 `json:"minority_share"`
	// Human-readable name.
	Name string `json:"name"`
	// The interquartile range of the clusters' stances, from 0 to 2: how far apart the middle half of the clusters sit. Breaks ties in minority_share.
	StanceIqr float64 `json:"stance_iqr"`
	// Display ticker or market symbol.
	Symbol string `json:"symbol"`
}

// DivisiveAssetsResponseAssetClusterStance models an object of the `clusters` service.
type DivisiveAssetsResponseAssetClusterStance struct {
	// Canonical author-cluster identifier.
	ClusterID string `json:"cluster_id"`
	// The cluster's posts on the asset over the window.
	Posts int64 `json:"posts"`
	// The cluster's net stance on the asset: (bullish − bearish) ÷ (bullish + bearish) posts, from −1 to +1.
	Stance float64 `json:"stance"`
}

// ComparisonHypeCyclesResponse models an object of the `clusters` service.
type ComparisonHypeCyclesResponse struct {
	BaseResponse
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// The window's biggest hype cycles, most recent peak first, up to limit.
	Cycles []ComparisonHypeCyclesResponseCycle `json:"cycles"`
	// Window length in complete UTC days, ending at the latest complete day.
	Days int `json:"days"`
}

// ComparisonHypeCyclesResponseCycle models an object of the `clusters` service.
type ComparisonHypeCyclesResponseCycle struct {
	// Every cluster with 3 or more posts in the 48 hours around the peak hour, earliest first.
	Arrivals []ComparisonHypeCyclesResponseCycleArrival `json:"arrivals"`
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID string `json:"asset_id"`
	// Human-readable name.
	Name string `json:"name"`
	// The market's posts on the asset on the peak day, bots left out.
	PeakDayPosts int64 `json:"peak_day_posts"`
	// UTC start of the market's busiest hour on the peak day.
	PeakTimestamp time.Time `json:"peak_timestamp"`
	// The market's posts on the asset hour by hour, from 24 hours before the peak hour to 23 hours after it: 48 observations.
	Series []ComparisonHypeCyclesResponseCycleSeriesPoint `json:"series"`
	// Display ticker or market symbol.
	Symbol string `json:"symbol"`
}

// ComparisonHypeCyclesResponseCycleArrival models an object of the `clusters` service.
type ComparisonHypeCyclesResponseCycleArrival struct {
	// Canonical author-cluster identifier.
	ClusterID string `json:"cluster_id"`
	// Hours from the market's peak hour to the posts-weighted average hour of the cluster's posts in the 48 hours around it. Negative is before the peak.
	OffsetHours float64 `json:"offset_hours"`
	// The cluster's posts on the asset in the 48 hours around the peak hour.
	Posts int64 `json:"posts"`
}

// ComparisonHypeCyclesResponseCycleSeriesPoint models an object of the `clusters` service.
type ComparisonHypeCyclesResponseCycleSeriesPoint struct {
	// The market's posts on the asset in the hour, bots left out.
	Posts int `json:"posts"`
	// UTC start of the hour.
	Timestamp time.Time `json:"timestamp"`
}

// ComparisonTraitsResponse models an object of the `clusters` service.
type ComparisonTraitsResponse struct {
	BaseResponse
	// Clusters active in the window: the population traits are ranked within.
	ActiveClusters int `json:"active_clusters"`
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// Window length in complete UTC days, ending at the latest complete day.
	Days int `json:"days"`
	// Market hype cycles in the window: assets whose posts peaked at 50 or more in a day and at least 3× their median day.
	HypeCycles int `json:"hype_cycles"`
	// One entry per trait, in the order earliness, contrarian, hype_fade, conviction, focus, breadth, bullishness, discovery, turnover, coverage.
	Traits []ComparisonTraitsResponseTrait `json:"traits"`
}

// ComparisonTraitsResponseTrait models an object of the `clusters` service.
type ComparisonTraitsResponseTrait struct {
	// Every cluster with a value for the trait, highest first.
	Clusters []ComparisonTraitsResponseTraitCluster `json:"clusters"`
	// The trait measured; the operation description defines each one.
	Trait string `json:"trait"`
	// days, share (0 to 1) or count; fixed per trait.
	Unit string `json:"unit"`
}

// ComparisonTraitsResponseTraitCluster models an object of the `clusters` service.
type ComparisonTraitsResponseTraitCluster struct {
	// Canonical author-cluster identifier.
	ClusterID string `json:"cluster_id"`
	// The cluster's value for the trait, in unit.
	Value float64 `json:"value"`
}

// ClusterHypeCyclesResponse models an object of the `clusters` service.
type ClusterHypeCyclesResponse struct {
	BaseResponse
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// The cluster analysed, as its canonical cluster_id.
	ClusterID string `json:"cluster_id"`
	// Every cluster with a lead (5 or more cycles), earliest first, the requested one included.
	Clusters []ClusterHypeCyclesResponseCluster `json:"clusters"`
	// The average cumulative build-up, one point per day_offset from −7 to +7; empty when no cycle has completed.
	Curve []ClusterHypeCyclesResponseCurveItem `json:"curve"`
	// Completed hype cycles the cluster joined: 3 or more posts within 7 days of the peak.
	CyclesJoined int `json:"cycles_joined"`
	// Window length in complete UTC days, ending at the latest complete day.
	Days int `json:"days"`
	// How far ahead of everyone else the cluster's posts arrive within 24 hours of the market's peak hour, in days. Positive is early; null under 20 cycles.
	LeadDays *float64 `json:"lead_days"`
}

// ClusterHypeCyclesResponseCluster models an object of the `clusters` service.
type ClusterHypeCyclesResponseCluster struct {
	// The cluster, as its canonical cluster_id.
	ClusterID string `json:"cluster_id"`
	// Completed hype cycles in which the cluster posted within 24 hours of the market's peak hour.
	CyclesJoined int `json:"cycles_joined"`
	// The cluster's lead in days, computed as the top-level lead_days. Positive is early.
	LeadDays float64 `json:"lead_days"`
	// High end of the likely range of lead_days.
	LeadHighDays float64 `json:"lead_high_days"`
	// Low end of the likely range of lead_days.
	LeadLowDays float64 `json:"lead_low_days"`
}

// ClusterHypeCyclesResponseCurveItem models an object of the `clusters` service.
type ClusterHypeCyclesResponseCurveItem struct {
	// Percentage of the cluster's posts in the 15 days around the peak made by the end of this day, averaged over cycles; 100 at day_offset 7.
	ClusterCumulativePct float64 `json:"cluster_cumulative_pct"`
	// Days from the market peak: negative before it, 0 on the peak day, positive after it.
	DayOffset int `json:"day_offset"`
	// The same for the market, over the same cycles.
	MarketCumulativePct float64 `json:"market_cumulative_pct"`
}

// ClusterOverlapResponse models an object of the `clusters` service.
type ClusterOverlapResponse struct {
	BaseResponse
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// The cluster analysed, as its canonical cluster_id.
	ClusterID string `json:"cluster_id"`
	// Window length in complete UTC days, ending at the latest complete day.
	Days int `json:"days"`
	// Every other active cluster, most-overlapping first.
	Peers []ClusterOverlapResponsePeer `json:"peers"`
	// Size of each cluster's focus set: its most-discussed assets, up to this many.
	TopAssetCount int `json:"top_asset_count"`
}

// ClusterOverlapResponsePeer models an object of the `clusters` service.
type ClusterOverlapResponsePeer struct {
	// Compared assets on which both clusters lean the same way relative to the market.
	AgreeingAssetCount int `json:"agreeing_asset_count"`
	// Canonical author-cluster identifier.
	ClusterID string `json:"cluster_id"`
	// Shared assets, across the whole overlap, on which both clusters lean at least 0.05 from the market's stance.
	ComparedAssetCount int `json:"compared_asset_count"`
	// On shared assets where both clusters lean at least 0.05 from the market's stance, the share where they lean the same way: agreeing_asset_count ÷ compared_asset_count. Null under 3
	DirectionAgreementShare *float64 `json:"direction_agreement_share"`
	// Share of this cluster's focus set also in the peer's, from 0 to 1.
	OverlapShare float64 `json:"overlap_share"`
	// Up to 5 assets in both focus sets, in this cluster's order of attention.
	SharedAssets []ClusterOverlapResponsePeerSharedAsset `json:"shared_assets"`
	// This cluster's stance minus the peer's, averaged over shared assets both have a stance on, from −2 to +2. Positive means this cluster is more bullish. Null under 3 assets.
	StanceGap *float64 `json:"stance_gap"`
}

// ClusterOverlapResponsePeerSharedAsset models an object of the `clusters` service.
type ClusterOverlapResponsePeerSharedAsset struct {
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID string `json:"asset_id"`
	// This cluster's net stance on the asset, from −1 to +1. Null under 3 bullish or bearish posts.
	ClusterStance *float64 `json:"cluster_stance"`
	// The market's net stance on the asset, from −1 to +1: the baseline both leans are measured against. Null under 3 bullish or bearish posts.
	MarketStance *float64 `json:"market_stance"`
	// Human-readable name.
	Name string `json:"name"`
	// The peer's net stance on the asset, from −1 to +1. Null under 3 bullish or bearish posts.
	PeerStance *float64 `json:"peer_stance"`
	// Display ticker or market symbol.
	Symbol string `json:"symbol"`
}

// ClusterProfileResponse models an object of the `clusters` service.
type ClusterProfileResponse struct {
	BaseResponse
	// Clusters active in the window: the population traits are ranked within.
	ActiveClusters int `json:"active_clusters"`
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// The cluster analysed, as its canonical cluster_id.
	ClusterID string `json:"cluster_id"`
	// Window length in complete UTC days, ending at the latest complete day.
	Days int `json:"days"`
	// Market hype cycles in the window: assets whose posts peaked at 50 or more in a day and at least 3× their median day.
	HypeCycles int `json:"hype_cycles"`
	// One entry per trait, in the order earliness, contrarian, hype_fade, conviction, focus, breadth, bullishness, discovery, turnover, coverage.
	Traits []ClusterProfileResponseTrait `json:"traits"`
}

// ClusterProfileResponseTrait models an object of the `clusters` service.
type ClusterProfileResponseTrait struct {
	// Every other cluster's value for this trait, ascending.
	PeerValues []float64 `json:"peer_values"`
	// From 0 (lowest) to 1 (highest) among clusters with a value. Null when value is null.
	Percentile *float64 `json:"percentile"`
	// 1 is the highest value; tied clusters share the better rank. Null when value is null.
	Rank *int `json:"rank"`
	// The trait measured; the operation description defines each one.
	Trait string `json:"trait"`
	// days, share (0 to 1) or count; fixed per trait.
	Unit string `json:"unit"`
	// The cluster's value for this trait, in unit; null when the cluster has too little data for it.
	Value *float64 `json:"value"`
}

// ClusterRotationResponse models an object of the `clusters` service.
type ClusterRotationResponse struct {
	BaseResponse
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// The cluster analysed, as its canonical cluster_id.
	ClusterID string `json:"cluster_id"`
	// Assets whose share of the cluster's posts rose most, largest gain first. Only assets whose share rose are listed.
	Gainers []ClusterRotationResponseGainer `json:"gainers"`
	// Lookback window represented by the response, in hours.
	Hours int `json:"hours"`
	// Assets whose share of the cluster's posts fell most, largest fall first. Only assets whose share fell are listed.
	Losers []ClusterRotationResponseGainer `json:"losers"`
	// Each topic's share of the current window and its change against the previous one, in the order equity, commodity, index, crypto, people, events.
	TopicChanges []ClusterRotationResponseTopicChange `json:"topic_changes"`
	// Each topic's share of the cluster's posts per complete UTC day, oldest first.
	TopicSeries []ClusterRotationResponseTopicSeriesPoint `json:"topic_series"`
}

// ClusterRotationResponseGainer models an object of the `clusters` service.
type ClusterRotationResponseGainer struct {
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID string `json:"asset_id"`
	// Posts as a fraction of all the cluster's posts about the selected asset classes in the current window, from 0 to 1.
	CurrentShare float64 `json:"current_share"`
	// Human-readable name.
	Name string `json:"name"`
	// The asset's posts as a fraction of all the cluster's posts in the comparison window, from 0 to 1.
	PreviousShare float64 `json:"previous_share"`
	// current_share minus the share in the comparison window, in basis points: 100 is one percentage point. Negative means the share fell.
	ShareChangeBps int `json:"share_change_bps"`
	// Display ticker or market symbol.
	Symbol string `json:"symbol"`
}

// ClusterRotationResponseTopicChange models an object of the `clusters` service.
type ClusterRotationResponseTopicChange struct {
	// The topic's share of the cluster's posts in the current window, from 0 to 1.
	CurrentShare float64 `json:"current_share"`
	// Change in share against the previous window, in basis points (100 is one percentage point).
	ShareChangeBps int `json:"share_change_bps"`
	// equity, commodity, index, crypto, people (people and their accounts) or events (places, themes and events).
	Topic string `json:"topic"`
}

// ClusterRotationResponseTopicSeriesPoint models an object of the `clusters` service.
type ClusterRotationResponseTopicSeriesPoint struct {
	// UTC start of the day: 00:00.
	Timestamp time.Time `json:"timestamp"`
	// Each topic's share of the cluster's posts that day, in the order of topic_changes. The shares add up to 1 when the cluster posted that day.
	Topics []ClusterRotationResponseTopicSeriesPointTopic `json:"topics"`
}

// ClusterRotationResponseTopicSeriesPointTopic models an object of the `clusters` service.
type ClusterRotationResponseTopicSeriesPointTopic struct {
	// equity, commodity, index, crypto, people (people and their accounts) or events (places, themes and events).
	Topic string `json:"topic"`
	// The topic's share of the cluster's posts that day, from 0 to 1.
	TopicShare float64 `json:"topic_share"`
}

// ClusterStancesResponse models an object of the `clusters` service.
type ClusterStancesResponse struct {
	BaseResponse
	// The markets analysed, in the order equity, commodity, index, crypto: the asset_class values requested, or all four when asset_class was omitted.
	AssetClasses []string `json:"asset_classes"`
	// The cluster's most-discussed eligible assets (3 or more cluster posts, 20 or more market posts), up to limit.
	Assets []ClusterStancesResponseAsset `json:"assets"`
	// The cluster analysed, as its canonical cluster_id.
	ClusterID string `json:"cluster_id"`
	// Share of eligible assets on which the cluster and the market lean opposite ways, both at least 0.1 from neutral. Null under 5 assets.
	ContrarianShare *float64 `json:"contrarian_share"`
	// Window length in complete UTC days, ending at the latest complete day.
	Days int `json:"days"`
	// contrarian_share averaged over the other clusters; null when none has one.
	PeerAverageContrarianShare *float64 `json:"peer_average_contrarian_share"`
}

// ClusterStancesResponseAsset models an object of the `clusters` service.
type ClusterStancesResponseAsset struct {
	// The asset's class: equity (stocks and ETFs), commodity, index, or crypto. Send it with asset_id on every asset-scoped request.
	AssetClass AssetClass `json:"asset_class"`
	// Canonical public asset identifier: the ticker or market symbol for stocks, commodities, and indices, and the CoinGecko slug (or, without one, a numeric id) for crypto.
	AssetID string `json:"asset_id"`
	// The cluster's posts about the asset over the window.
	ClusterPosts int64 `json:"cluster_posts"`
	// The cluster's net stance on the asset over the window: (bullish − bearish) / (bullish + bearish) posts, from −1 (all bearish) to +1 (all bullish).
	ClusterStance float64 `json:"cluster_stance"`
	// The market's net stance on the same asset over the same window, on the same −1 to +1 scale, from every author except bots.
	MarketStance float64 `json:"market_stance"`
	// Human-readable name.
	Name string `json:"name"`
	// Display ticker or market symbol.
	Symbol string `json:"symbol"`
}
