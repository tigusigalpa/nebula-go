package nebula

import (
	"net/url"
	"strconv"
	"strings"
)

// validateRelations keeps rules described in prose separate from generated
// schema constraints. OpenAPI 3.0 parameters cannot express these relationships.
func validateRelations(op string, q url.Values) error {
	invalid := func(field, msg string) error { return &ValidationError{Op: op, Field: field, Msg: msg} }
	if q.Has("asset_id") {
		if strings.TrimSpace(q.Get("asset_id")) == "" || strings.TrimSpace(q.Get("asset_class")) == "" {
			return invalid("asset_id", "asset identifier and class must both be non-empty")
		}
		if q.Has("subject_id") {
			return invalid("subject_id", "subject and asset selectors are mutually exclusive")
		}
	}
	if q.Has("subject_id") {
		if strings.TrimSpace(q.Get("subject_id")) == "" {
			return invalid("subject_id", "must be non-empty")
		}
		if q.Has("asset_class") {
			return invalid("asset_class", "cannot be combined with a subject selector")
		}
	}
	switch op {
	case "get_emotions", "get_mindshare", "get_related_subjects", "get_sentiment_timeframes":
		if !q.Has("asset_id") && !q.Has("subject_id") {
			return invalid("selector", "asset or subject selector is required")
		}
	case "list_assets":
		if q.Has("asset_id") && q.Has("query") {
			return invalid("query", "asset_id cannot be combined with query")
		}
	case "list_subjects":
		if q.Has("subject_id") && q.Has("query") {
			return invalid("query", "subject_id cannot be combined with query")
		}
	case "screener", "get_top_momentum", "get_top_trending":
		if q.Get("kind") == string(KindSubject) && q.Has("asset_class") {
			return invalid("asset_class", "applies only to kind=asset")
		}
		if op == "screener" {
			if q.Get("kind") == string(KindSubject) && (q.Get("sort") == string(ScreenerSortPriceChange24hPct) || q.Has("min_price_change_24h_pct") || q.Has("max_price_change_24h_pct")) {
				return invalid("price_change_24h_pct", "price filters and sorting require kind=asset")
			}
			for _, field := range []string{"share", "share_change", "net_bullish_posts", "price_change_24h_pct"} {
				lo, hi := "min_"+field, "max_"+field
				if q.Has(lo) && q.Has(hi) {
					minValue, _ := strconv.ParseFloat(q.Get(lo), 64)
					maxValue, _ := strconv.ParseFloat(q.Get(hi), 64)
					if minValue > maxValue {
						return invalid(lo, "must not exceed "+hi)
					}
				}
			}
		}
	}
	if hasCommunity(q["clusters"]) && (q.Get("kind") == string(KindSubject) || q.Has("subject_id") || q.Get("asset_class") == string(AssetClassCrypto)) {
		return invalid("clusters", "posting communities require assets in equity, commodity or index")
	}
	if isCommunity(q.Get("cluster_id")) && q.Has("asset_class") {
		for _, class := range strings.Split(q.Get("asset_class"), ",") {
			if class != string(AssetClassEquity) {
				return invalid("asset_class", "posting communities apply only to equity")
			}
		}
	}
	return nil
}
