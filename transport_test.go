package nebula

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) (*Client, func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	opts = append([]Option{WithBaseURL(srv.URL), WithRetry(1, 0, 0)}, opts...)
	c, err := New("test-key", opts...)
	if err != nil {
		srv.Close()
		t.Fatalf("New: %v", err)
	}
	return c, srv.Close
}

func TestAPIKeyAuthHeader(t *testing.T) {
	var got string
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("X-API-Key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"plan":"pro","history_window_days":3650,"rate_limit_per_minute":300}`))
	})
	defer closeFn()
	if _, err := c.Account.Get(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got != "test-key" {
		t.Fatalf("X-API-Key = %q", got)
	}
}

func TestBearerAuthHeader(t *testing.T) {
	var gotAuth, gotKey string
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotKey = r.Header.Get("X-API-Key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"plan":"pro","history_window_days":null,"rate_limit_per_minute":null}`))
	}, WithAuthMode(AuthModeBearer))
	defer closeFn()
	if _, err := c.Account.Get(context.Background()); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer test-key" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if gotKey != "" {
		t.Fatalf("X-API-Key must not be sent in bearer mode, got %q", gotKey)
	}
}

func TestUserAgentHeader(t *testing.T) {
	var got string
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("User-Agent")
		_, _ = w.Write([]byte(`{"plan":"free","history_window_days":null,"rate_limit_per_minute":60}`))
	})
	defer closeFn()
	if _, err := c.Account.Get(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got != "nebula-go/"+Version {
		t.Fatalf("User-Agent = %q", got)
	}
}

func TestLocaleHeaderAndQuery(t *testing.T) {
	var gotHeader, gotQuery string
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Nebula-Locale")
		gotQuery = r.URL.Query().Get("locale")
		_, _ = w.Write([]byte(`{"plan":"free","history_window_days":null,"rate_limit_per_minute":60}`))
	}, WithLocale(LocaleRussian))
	defer closeFn()
	if _, err := c.Account.Get(context.Background()); err != nil {
		t.Fatal(err)
	}
	if gotHeader != "ru" || gotQuery != "" {
		t.Fatalf("header locale=%q query=%q", gotHeader, gotQuery)
	}

	gotHeader, gotQuery = "", ""
	c2, close2 := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Nebula-Locale")
		gotQuery = r.URL.Query().Get("locale")
		_, _ = w.Write([]byte(`{"plan":"free","history_window_days":null,"rate_limit_per_minute":60}`))
	}, WithLocale(LocaleJapanese), WithLocaleMode(LocaleModeQuery))
	defer close2()
	if _, err := c2.Account.Get(context.Background()); err != nil {
		t.Fatal(err)
	}
	if gotHeader != "" || gotQuery != "ja" {
		t.Fatalf("query mode: header=%q query=%q", gotHeader, gotQuery)
	}
}

func TestPerRequestLocale(t *testing.T) {
	var got string
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("X-Nebula-Locale")
		_, _ = w.Write([]byte(`{"plan":"free","history_window_days":null,"rate_limit_per_minute":60}`))
	})
	defer closeFn()
	ctx := WithRequestLocale(context.Background(), LocaleSpanish)
	if _, err := c.Account.Get(ctx); err != nil {
		t.Fatal(err)
	}
	if got != "es" {
		t.Fatalf("per-request locale header = %q", got)
	}
}

func TestMetadataParsing(t *testing.T) {
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("RateLimit-Limit", "300")
		w.Header().Set("RateLimit-Remaining", "42")
		w.Header().Set("RateLimit-Reset", "17")
		w.Header().Set("X-Credits-Charged", "5")
		w.Header().Set("X-Credits-Remaining", "1000")
		w.Header().Set("X-Request-Id", "req_abc")
		_, _ = w.Write([]byte(`{"plan":"pro","history_window_days":3650,"rate_limit_per_minute":300}`))
	})
	defer closeFn()
	res, err := c.Account.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	m := res.Metadata
	if m.RateLimitLimit == nil || *m.RateLimitLimit != 300 {
		t.Fatalf("RateLimitLimit=%v", m.RateLimitLimit)
	}
	if m.RateLimitRemaining == nil || *m.RateLimitRemaining != 42 {
		t.Fatalf("RateLimitRemaining=%v", m.RateLimitRemaining)
	}
	if m.CreditsCharged == nil || *m.CreditsCharged != 5 {
		t.Fatalf("CreditsCharged=%v", m.CreditsCharged)
	}
	if m.RequestID != "req_abc" {
		t.Fatalf("RequestID=%q", m.RequestID)
	}
	if m.StatusCode != 200 {
		t.Fatalf("StatusCode=%d", m.StatusCode)
	}
	if res.Raw["plan"] != "pro" {
		t.Fatalf("Raw missing plan: %v", res.Raw)
	}
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		status   int
		body     string
		sentinel error
	}{
		{400, `{"error":true,"message":"bad param"}`, ErrBadRequest},
		{401, `{"error":true,"message":"revoked"}`, ErrUnauthorized},
		{402, `{"error":true,"message":"out of credits","plan_renews_at":"2026-11-01T00:00:00Z","purchased_credits":100}`, ErrOutOfCredits},
		{403, `{"error":true,"message":"history beyond plan"}`, ErrForbidden},
		{404, `{"error":true,"message":"nope"}`, ErrNotFound},
		{405, `{"error":true,"message":"use GET"}`, ErrMethodNotAllowed},
		{429, `{"error":true,"message":"slow down"}`, ErrRateLimited},
		{500, `{"error":true,"message":"boom"}`, ErrServer},
		{503, `{"error":true,"message":"dependency down"}`, ErrServiceUnavailable},
	}
	for _, tc := range cases {
		c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(tc.status)
			_, _ = w.Write([]byte(tc.body))
		})
		_, err := c.Account.Get(context.Background())
		closeFn()
		if err == nil {
			t.Fatalf("status %d: expected error", tc.status)
		}
		if !errors.Is(err, tc.sentinel) {
			t.Fatalf("status %d: errors.Is(%v) failed, got %v", tc.status, tc.sentinel, err)
		}
		var ae *APIError
		if !errors.As(err, &ae) {
			t.Fatalf("status %d: not APIError", tc.status)
		}
		if ae.StatusCode != tc.status {
			t.Fatalf("status %d != %d", ae.StatusCode, tc.status)
		}
	}
	// 402 detail fields
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(402)
		_, _ = w.Write([]byte(`{"error":true,"message":"out","plan_renews_at":"2026-11-01T00:00:00Z","purchased_credits":100}`))
	})
	defer closeFn()
	_, err := c.Account.Get(context.Background())
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatal("expected APIError")
	}
	if ae.PlanRenewsAt != "2026-11-01T00:00:00Z" {
		t.Fatalf("PlanRenewsAt=%q", ae.PlanRenewsAt)
	}
	if ae.PurchasedCredits == nil || *ae.PurchasedCredits != 100 {
		t.Fatalf("PurchasedCredits=%v", ae.PurchasedCredits)
	}
}

func TestRetryOn429HonorsRetryAfter(t *testing.T) {
	var calls int32
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			_, _ = w.Write([]byte(`{"error":true,"message":"slow"}`))
			return
		}
		_, _ = w.Write([]byte(`{"plan":"free","history_window_days":null,"rate_limit_per_minute":60}`))
	}, WithRetry(3, time.Millisecond, 10*time.Millisecond))
	defer closeFn()
	if _, err := c.Account.Get(context.Background()); err != nil {
		t.Fatalf("expected success after retry: %v", err)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("calls=%d", calls)
	}
}

func TestRetryOn503(t *testing.T) {
	var calls int32
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(503)
			_, _ = w.Write([]byte(`{"error":true,"message":"dep down"}`))
			return
		}
		_, _ = w.Write([]byte(`{"plan":"free","history_window_days":null,"rate_limit_per_minute":60}`))
	}, WithRetry(4, time.Millisecond, 10*time.Millisecond))
	defer closeFn()
	if _, err := c.Account.Get(context.Background()); err != nil {
		t.Fatalf("expected success: %v", err)
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Fatalf("calls=%d", calls)
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	var calls int32
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error":true,"message":"bad"}`))
	}, WithRetry(5, time.Millisecond, 10*time.Millisecond))
	defer closeFn()
	if _, err := c.Account.Get(context.Background()); err == nil {
		t.Fatal("expected error")
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestNoRetryOnPost(t *testing.T) {
	var calls int32
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(503)
		_, _ = w.Write([]byte(`{"error":true,"message":"down"}`))
	}, WithRetry(5, time.Millisecond, 10*time.Millisecond))
	defer closeFn()
	req := &ChatCompletionRequest{Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hi"}}}
	_, err := c.Agent.CreateChatCompletion(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("POST must not retry, calls=%d", calls)
	}
}

func TestContextCancellation(t *testing.T) {
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
		_, _ = w.Write([]byte(`{}`))
	})
	defer closeFn()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := c.Account.Get(ctx)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}

func TestTimeoutCoversRetries(t *testing.T) {
	var calls int32
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(503)
		_, _ = w.Write([]byte(`{"error":true,"message":"down"}`))
	}, WithRetry(10, 200*time.Millisecond, 5*time.Second), WithTimeout(120*time.Millisecond))
	defer closeFn()
	start := time.Now()
	_, err := c.Account.Get(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("timeout exceeded: %v", elapsed)
	}
	if atomic.LoadInt32(&calls) < 1 {
		t.Fatal("no calls")
	}
}

func TestX402EligibilityGuard(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	c, err := NewX402Client(nil, WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	// non-eligible ops must fail locally
	if _, err := c.Account.Get(context.Background()); !errors.Is(err, ErrNotX402Eligible) {
		t.Fatalf("account: %v", err)
	}
	if _, err := c.Forecasts.Mindshare(context.Background(), &MindshareForecastRequest{Asset: &AssetRef{ID: "BTC", Class: AssetClassCrypto}}); !errors.Is(err, ErrNotX402Eligible) {
		t.Fatalf("forecast: %v", err)
	}
	req := &ChatCompletionRequest{Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hi"}}}
	if _, err := c.Agent.CreateChatCompletion(context.Background(), req); !errors.Is(err, ErrNotX402Eligible) {
		t.Fatalf("chat: %v", err)
	}
	if atomic.LoadInt32(&calls) != 0 {
		t.Fatalf("local guard must prevent HTTP calls, got %d", calls)
	}
	if !c.X402() {
		t.Fatal("X402() should report true")
	}
}

func TestX402NoAPIKeySent(t *testing.T) {
	var gotKey, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-API-Key")
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"assets":[],"has_more":false,"limit":50,"offset":0}`))
	}))
	defer srv.Close()
	c, err := NewX402Client(nil, WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Lookup.ListAssets(context.Background(), &ListAssetsRequest{}); err != nil {
		t.Fatal(err)
	}
	if gotKey != "" || gotAuth != "" {
		t.Fatalf("x402 mode must not send credentials: %q %q", gotKey, gotAuth)
	}
}

func TestQuerySerialization(t *testing.T) {
	var gotURL string
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.String()
		_, _ = w.Write([]byte(`{"cohorts":{},"series":[]}`))
	})
	defer closeFn()
	_, err := c.Social.Sentiment(context.Background(), &SentimentRequest{
		Asset:    &AssetRef{ID: "NVDA", Class: AssetClassEquity},
		Hours:    IntPtr(168),
		Clusters: []string{"kols", "media"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotURL, "asset_id=NVDA") || !strings.Contains(gotURL, "asset_class=equity") {
		t.Fatalf("asset pair missing: %s", gotURL)
	}
	if !strings.Contains(gotURL, "clusters=kols") || !strings.Contains(gotURL, "clusters=media") {
		t.Fatalf("clusters must be repeated params: %s", gotURL)
	}
	if !strings.Contains(gotURL, "hours=168") {
		t.Fatalf("hours missing: %s", gotURL)
	}
}

func TestCSVSerialization(t *testing.T) {
	var got string
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query().Get("domains")
		_, _ = w.Write([]byte(`{"calendar_events":[]}`))
	})
	defer closeFn()
	_, err := c.Intelligence.CalendarEvents(context.Background(), &CalendarEventsRequest{
		Domains: []CalendarDomain{CalendarDomainCrypto, CalendarDomainEquities},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != "crypto,equities" {
		t.Fatalf("domains must be CSV: %q", got)
	}
}

func TestValidationBlocksBeforeHTTP(t *testing.T) {
	var calls int32
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`{}`))
	})
	defer closeFn()
	// emotions requires a selector
	_, err := c.Social.Emotions(context.Background(), &EmotionsRequest{})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	// asset+subject mixing
	_, err = c.Social.Sentiment(context.Background(), &SentimentRequest{
		Asset:   &AssetRef{ID: "NVDA", Class: AssetClassEquity},
		Subject: &SubjectRef{ID: "s1"},
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	// asset_id + query on /assets
	_, err = c.Lookup.ListAssets(context.Background(), &ListAssetsRequest{
		Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity},
		Query: StringPtr("x"),
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	// screener: price filter on subject kind
	_, err = c.Discovery.Screener(context.Background(), &ScreenerRequest{
		Kind:                 KindPtr(KindSubject),
		MinPriceChange24hPct: Float64Ptr(1),
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	// min > max
	_, err = c.Discovery.Screener(context.Background(), &ScreenerRequest{
		MinShare: Float64Ptr(0.9), MaxShare: Float64Ptr(0.1),
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	// communities on crypto
	_, err = c.Discovery.Screener(context.Background(), &ScreenerRequest{
		AssetClass: AssetClassPtr(AssetClassCrypto),
		Clusters:   []string{"wallstreetbets"},
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	// out of range hours
	_, err = c.Social.Sentiment(context.Background(), &SentimentRequest{Hours: IntPtr(99999)})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if atomic.LoadInt32(&calls) != 0 {
		t.Fatalf("validation must happen before HTTP, got %d calls", calls)
	}
}

func TestForecastNotForecastableDecodes(t *testing.T) {
	c, closeFn := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"not_forecastable","reason":"too_little_history","series":[],"sample_paths":[],"refreshed":false,"horizon_hours":72}`))
	})
	defer closeFn()
	res, err := c.Forecasts.Mindshare(context.Background(), &MindshareForecastRequest{
		Asset: &AssetRef{ID: "NEW", Class: AssetClassCrypto},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != ForecastStatusNotForecastable {
		t.Fatalf("status=%q", res.Status)
	}
	if res.Reason == nil || *res.Reason != ForecastReasonTooLittleHistory {
		t.Fatalf("reason=%v", res.Reason)
	}
}
