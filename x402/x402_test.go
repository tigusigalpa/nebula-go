package x402

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	nebula "github.com/tigusigalpa/nebula-go"
	evmmech "github.com/x402-foundation/x402/go/v2/mechanisms/evm"
)

// fakeSigner implements evmmech.ClientEvmSigner without any real key or
// network. It records that a signature was requested; the signature bytes are
// throwaway test data.
type fakeSigner struct {
	address string
	signed  int32
}

func (f *fakeSigner) Address() string { return f.address }

func (f *fakeSigner) SignTypedData(_ context.Context, domain evmmech.TypedDataDomain, types map[string][]evmmech.TypedDataField, primaryType string, message map[string]interface{}) ([]byte, error) {
	atomic.AddInt32(&f.signed, 1)
	return []byte("test-signature-not-real"), nil
}

func offerHeader(t *testing.T, resourceURL, network, asset, amount, payTo string) string {
	t.Helper()
	offer := map[string]any{
		"x402Version": 2,
		"resource":    map[string]any{"url": resourceURL},
		"accepts": []map[string]any{{
			"scheme":            "exact",
			"network":           network,
			"asset":             asset,
			"amount":            amount,
			"payTo":             payTo,
			"maxTimeoutSeconds": 60,
			"extra":             map[string]any{"name": "USD Coin", "version": "2"},
		}},
	}
	b, err := json.Marshal(offer)
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// testServer mimics the Nebula x402 surface: 402 offer then settled data.
func testServer(t *testing.T, network, asset, amount, payTo string, onSigned func(sig string)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("PAYMENT-SIGNATURE") == "" {
			resource := "http://" + r.Host + r.URL.String()
			w.Header().Set("PAYMENT-REQUIRED", offerHeader(t, resource, network, asset, amount, payTo))
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		onSigned(r.Header.Get("PAYMENT-SIGNATURE"))
		settle, _ := json.Marshal(map[string]any{"success": true, "transaction": "0xtesttx", "network": Network, "payer": "0xpayer"})
		w.Header().Set("PAYMENT-RESPONSE", base64.StdEncoding.EncodeToString(settle))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"assets":[{"asset_id":"NVDA","asset_class":"equity","name":"NVIDIA","symbol":"NVDA"}],"has_more":false,"limit":50,"offset":0}`))
	}))
}

func testPolicy(srv *httptest.Server) SpendPolicy {
	return SpendPolicy{
		MaxUSDCPerPayment: "0.01",
		AllowedHosts:      []string{serverHost(srv)},
	}
}

func serverHost(srv *httptest.Server) string {
	return strings.TrimPrefix(srv.URL, "http://")
}

func TestOfferSignatureReceiptFlow(t *testing.T) {
	var sawSignature int32
	srv := testServer(t, Network, USDCAsset, "1250", "0xRecipient1234", func(sig string) {
		atomic.AddInt32(&sawSignature, 1)
		if sig == "" {
			t.Error("PAYMENT-SIGNATURE must not be empty")
		}
	})
	defer srv.Close()

	signer := &fakeSigner{address: "0xBuyer"}
	httpClient, err := NewHTTPClient(nil, signer, testPolicy(srv))
	if err != nil {
		t.Fatal(err)
	}
	client, err := nebula.NewX402Client(httpClient, nebula.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Lookup.ListAssets(context.Background(), &nebula.ListAssetsRequest{Query: nebula.StringPtr("NVIDIA")})
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if atomic.LoadInt32(&sawSignature) != 1 {
		t.Fatalf("expected one signed retry, got %d", sawSignature)
	}
	if res.Metadata.PaymentReceipt == nil {
		t.Fatal("PAYMENT-RESPONSE receipt missing from metadata")
	}
	var receipt map[string]any
	if err := json.Unmarshal(res.Metadata.PaymentReceipt, &receipt); err != nil {
		t.Fatalf("receipt decode: %v", err)
	}
	if receipt["transaction"] != "0xtesttx" {
		t.Fatalf("receipt: %v", receipt)
	}
}

func TestPolicyRejectsWrongNetwork(t *testing.T) {
	srv := testServer(t, "eip155:1", USDCAsset, "1250", "0xR", func(string) {
		t.Error("signature must never happen for a wrong-network offer")
	})
	defer srv.Close()
	signer := &fakeSigner{address: "0xBuyer"}
	httpClient, err := NewHTTPClient(nil, signer, testPolicy(srv))
	if err != nil {
		t.Fatal(err)
	}
	client, err := nebula.NewX402Client(httpClient, nebula.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Lookup.ListAssets(context.Background(), &nebula.ListAssetsRequest{})
	if err == nil {
		t.Fatal("expected payment rejection")
	}
	if atomic.LoadInt32(&signer.signed) != 0 {
		t.Fatal("signer must not be called for a policy-rejected offer")
	}
}

func TestPolicyRejectsOverspend(t *testing.T) {
	srv := testServer(t, Network, USDCAsset, "999999999", "0xR", func(string) {
		t.Error("signature must never happen over cap")
	})
	defer srv.Close()
	signer := &fakeSigner{address: "0xBuyer"}
	httpClient, err := NewHTTPClient(nil, signer, testPolicy(srv))
	if err != nil {
		t.Fatal(err)
	}
	client, err := nebula.NewX402Client(httpClient, nebula.WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Lookup.ListAssets(context.Background(), &nebula.ListAssetsRequest{})
	if err == nil {
		t.Fatal("expected overspend rejection")
	}
	if atomic.LoadInt32(&signer.signed) != 0 {
		t.Fatal("signer must not be called over cap")
	}
}

func TestSpendPolicyCheck(t *testing.T) {
	p := SpendPolicy{
		MaxUSDCPerPayment:  "0.01",
		AllowedRecipients:  []string{"0xAllowed"},
		AllowedHosts:       []string{"nebula-api.hiddensystems.ai"},
		RequireExactScheme: true,
	}
	good := &PaymentOffer{
		ResourceURL: "https://nebula-api.hiddensystems.ai/api/v1/x402/assets",
		Network:     Network, Scheme: "exact", Asset: USDCAsset,
		Amount: mustBig("1250"), PayTo: "0xAllowed",
	}
	if err := p.Check(good); err != nil {
		t.Fatalf("good offer rejected: %v", err)
	}
	bad := *good
	bad.PayTo = "0xOther"
	if err := p.Check(&bad); err == nil {
		t.Fatal("recipient check failed")
	}
	bad = *good
	bad.Asset = "0xDeadBeef"
	if err := p.Check(&bad); err == nil {
		t.Fatal("asset check failed")
	}
	bad = *good
	bad.Amount = mustBig("99999999")
	if err := p.Check(&bad); err == nil {
		t.Fatal("amount check failed")
	}
	bad = *good
	bad.ResourceURL = "https://evil.example.com/api"
	if err := p.Check(&bad); err == nil {
		t.Fatal("resource host check failed")
	}
	bad = *good
	bad.Scheme = "upto"
	if err := p.Check(&bad); err == nil {
		t.Fatal("scheme check failed")
	}
}

func TestPreviewOffer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("PAYMENT-REQUIRED", offerHeader(t, r.URL.String(), Network, USDCAsset, "2500", "0xR"))
		w.WriteHeader(http.StatusPaymentRequired)
	}))
	defer srv.Close()
	offer, err := PreviewOffer(context.Background(), nil, http.MethodGet, srv.URL+"/api/v1/x402/assets")
	if err != nil {
		t.Fatal(err)
	}
	if offer.Network != Network || offer.Asset != USDCAsset {
		t.Fatalf("offer: %+v", offer)
	}
	if offer.Amount == nil || offer.Amount.String() != "2500" {
		t.Fatalf("amount: %v", offer.Amount)
	}
}

func TestMaxAtomicValidation(t *testing.T) {
	if _, err := NewHTTPClient(nil, &fakeSigner{}, SpendPolicy{}); err == nil {
		t.Fatal("empty cap must be rejected")
	}
	if _, err := NewHTTPClient(nil, &fakeSigner{}, SpendPolicy{MaxUSDCPerPayment: "0.0000001"}); err == nil {
		t.Fatal("sub-micro cap must be rejected")
	}
	if _, err := NewHTTPClient(nil, &fakeSigner{}, SpendPolicy{MaxUSDCPerPayment: "0.01"}); err != nil {
		t.Fatalf("valid cap rejected: %v", err)
	}
}

func TestEligibilityLocal(t *testing.T) {
	if !nebula.X402Eligible("list_assets") {
		t.Fatal("list_assets must be eligible")
	}
	for _, op := range []string{"get_account", "get_mindshare_forecast", "get_sentiment_forecast", "create_chat_completion"} {
		if nebula.X402Eligible(op) {
			t.Fatalf("%s must not be eligible", op)
		}
	}
	// 37 eligible operations
	count := 0
	for _, op := range ops41() {
		if nebula.X402Eligible(op) {
			count++
		}
	}
	if count != 37 {
		t.Fatalf("expected 37 eligible ops, got %d", count)
	}
}

func ops41() []string {
	return []string{
		"screener", "get_top_momentum", "get_top_trending",
		"list_assets", "search", "list_subjects",
		"get_emotions", "get_mindshare", "get_related_subjects", "get_sentiment", "get_sentiment_timeframes", "get_social_momentum",
		"list_forecast_assets", "get_mindshare_forecast", "get_sentiment_forecast",
		"get_calls", "get_call_returns", "get_conviction_index", "get_fear_greed_index", "get_price_levels", "get_social_volume_index", "get_social_volume_profile",
		"get_price_history",
		"list_calendar_events", "list_insights", "get_insight_attention",
		"list_clusters", "compare_clusters", "get_divisive_assets", "compare_hype_cycles", "compare_cluster_traits",
		"get_cluster_hype_cycles", "get_cluster_overlap", "get_cluster_profile", "get_cluster_rotation", "get_cluster_stances",
		"create_chat_completion", "get_account",
		"list_prediction_markets", "get_prediction_history", "get_prediction_market",
	}
}

func mustBig(s string) *big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	return n
}

// decodePaymentResponse sanity
func TestDecodePaymentResponse(t *testing.T) {
	h := http.Header{}
	settle, _ := json.Marshal(map[string]any{"success": true, "transaction": "0xabc"})
	h.Set("PAYMENT-RESPONSE", base64.StdEncoding.EncodeToString(settle))
	raw, err := DecodePaymentResponse(h)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil || m["transaction"] != "0xabc" {
		t.Fatalf("receipt: %v %v", m, err)
	}
}

var _ = errors.Is
