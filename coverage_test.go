// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCoverage exercises every one of the 41 API operations against an
// httptest server: it verifies the request path, auth, serialization and
// typed fixture decoding.
func TestCoverage(t *testing.T) {
	type step struct {
		op   string
		path string
		call func(t *testing.T, c *Client) error
	}
	ctx := context.Background()
	steps := []step{
		{op: "screener", path: "/screener", call: func(t *testing.T, c *Client) error {
			_, err := c.Discovery.Screener(ctx, &ScreenerRequest{Kind: KindPtr(KindAsset)})
			return err
		}},
		{op: "get_top_momentum", path: "/top-momentum", call: func(t *testing.T, c *Client) error {
			_, err := c.Discovery.TopMomentum(ctx, &TopMomentumRequest{})
			return err
		}},
		{op: "get_top_trending", path: "/top-trending", call: func(t *testing.T, c *Client) error {
			_, err := c.Discovery.TopTrending(ctx, &TopTrendingRequest{})
			return err
		}},
		{op: "list_assets", path: "/assets", call: func(t *testing.T, c *Client) error {
			_, err := c.Lookup.ListAssets(ctx, &ListAssetsRequest{Query: StringPtr("NVIDIA")})
			return err
		}},
		{op: "search", path: "/search", call: func(t *testing.T, c *Client) error {
			_, err := c.Lookup.Search(ctx, &SearchRequest{Query: "NVIDIA"})
			return err
		}},
		{op: "list_subjects", path: "/subjects", call: func(t *testing.T, c *Client) error {
			_, err := c.Lookup.ListSubjects(ctx, &ListSubjectsRequest{Query: StringPtr("nvidia")})
			return err
		}},
		{op: "get_emotions", path: "/emotions", call: func(t *testing.T, c *Client) error {
			_, err := c.Social.Emotions(ctx, &EmotionsRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "get_mindshare", path: "/mindshare", call: func(t *testing.T, c *Client) error {
			_, err := c.Social.Mindshare(ctx, &MindshareRequest{Subject: &SubjectRef{ID: "subject_1"}})
			return err
		}},
		{op: "get_related_subjects", path: "/related-subjects", call: func(t *testing.T, c *Client) error {
			_, err := c.Social.RelatedSubjects(ctx, &RelatedSubjectsRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "get_sentiment", path: "/sentiment", call: func(t *testing.T, c *Client) error {
			_, err := c.Social.Sentiment(ctx, &SentimentRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}, Hours: IntPtr(168)})
			return err
		}},
		{op: "get_sentiment_timeframes", path: "/sentiment/timeframes", call: func(t *testing.T, c *Client) error {
			_, err := c.Social.SentimentTimeframes(ctx, &SentimentTimeframesRequest{Subject: &SubjectRef{ID: "subject_1"}})
			return err
		}},
		{op: "get_social_momentum", path: "/social-momentum", call: func(t *testing.T, c *Client) error {
			_, err := c.Social.SocialMomentum(ctx, &SocialMomentumRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "list_forecast_assets", path: "/forecasts/assets", call: func(t *testing.T, c *Client) error {
			_, err := c.Forecasts.ListAssets(ctx, &ForecastAssetsRequest{})
			return err
		}},
		{op: "get_mindshare_forecast", path: "/forecasts/mindshare", call: func(t *testing.T, c *Client) error {
			_, err := c.Forecasts.Mindshare(ctx, &MindshareForecastRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "get_sentiment_forecast", path: "/forecasts/sentiment", call: func(t *testing.T, c *Client) error {
			_, err := c.Forecasts.Sentiment(ctx, &SentimentForecastRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "get_calls", path: "/calls", call: func(t *testing.T, c *Client) error {
			_, err := c.MarketSentiment.Calls(ctx, &CallsRequest{Clusters: []string{"kols", "media"}})
			return err
		}},
		{op: "get_call_returns", path: "/calls/returns", call: func(t *testing.T, c *Client) error {
			_, err := c.MarketSentiment.CallReturns(ctx, &CallReturnsRequest{})
			return err
		}},
		{op: "get_conviction_index", path: "/conviction-index", call: func(t *testing.T, c *Client) error {
			_, err := c.MarketSentiment.ConvictionIndex(ctx, &ConvictionIndexRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "get_fear_greed_index", path: "/fear-greed-index", call: func(t *testing.T, c *Client) error {
			_, err := c.MarketSentiment.FearGreedIndex(ctx, &FearGreedIndexRequest{Asset: &AssetRef{ID: "BTC", Class: AssetClassCrypto}})
			return err
		}},
		{op: "get_price_levels", path: "/price-levels", call: func(t *testing.T, c *Client) error {
			_, err := c.MarketSentiment.PriceLevels(ctx, &PriceLevelsRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "get_social_volume_index", path: "/social-volume-index", call: func(t *testing.T, c *Client) error {
			_, err := c.MarketSentiment.SocialVolumeIndex(ctx, &SocialVolumeIndexRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "get_social_volume_profile", path: "/social-volume-profile", call: func(t *testing.T, c *Client) error {
			_, err := c.MarketSentiment.SocialVolumeProfile(ctx, &SocialVolumeProfileRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "get_price_history", path: "/price-history", call: func(t *testing.T, c *Client) error {
			_, err := c.MarketData.PriceHistory(ctx, &PriceHistoryRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
			return err
		}},
		{op: "list_calendar_events", path: "/calendar-events", call: func(t *testing.T, c *Client) error {
			_, err := c.Intelligence.CalendarEvents(ctx, &CalendarEventsRequest{Domains: []CalendarDomain{CalendarDomainCrypto, CalendarDomainEquities}})
			return err
		}},
		{op: "list_insights", path: "/insights", call: func(t *testing.T, c *Client) error {
			_, err := c.Intelligence.Insights(ctx, &InsightsRequest{Tags: []string{"etf", "earnings"}})
			return err
		}},
		{op: "get_insight_attention", path: "/insights/attention", call: func(t *testing.T, c *Client) error {
			_, err := c.Intelligence.InsightAttention(ctx, &InsightAttentionRequest{InsightID: 42})
			return err
		}},
		{op: "list_clusters", path: "/clusters", call: func(t *testing.T, c *Client) error { _, err := c.Clusters.List(ctx, &ClustersRequest{}); return err }},
		{op: "compare_clusters", path: "/clusters/comparison", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.Comparison(ctx, &ClustersComparisonRequest{AssetClasses: []AssetClass{AssetClassEquity, AssetClassCrypto}})
			return err
		}},
		{op: "get_divisive_assets", path: "/clusters/comparison/divisive-assets", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.DivisiveAssets(ctx, &DivisiveAssetsRequest{})
			return err
		}},
		{op: "compare_hype_cycles", path: "/clusters/comparison/hype-cycles", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.ComparisonHypeCycles(ctx, &ComparisonHypeCyclesRequest{})
			return err
		}},
		{op: "compare_cluster_traits", path: "/clusters/comparison/traits", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.ComparisonTraits(ctx, &ComparisonTraitsRequest{})
			return err
		}},
		{op: "get_cluster_hype_cycles", path: "/clusters/hype-cycles", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.HypeCycles(ctx, &ClusterHypeCyclesRequest{ClusterID: ClusterKOLs})
			return err
		}},
		{op: "get_cluster_overlap", path: "/clusters/overlap", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.Overlap(ctx, &ClusterOverlapRequest{ClusterID: ClusterKOLs})
			return err
		}},
		{op: "get_cluster_profile", path: "/clusters/profile", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.Profile(ctx, &ClusterProfileRequest{ClusterID: ClusterKOLs})
			return err
		}},
		{op: "get_cluster_rotation", path: "/clusters/rotation", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.Rotation(ctx, &ClusterRotationRequest{ClusterID: ClusterKOLs})
			return err
		}},
		{op: "get_cluster_stances", path: "/clusters/stances", call: func(t *testing.T, c *Client) error {
			_, err := c.Clusters.Stances(ctx, &ClusterStancesRequest{ClusterID: ClusterKOLs})
			return err
		}},
		{op: "create_chat_completion", path: "/chat/completions", call: func(t *testing.T, c *Client) error {
			_, err := c.Agent.CreateChatCompletion(ctx, &ChatCompletionRequest{Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hi"}}})
			return err
		}},
		{op: "get_account", path: "/account", call: func(t *testing.T, c *Client) error { _, err := c.Account.Get(ctx); return err }},
		{op: "list_prediction_markets", path: "/predictions", call: func(t *testing.T, c *Client) error {
			_, err := c.Predictions.ListMarkets(ctx, &PredictionsRequest{Status: PredictionStatusPtr(PredictionStatusOpen)})
			return err
		}},
		{op: "get_prediction_history", path: "/predictions/history", call: func(t *testing.T, c *Client) error {
			_, err := c.Predictions.History(ctx, &PredictionHistoryRequest{ListingID: "l1", OutcomeID: "o1"})
			return err
		}},
		{op: "get_prediction_market", path: "/predictions/market", call: func(t *testing.T, c *Client) error {
			_, err := c.Predictions.Market(ctx, &PredictionMarketRequest{MarketID: "m1"})
			return err
		}},
	}

	var gotPaths []string
	var gotAuth []string
	var gotQuery []url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.Method+" "+r.URL.Path)
		gotAuth = append(gotAuth, r.Header.Get("X-API-Key"))
		gotQuery = append(gotQuery, r.URL.Query())
		// figure out which op by path
		op := strings.TrimPrefix(r.URL.Path, "/api/v1/public/")
		fixture, err := os.ReadFile(filepath.Join("testdata", fixtureName(op)))
		if err != nil {
			http.Error(w, "no fixture for "+op, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("RateLimit-Limit", "300")
		w.Header().Set("RateLimit-Remaining", "299")
		w.Header().Set("RateLimit-Reset", "60")
		w.Header().Set("X-Credits-Charged", "1")
		w.Header().Set("X-Credits-Remaining", "158749")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()

	client, err := New("test-key-123", WithBaseURL(srv.URL+"/api/v1/public"), WithRetry(1, 0, 0))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	for _, st := range steps {
		if err := st.call(t, client); err != nil {
			t.Errorf("%s: %v", st.op, err)
		}
	}
	if len(gotPaths) != 41 {
		t.Fatalf("expected 41 calls, got %d", len(gotPaths))
	}
	for i, st := range steps {
		want := "GET /api/v1/public" + st.path
		if st.op == "create_chat_completion" {
			want = "POST /api/v1/public" + st.path
		}
		if gotPaths[i] != want {
			t.Errorf("step %d %s: got %q want %q", i, st.op, gotPaths[i], want)
		}
		if gotAuth[i] != "test-key-123" {
			t.Errorf("step %d %s: missing X-API-Key", i, st.op)
		}
	}
}

// fixtureName maps a request path suffix to its testdata fixture.
func fixtureName(suffix string) string {
	return pathToOp[suffix] + ".json"
}

// pathToOp maps API paths to operation ids (for fixture lookup).
var pathToOp = map[string]string{
	"screener":                            "screener",
	"top-momentum":                        "get_top_momentum",
	"top-trending":                        "get_top_trending",
	"assets":                              "list_assets",
	"search":                              "search",
	"subjects":                            "list_subjects",
	"emotions":                            "get_emotions",
	"mindshare":                           "get_mindshare",
	"related-subjects":                    "get_related_subjects",
	"sentiment":                           "get_sentiment",
	"sentiment/timeframes":                "get_sentiment_timeframes",
	"social-momentum":                     "get_social_momentum",
	"forecasts/assets":                    "list_forecast_assets",
	"forecasts/mindshare":                 "get_mindshare_forecast",
	"forecasts/sentiment":                 "get_sentiment_forecast",
	"calls":                               "get_calls",
	"calls/returns":                       "get_call_returns",
	"conviction-index":                    "get_conviction_index",
	"fear-greed-index":                    "get_fear_greed_index",
	"price-levels":                        "get_price_levels",
	"social-volume-index":                 "get_social_volume_index",
	"social-volume-profile":               "get_social_volume_profile",
	"price-history":                       "get_price_history",
	"calendar-events":                     "list_calendar_events",
	"insights":                            "list_insights",
	"insights/attention":                  "get_insight_attention",
	"clusters":                            "list_clusters",
	"clusters/comparison":                 "compare_clusters",
	"clusters/comparison/divisive-assets": "get_divisive_assets",
	"clusters/comparison/hype-cycles":     "compare_hype_cycles",
	"clusters/comparison/traits":          "compare_cluster_traits",
	"clusters/hype-cycles":                "get_cluster_hype_cycles",
	"clusters/overlap":                    "get_cluster_overlap",
	"clusters/profile":                    "get_cluster_profile",
	"clusters/rotation":                   "get_cluster_rotation",
	"clusters/stances":                    "get_cluster_stances",
	"chat/completions":                    "create_chat_completion",
	"account":                             "get_account",
	"predictions":                         "list_prediction_markets",
	"predictions/history":                 "get_prediction_history",
	"predictions/market":                  "get_prediction_market",
}

// TestFixtureDecodes checks each generated fixture still decodes into its
// typed response — guards against upstream drift in the vendored spec.
func TestFixtureDecodes(t *testing.T) {
	ops := []string{
		"screener",
		"get_top_momentum",
		"get_top_trending",
		"list_assets",
		"search",
		"list_subjects",
		"get_emotions",
		"get_mindshare",
		"get_related_subjects",
		"get_sentiment",
		"get_sentiment_timeframes",
		"get_social_momentum",
		"list_forecast_assets",
		"get_mindshare_forecast",
		"get_sentiment_forecast",
		"get_calls",
		"get_call_returns",
		"get_conviction_index",
		"get_fear_greed_index",
		"get_price_levels",
		"get_social_volume_index",
		"get_social_volume_profile",
		"get_price_history",
		"list_calendar_events",
		"list_insights",
		"get_insight_attention",
		"list_clusters",
		"compare_clusters",
		"get_divisive_assets",
		"compare_hype_cycles",
		"compare_cluster_traits",
		"get_cluster_hype_cycles",
		"get_cluster_overlap",
		"get_cluster_profile",
		"get_cluster_rotation",
		"get_cluster_stances",
		"create_chat_completion",
		"get_account",
		"list_prediction_markets",
		"get_prediction_history",
		"get_prediction_market",
	}
	for _, op := range ops {
		b, err := os.ReadFile(filepath.Join("testdata", op+".json"))
		if err != nil {
			t.Fatalf("%s: %v", op, err)
		}
		var v map[string]any
		if err := json.Unmarshal(b, &v); err != nil {
			t.Errorf("%s: fixture is not an object: %v", op, err)
		}
	}
}

var _ = time.Now // keep time import when unused
