package nebula

import (
	"context"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

// Headers must not cancel the context: real responses may arrive in many reads.
func TestDelayedResponseBody(t *testing.T) {
	c, closeServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		select {
		case <-time.After(20 * time.Millisecond):
			_, _ = io.WriteString(w, `{"plan":"pro"}`)
		case <-r.Context().Done():
		}
	})
	defer closeServer()
	res, err := c.Account.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res.Plan != "pro" {
		t.Fatalf("plan=%q", res.Plan)
	}
}

func TestBodyReadUsesOverallDeadline(t *testing.T) {
	c, closeServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}, WithTimeout(30*time.Millisecond))
	defer closeServer()
	_, err := c.Account.Get(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err=%v", err)
	}
}

func TestCancellationDuringBackoff(t *testing.T) {
	c, closeServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusServiceUnavailable)
	}, WithTimeout(40*time.Millisecond), WithRetry(3, time.Millisecond, time.Second))
	defer closeServer()
	_, err := c.Account.Get(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("lost deadline: %v", err)
	}
}

func TestRequestLocaleOverridesQueryAndMode(t *testing.T) {
	for _, stream := range []bool{false, true} {
		for _, header := range []bool{false, true} {
			c, closeServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if header {
					if r.Header.Get("X-Nebula-Locale") != "ja" || r.URL.Query().Has("locale") {
						t.Errorf("locale=%v %s", r.Header, r.URL)
					}
				} else if r.Header.Get("X-Nebula-Locale") != "" || r.URL.Query().Get("locale") != "ja" {
					t.Errorf("locale=%v %s", r.Header, r.URL)
				}
				if stream {
					w.Header().Set("Content-Type", "text/event-stream")
					_, _ = io.WriteString(w, "data: [DONE]\n\n")
				} else {
					_, _ = io.WriteString(w, `{}`)
				}
			}, WithLocale(LocaleRussian), WithLocaleMode(LocaleModeQuery))
			ctx := WithRequestLocale(context.Background(), LocaleJapanese)
			if header {
				ctx = WithRequestLocaleMode(ctx, LocaleModeHeader)
			}
			if stream {
				s, err := c.Agent.StreamChatCompletion(ctx, streamRequest())
				if err != nil {
					t.Fatal(err)
				}
				for range s.Events() {
				}
				if err := s.Err(); err != nil {
					t.Fatal(err)
				}
			} else if _, err := c.Account.Get(ctx); err != nil {
				t.Fatal(err)
			}
			closeServer()
		}
	}
}

func TestRetryAfterHTTPDatePrecedesReset(t *testing.T) {
	h := http.Header{}
	h.Set("Retry-After", time.Now().Add(4*time.Second).UTC().Format(http.TimeFormat))
	h.Set("RateLimit-Reset", "60")
	err := decodeAPIError(&http.Response{StatusCode: 503, Header: h, Body: io.NopCloser(strings.NewReader(`{}`))}, "test")
	tr := newTransport(defaultConfig(), "key")
	delay := tr.backoff(1, err)
	if delay < 3*time.Second || delay > 4*time.Second {
		t.Fatalf("delay=%s", delay)
	}
}

func TestConcurrentRetryBackoff(t *testing.T) {
	tr := newTransport(defaultConfig(), "key")
	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for attempt := 1; attempt <= 100; attempt++ {
				delay := tr.backoff(attempt, nil)
				if delay < 0 || delay > tr.cfg.retry.MaxDelay {
					t.Errorf("delay=%s", delay)
				}
			}
		}()
	}
	wg.Wait()
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestX402NeverRetriesTransportFailure(t *testing.T) {
	var calls atomic.Int32
	hc := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, io.ErrUnexpectedEOF
	})}
	c, err := NewX402Client(hc, WithRetry(3, time.Nanosecond, time.Nanosecond))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Lookup.ListAssets(context.Background(), nil)
	if err == nil || calls.Load() != 1 {
		t.Fatalf("calls=%d err=%v", calls.Load(), err)
	}
}

func TestRetryOnlyTransientFailures(t *testing.T) {
	tr := newTransport(defaultConfig(), "key")
	for _, tc := range []struct {
		err   error
		retry bool
	}{
		{io.ErrUnexpectedEOF, true},
		{syscall.ECONNRESET, true},
		{&net.DNSError{IsTemporary: true}, true},
		{&net.DNSError{IsNotFound: true}, false},
		{errors.New("permanent failure"), false},
		{context.Canceled, false},
		{&PaymentUncertainError{Err: io.ErrUnexpectedEOF}, false},
	} {
		err := &TransportError{Op: "test", Err: tc.err}
		if got := tr.shouldRetry(http.MethodGet, err); got != tc.retry {
			t.Errorf("%v: retry=%v want %v", tc.err, got, tc.retry)
		}
	}
}

func TestCrossOriginRedirectDoesNotLeakCredentials(t *testing.T) {
	var leaked atomic.Bool
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		leaked.Store(true)
		_, _ = io.WriteString(w, `{}`)
	}))
	defer destination.Close()
	c, closeServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, destination.URL, http.StatusFound) })
	defer closeServer()
	_, err := c.Account.Get(context.Background())
	var ae *APIError
	if !errors.As(err, &ae) || ae.StatusCode != http.StatusFound || leaked.Load() {
		t.Fatalf("leaked=%v err=%v", leaked.Load(), err)
	}
}

func TestInvalidSuccessBodies(t *testing.T) {
	for _, body := range []string{"", "null", "[]", "not JSON"} {
		c, closeServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, body) })
		_, err := c.Account.Get(context.Background())
		closeServer()
		var de *DecodeError
		if !errors.As(err, &de) {
			t.Errorf("body=%q err=%v", body, err)
		}
	}
}

func TestInvalidRequestsNeverReachNetwork(t *testing.T) {
	var calls atomic.Int32
	c, closeServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1); _, _ = io.WriteString(w, `{}`) })
	defer closeServer()
	checks := []func() error{
		func() error { _, err := c.Lookup.Search(context.Background(), nil); return err },
		func() error { _, err := c.Social.Emotions(context.Background(), nil); return err },
		func() error { _, err := c.Forecasts.Mindshare(context.Background(), nil); return err },
		func() error { _, err := c.Clusters.Profile(context.Background(), nil); return err },
		func() error {
			_, err := c.Social.Mindshare(context.Background(), &MindshareRequest{Subject: &SubjectRef{}})
			return err
		},
		func() error {
			_, err := c.Lookup.Search(context.Background(), &SearchRequest{Query: "NVDA", Kind: KindPtr("invalid")})
			return err
		},
		func() error {
			_, err := c.Clusters.List(context.Background(), &ClustersRequest{AssetClass: AssetClassPtr(AssetClassCommodity)})
			return err
		},
		func() error {
			_, err := c.Clusters.Comparison(context.Background(), &ClustersComparisonRequest{AssetClasses: []AssetClass{"invalid"}})
			return err
		},
		func() error {
			_, err := c.Discovery.Screener(context.Background(), &ScreenerRequest{MinShare: Float64Ptr(math.NaN())})
			return err
		},
		func() error {
			_, err := c.Agent.CreateChatCompletion(context.Background(), &ChatCompletionRequest{Messages: []ChatMessage{{Role: "invalid", Content: "hello"}}})
			return err
		},
		func() error {
			req := streamRequest()
			req.ThinkingMode = ThinkingModePtr("invalid")
			_, err := c.Agent.CreateChatCompletion(context.Background(), req)
			return err
		},
	}
	for i, check := range checks {
		if err := check(); !errors.Is(err, ErrValidation) {
			t.Errorf("case %d: %v", i, err)
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("HTTP calls=%d", calls.Load())
	}
}

func TestClassFiltersAndHTTPMetadata(t *testing.T) {
	c, closeServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/emotions" && (r.URL.Query().Get("asset_class") != "crypto" || r.URL.Query().Has("asset_id")) {
			t.Errorf("query=%s", r.URL.RawQuery)
		}
		w.Header().Set("X-Request-Id", "req-filter")
		if r.URL.Path == "/sentiment" {
			_, _ = io.WriteString(w, `{"series":[]}`)
		} else {
			_, _ = io.WriteString(w, `{"assets":[],"series":{},"metadata":{"hours":24}}`)
		}
	})
	defer closeServer()
	ctx := context.Background()
	if _, err := c.Lookup.ListAssets(ctx, &ListAssetsRequest{AssetClass: AssetClassPtr(AssetClassCrypto)}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Social.Sentiment(ctx, &SentimentRequest{AssetClass: AssetClassPtr(AssetClassCrypto)}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Social.SocialMomentum(ctx, &SocialMomentumRequest{AssetClass: AssetClassPtr(AssetClassCrypto)}); err != nil {
		t.Fatal(err)
	}
	res, err := c.Social.Emotions(ctx, &EmotionsRequest{Subject: &SubjectRef{ID: "text:example"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Metadata == nil || res.Metadata.Hours != 24 || res.HTTPMetadata().RequestID != "req-filter" {
		t.Fatalf("body=%+v HTTP=%+v", res.Metadata, res.HTTPMetadata())
	}
	if _, err := c.Social.Sentiment(ctx, &SentimentRequest{AssetClass: AssetClassPtr(AssetClassCrypto), Subject: &SubjectRef{ID: "text:example"}}); !errors.Is(err, ErrValidation) {
		t.Fatalf("selector err=%v", err)
	}
}
