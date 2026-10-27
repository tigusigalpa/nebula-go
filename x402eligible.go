package nebula

// x402EligibleOps mirrors x-x402-eligible in the vendored OpenAPI spec
// (v1.0.2): the 37 fixed-price GET operations served on /api/v1/x402.
// /account, both forecasts, and POST /chat/completions are excluded.
var x402EligibleOps = map[string]bool{
	"list_assets":               true,
	"list_calendar_events":      true,
	"get_calls":                 true,
	"get_call_returns":          true,
	"list_clusters":             true,
	"compare_clusters":          true,
	"get_divisive_assets":       true,
	"compare_hype_cycles":       true,
	"compare_cluster_traits":    true,
	"get_cluster_hype_cycles":   true,
	"get_cluster_overlap":       true,
	"get_cluster_profile":       true,
	"get_cluster_rotation":      true,
	"get_cluster_stances":       true,
	"get_conviction_index":      true,
	"get_emotions":              true,
	"get_fear_greed_index":      true,
	"list_forecast_assets":      true,
	"list_insights":             true,
	"get_insight_attention":     true,
	"get_mindshare":             true,
	"list_prediction_markets":   true,
	"get_prediction_history":    true,
	"get_prediction_market":     true,
	"get_price_history":         true,
	"get_price_levels":          true,
	"get_related_subjects":      true,
	"screener":                  true,
	"search":                    true,
	"get_sentiment":             true,
	"get_sentiment_timeframes":  true,
	"get_social_momentum":       true,
	"get_social_volume_index":   true,
	"get_social_volume_profile": true,
	"list_subjects":             true,
	"get_top_momentum":          true,
	"get_top_trending":          true,
}

// x402Eligible reports whether op is served on the x402 base path.
func x402Eligible(op string) bool { return x402EligibleOps[op] }

// X402Eligible reports whether an operation id is served on the x402 surface.
// Exported for callers building eligibility UI or guards.
func X402Eligible(op string) bool { return x402EligibleOps[op] }
