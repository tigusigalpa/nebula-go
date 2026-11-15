package nebula

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Version is the SDK release.
const Version = "0.2.0"

// BaseURL is the account-authenticated API base.
const BaseURL = "https://nebula-api.hiddensystems.ai/api/v1/public"

// X402BaseURL is the pay-per-request x402 API base. Only the 37
// x-x402-eligible GET operations are served there.
const X402BaseURL = "https://nebula-api.hiddensystems.ai/api/v1/x402"

// Client is the Nebula API client. Create it with New; it is immutable after
// construction and safe for concurrent use across goroutines.
type Client struct {
	cfg *config
	tr  *transport

	Discovery       *DiscoveryService
	Lookup          *LookupService
	Social          *SocialService
	Forecasts       *ForecastsService
	MarketSentiment *MarketSentimentService
	MarketData      *MarketDataService
	Intelligence    *IntelligenceService
	Clusters        *ClustersService
	Agent           *AgentService
	Account         *AccountService
	Predictions     *PredictionsService
}

// New creates an account-mode Client authenticated by apiKey.
// The key must come from a secure source (env var, secret store); it is sent
// as X-API-Key unless WithAuthMode selects the Bearer form.
func New(apiKey string, opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	if apiKey == "" && !cfg.x402Mode {
		return nil, fmt.Errorf("nebula: New: empty API key (pass a key or use the x402 adapter)")
	}
	if cfg.httpClient == nil {
		cfg.httpClient = &http.Client{}
	}
	cfg.httpClient = secureHTTPClient(cfg.httpClient)
	base, err := normalizeBase(cfg.baseURL)
	if err != nil {
		return nil, fmt.Errorf("nebula: invalid base URL: %w", err)
	}
	cfg.baseURL = base
	c := &Client{cfg: cfg}
	c.tr = newTransport(cfg, apiKey)
	svc := &service{client: c}
	c.Discovery = &DiscoveryService{service: svc}
	c.Lookup = &LookupService{service: svc}
	c.Social = &SocialService{service: svc}
	c.Forecasts = &ForecastsService{service: svc}
	c.MarketSentiment = &MarketSentimentService{service: svc}
	c.MarketData = &MarketDataService{service: svc}
	c.Intelligence = &IntelligenceService{service: svc}
	c.Clusters = &ClustersService{service: svc}
	c.Agent = &AgentService{service: svc}
	c.Account = &AccountService{service: svc}
	c.Predictions = &PredictionsService{service: svc}
	return c, nil
}

// NewX402Client builds a client pinned to the x402 surface
// (/api/v1/x402) without an API key. It is used by the x402 adapter package
// with a payment-wrapped *http.Client; non-eligible operations
// (/account, forecasts, chat completions) fail locally before any HTTP call.
func NewX402Client(httpClient *http.Client, opts ...Option) (*Client, error) {
	opts = append([]Option{WithBaseURL(X402BaseURL)}, opts...)
	cfg := defaultConfig()
	cfg.x402Mode = true
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	if httpClient != nil {
		cfg.httpClient = httpClient
	}
	if cfg.httpClient == nil {
		cfg.httpClient = &http.Client{}
	}
	cfg.httpClient = secureHTTPClient(cfg.httpClient)
	base, err := normalizeBase(cfg.baseURL)
	if err != nil {
		return nil, err
	}
	cfg.baseURL = base
	c := &Client{cfg: cfg}
	c.tr = newTransport(cfg, "")
	svc := &service{client: c}
	c.Discovery = &DiscoveryService{service: svc}
	c.Lookup = &LookupService{service: svc}
	c.Social = &SocialService{service: svc}
	c.Forecasts = &ForecastsService{service: svc}
	c.MarketSentiment = &MarketSentimentService{service: svc}
	c.MarketData = &MarketDataService{service: svc}
	c.Intelligence = &IntelligenceService{service: svc}
	c.Clusters = &ClustersService{service: svc}
	c.Agent = &AgentService{service: svc}
	c.Account = &AccountService{service: svc}
	c.Predictions = &PredictionsService{service: svc}
	return c, nil
}

// X402 reports whether the client is bound to the x402 surface.
func (c *Client) X402() bool { return c.cfg.x402Mode }

// normalizeBase parses u and returns a clean absolute base URL without a
// trailing slash, so joins never produce double slashes.
func normalizeBase(u string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(u))
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return "", fmt.Errorf("base URL must be http(s): %q", u)
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("base URL has no host: %q", u)
	}
	if parsed.User != nil {
		return "", fmt.Errorf("base URL must not contain credentials")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	parsed.RawPath = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

// Copy the client so later caller edits to its settings do not change the SDK.
// Cross-origin redirects could leak X-API-Key or a signed payment header.
func secureHTTPClient(hc *http.Client) *http.Client {
	copy := *hc
	checkRedirect := hc.CheckRedirect
	copy.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 && (req.URL.Scheme != via[0].URL.Scheme || !strings.EqualFold(req.URL.Host, via[0].URL.Host)) {
			return http.ErrUseLastResponse
		}
		if checkRedirect != nil {
			return checkRedirect(req, via)
		}
		if len(via) >= 10 {
			return fmt.Errorf("stopped after 10 redirects")
		}
		return nil
	}
	return &copy
}

// service is the shared plumbing embedded in every service group.
type service struct {
	client *Client
}

// Service group structs.
type DiscoveryService struct{ *service }
type LookupService struct{ *service }
type SocialService struct{ *service }
type ForecastsService struct{ *service }
type MarketSentimentService struct{ *service }
type MarketDataService struct{ *service }
type IntelligenceService struct{ *service }
type ClustersService struct{ *service }
type AgentService struct{ *service }
type AccountService struct{ *service }
type PredictionsService struct{ *service }
