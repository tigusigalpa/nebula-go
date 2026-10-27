package nebula

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// TestConcurrentUse runs parallel calls across services to catch data races
// under -race. The client must be safe for concurrent use.
func TestConcurrentUse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/assets":
			_, _ = w.Write([]byte(`{"assets":[],"has_more":false,"limit":50,"offset":0}`))
		default:
			_, _ = w.Write([]byte(`{"plan":"free","history_window_days":null,"rate_limit_per_minute":60}`))
		}
	}))
	defer srv.Close()

	c, err := New("key", WithBaseURL(srv.URL), WithRetry(2, 0, 0))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				_, _ = c.Account.Get(context.Background())
			} else {
				_, _ = c.Lookup.ListAssets(context.Background(), &ListAssetsRequest{Limit: IntPtr(1)})
			}
		}(i)
	}
	wg.Wait()
}
