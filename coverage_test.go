// Service-call matrix for the pinned Nebula Public API contract.
// Field-level request/response audits live in the contract test files.

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
)

// coverageStep shares actual typed service calls with the fixture audits.
type coverageStep struct {
	op   string
	path string
	call func(*testing.T, *Client) (any, error)
}

func coverageSteps(ctx context.Context) []coverageStep {
	return []coverageStep{
		{op: "screener", path: "/screener", call: func(t *testing.T, c *Client) (any, error) {
			return c.Discovery.Screener(ctx, &ScreenerRequest{Kind: KindPtr(KindAsset)})
		}},
		{op: "get_top_momentum", path: "/top-momentum", call: func(t *testing.T, c *Client) (any, error) {
			return c.Discovery.TopMomentum(ctx, &TopMomentumRequest{})
		}},
		{op: "get_top_trending", path: "/top-trending", call: func(t *testing.T, c *Client) (any, error) {
			return c.Discovery.TopTrending(ctx, &TopTrendingRequest{})
		}},
		{op: "list_assets", path: "/assets", call: func(t *testing.T, c *Client) (any, error) {
			return c.Lookup.ListAssets(ctx, &ListAssetsRequest{Query: StringPtr("NVIDIA")})
		}},
		{op: "search", path: "/search", call: func(t *testing.T, c *Client) (any, error) {
			return c.Lookup.Search(ctx, &SearchRequest{Query: "NVIDIA"})
		}},
		{op: "list_subjects", path: "/subjects", call: func(t *testing.T, c *Client) (any, error) {
			return c.Lookup.ListSubjects(ctx, &ListSubjectsRequest{Query: StringPtr("nvidia")})
		}},
		{op: "get_emotions", path: "/emotions", call: func(t *testing.T, c *Client) (any, error) {
			return c.Social.Emotions(ctx, &EmotionsRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "get_mindshare", path: "/mindshare", call: func(t *testing.T, c *Client) (any, error) {
			return c.Social.Mindshare(ctx, &MindshareRequest{Subject: &SubjectRef{ID: "subject_1"}})
		}},
		{op: "get_related_subjects", path: "/related-subjects", call: func(t *testing.T, c *Client) (any, error) {
			return c.Social.RelatedSubjects(ctx, &RelatedSubjectsRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "get_sentiment", path: "/sentiment", call: func(t *testing.T, c *Client) (any, error) {
			return c.Social.Sentiment(ctx, &SentimentRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}, Hours: IntPtr(168)})
		}},
		{op: "get_sentiment_timeframes", path: "/sentiment/timeframes", call: func(t *testing.T, c *Client) (any, error) {
			return c.Social.SentimentTimeframes(ctx, &SentimentTimeframesRequest{Subject: &SubjectRef{ID: "subject_1"}})
		}},
		{op: "get_social_momentum", path: "/social-momentum", call: func(t *testing.T, c *Client) (any, error) {
			return c.Social.SocialMomentum(ctx, &SocialMomentumRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "list_forecast_assets", path: "/forecasts/assets", call: func(t *testing.T, c *Client) (any, error) {
			return c.Forecasts.ListAssets(ctx, &ForecastAssetsRequest{})
		}},
		{op: "get_mindshare_forecast", path: "/forecasts/mindshare", call: func(t *testing.T, c *Client) (any, error) {
			return c.Forecasts.Mindshare(ctx, &MindshareForecastRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "get_sentiment_forecast", path: "/forecasts/sentiment", call: func(t *testing.T, c *Client) (any, error) {
			return c.Forecasts.Sentiment(ctx, &SentimentForecastRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "get_calls", path: "/calls", call: func(t *testing.T, c *Client) (any, error) {
			return c.MarketSentiment.Calls(ctx, &CallsRequest{Clusters: []string{"kols", "media"}})
		}},
		{op: "get_call_returns", path: "/calls/returns", call: func(t *testing.T, c *Client) (any, error) {
			return c.MarketSentiment.CallReturns(ctx, &CallReturnsRequest{})
		}},
		{op: "get_conviction_index", path: "/conviction-index", call: func(t *testing.T, c *Client) (any, error) {
			return c.MarketSentiment.ConvictionIndex(ctx, &ConvictionIndexRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "get_fear_greed_index", path: "/fear-greed-index", call: func(t *testing.T, c *Client) (any, error) {
			return c.MarketSentiment.FearGreedIndex(ctx, &FearGreedIndexRequest{Asset: &AssetRef{ID: "BTC", Class: AssetClassCrypto}})
		}},
		{op: "get_price_levels", path: "/price-levels", call: func(t *testing.T, c *Client) (any, error) {
			return c.MarketSentiment.PriceLevels(ctx, &PriceLevelsRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "get_social_volume_index", path: "/social-volume-index", call: func(t *testing.T, c *Client) (any, error) {
			return c.MarketSentiment.SocialVolumeIndex(ctx, &SocialVolumeIndexRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "get_social_volume_profile", path: "/social-volume-profile", call: func(t *testing.T, c *Client) (any, error) {
			return c.MarketSentiment.SocialVolumeProfile(ctx, &SocialVolumeProfileRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "get_price_history", path: "/price-history", call: func(t *testing.T, c *Client) (any, error) {
			return c.MarketData.PriceHistory(ctx, &PriceHistoryRequest{Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}})
		}},
		{op: "list_calendar_events", path: "/calendar-events", call: func(t *testing.T, c *Client) (any, error) {
			return c.Intelligence.CalendarEvents(ctx, &CalendarEventsRequest{Domains: []CalendarDomain{CalendarDomainCrypto, CalendarDomainEquities}})
		}},
		{op: "list_insights", path: "/insights", call: func(t *testing.T, c *Client) (any, error) {
			return c.Intelligence.Insights(ctx, &InsightsRequest{Tags: []string{"etf", "earnings"}})
		}},
		{op: "get_insight_attention", path: "/insights/attention", call: func(t *testing.T, c *Client) (any, error) {
			return c.Intelligence.InsightAttention(ctx, &InsightAttentionRequest{InsightID: 42})
		}},
		{op: "list_clusters", path: "/clusters", call: func(t *testing.T, c *Client) (any, error) { return c.Clusters.List(ctx, &ClustersRequest{}) }},
		{op: "compare_clusters", path: "/clusters/comparison", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.Comparison(ctx, &ClustersComparisonRequest{AssetClasses: []AssetClass{AssetClassEquity, AssetClassCrypto}})
		}},
		{op: "get_divisive_assets", path: "/clusters/comparison/divisive-assets", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.DivisiveAssets(ctx, &DivisiveAssetsRequest{})
		}},
		{op: "compare_hype_cycles", path: "/clusters/comparison/hype-cycles", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.ComparisonHypeCycles(ctx, &ComparisonHypeCyclesRequest{})
		}},
		{op: "compare_cluster_traits", path: "/clusters/comparison/traits", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.ComparisonTraits(ctx, &ComparisonTraitsRequest{})
		}},
		{op: "get_cluster_hype_cycles", path: "/clusters/hype-cycles", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.HypeCycles(ctx, &ClusterHypeCyclesRequest{ClusterID: ClusterKOLs})
		}},
		{op: "get_cluster_overlap", path: "/clusters/overlap", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.Overlap(ctx, &ClusterOverlapRequest{ClusterID: ClusterKOLs})
		}},
		{op: "get_cluster_profile", path: "/clusters/profile", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.Profile(ctx, &ClusterProfileRequest{ClusterID: ClusterKOLs})
		}},
		{op: "get_cluster_rotation", path: "/clusters/rotation", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.Rotation(ctx, &ClusterRotationRequest{ClusterID: ClusterKOLs})
		}},
		{op: "get_cluster_stances", path: "/clusters/stances", call: func(t *testing.T, c *Client) (any, error) {
			return c.Clusters.Stances(ctx, &ClusterStancesRequest{ClusterID: ClusterKOLs})
		}},
		{op: "create_chat_completion", path: "/chat/completions", call: func(t *testing.T, c *Client) (any, error) {
			return c.Agent.CreateChatCompletion(ctx, &ChatCompletionRequest{Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hi"}}})
		}},
		{op: "get_account", path: "/account", call: func(t *testing.T, c *Client) (any, error) { return c.Account.Get(ctx) }},
		{op: "list_prediction_markets", path: "/predictions", call: func(t *testing.T, c *Client) (any, error) {
			return c.Predictions.ListMarkets(ctx, &PredictionsRequest{Status: PredictionStatusPtr(PredictionStatusOpen)})
		}},
		{op: "get_prediction_history", path: "/predictions/history", call: func(t *testing.T, c *Client) (any, error) {
			return c.Predictions.History(ctx, &PredictionHistoryRequest{ListingID: "l1", OutcomeID: "o1"})
		}},
		{op: "get_prediction_market", path: "/predictions/market", call: func(t *testing.T, c *Client) (any, error) {
			return c.Predictions.Market(ctx, &PredictionMarketRequest{MarketID: "m1"})
		}},
	}
}

// TestCoverage exercises all 41 operations against an HTTP server and checks
// their routes, authentication, fixtures, documentation and x402 eligibility.
func TestCoverage(t *testing.T) {
	ctx := context.Background()
	steps := coverageSteps(ctx)
	// Tie the exercised service calls to the contract, docs, and fixtures. A new
	// upstream operation cannot silently disappear from the test matrix.
	specBytes, err := os.ReadFile("resources/openapi/nebula-openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		Paths map[string]map[string]struct {
			ID       string `json:"operationId"`
			Eligible bool   `json:"x-x402-eligible"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(specBytes, &spec); err != nil {
		t.Fatal(err)
	}
	docBytes, err := os.ReadFile("docs/API-COVERAGE.md")
	if err != nil {
		t.Fatal(err)
	}
	stepPaths := map[string]string{}
	for _, st := range steps {
		if _, duplicate := stepPaths[st.op]; duplicate {
			t.Fatalf("duplicate test for %s", st.op)
		}
		stepPaths[st.op] = st.path
	}
	operations, eligible := 0, 0
	for path, methods := range spec.Paths {
		for method, op := range methods {
			if method != "get" && method != "post" {
				continue
			}
			operations++
			if op.Eligible {
				eligible++
			}
			if stepPaths[op.ID] != path {
				t.Errorf("missing service fixture test for %s %s (%s)", method, path, op.ID)
			}
			if pathToOp[strings.TrimPrefix(path, "/")] != op.ID {
				t.Errorf("fixture mapping missing for %s", op.ID)
			}
			if !strings.Contains(string(docBytes), "`"+op.ID+"`") {
				t.Errorf("docs missing operation %s", op.ID)
			}
			if X402Eligible(op.ID) != op.Eligible {
				t.Errorf("x402 eligibility differs from OpenAPI for %s", op.ID)
			}
		}
	}
	if operations != 41 || len(steps) != operations || eligible != 37 {
		t.Fatalf("operations=%d tests=%d eligible=%d", operations, len(steps), eligible)
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
		if _, err := st.call(t, client); err != nil {
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
