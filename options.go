package nebula

import (
	"fmt"
	"net/http"
	"time"
)

// AuthMode selects how the API key is sent.
type AuthMode string

const (
	// AuthModeAPIKey sends X-API-Key: <key> (default).
	AuthModeAPIKey AuthMode = "api_key"
	// AuthModeBearer sends Authorization: Bearer <key>. Never combined with
	// X-API-Key unless the caller explicitly configures both.
	AuthModeBearer AuthMode = "bearer"
)

// LocaleMode selects whether a configured locale is sent as the
// X-Nebula-Locale header (default) or the `locale` query parameter.
type LocaleMode string

const (
	// LocaleModeHeader sends X-Nebula-Locale.
	LocaleModeHeader LocaleMode = "header"
	// LocaleModeQuery sends ?locale=.
	LocaleModeQuery LocaleMode = "query"
)

// RetryPolicy configures automatic retries of idempotent GET requests.
// Retries happen only on transient transport errors, HTTP 429, and HTTP 503
// (and other 5xx marked transient by the server contract). POST requests and
// deterministic 4xx failures are never retried automatically.
type RetryPolicy struct {
	// MaxAttempts is the total number of tries including the first. 1 disables
	// retries; values < 1 are normalized to 1.
	MaxAttempts int
	// BaseDelay is the first backoff delay; each attempt doubles it up to
	// MaxDelay, with +/-25% jitter.
	BaseDelay time.Duration
	// MaxDelay caps the computed backoff.
	MaxDelay time.Duration
}

func (p RetryPolicy) normalized() RetryPolicy {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.BaseDelay <= 0 {
		p.BaseDelay = 250 * time.Millisecond
	}
	if p.MaxDelay <= 0 {
		p.MaxDelay = 5 * time.Second
	}
	return p
}

// Option configures a Client at New() time.
type Option func(*config) error

type config struct {
	baseURL      string
	httpClient   *http.Client
	timeout      time.Duration
	retry        RetryPolicy
	authMode     AuthMode
	locale       Locale
	localeMode   LocaleMode
	userAgent    string
	x402Mode     bool
	x402BaseURL  string
	logger       Logger
	perRequestUA string
}

func defaultConfig() *config {
	return &config{
		baseURL:     BaseURL,
		timeout:     30 * time.Second,
		retry:       RetryPolicy{MaxAttempts: 3, BaseDelay: 250 * time.Millisecond, MaxDelay: 5 * time.Second},
		authMode:    AuthModeAPIKey,
		localeMode:  LocaleModeHeader,
		userAgent:   "nebula-go/" + Version,
		x402BaseURL: X402BaseURL,
	}
}

// WithBaseURL overrides the account API base URL (testing / gateways).
// The URL must be absolute; it is joined with endpoint paths safely.
func WithBaseURL(u string) Option {
	return func(c *config) error {
		if u == "" {
			return fmt.Errorf("nebula: WithBaseURL: empty URL")
		}
		c.baseURL = u
		return nil
	}
}

// WithHTTPClient supplies the *http.Client used for requests. Pass an
// x402-wrapped client (see the x402 adapter package) or a custom-instrumented
// client here; nil keeps a default client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *config) error {
		if hc != nil {
			c.httpClient = hc
		}
		return nil
	}
}

// WithTimeout sets the overall request timeout. The deadline covers all retry
// attempts of one call. Zero or negative keeps the default.
func WithTimeout(d time.Duration) Option {
	return func(c *config) error {
		if d < 0 {
			return fmt.Errorf("nebula: WithTimeout: negative duration")
		}
		if d > 0 {
			c.timeout = d
		}
		return nil
	}
}

// WithRetry configures automatic retries of idempotent GET requests.
// maxAttempts is the total tries including the first (1 = no retries);
// base is the initial backoff, max the cap. The delay honors the server's
// Retry-After and RateLimit-Reset hints before the computed backoff.
func WithRetry(maxAttempts int, base, max time.Duration) Option {
	return func(c *config) error {
		if maxAttempts < 1 {
			return fmt.Errorf("nebula: WithRetry: maxAttempts must be >= 1")
		}
		if base < 0 || max < 0 {
			return fmt.Errorf("nebula: WithRetry: negative delay")
		}
		c.retry = RetryPolicy{MaxAttempts: maxAttempts, BaseDelay: base, MaxDelay: max}
		return nil
	}
}

// WithRetryPolicy installs a fully specified retry policy.
func WithRetryPolicy(p RetryPolicy) Option {
	return func(c *config) error {
		c.retry = p.normalized()
		return nil
	}
}

// WithAuthMode switches between X-API-Key (default) and
// Authorization: Bearer <key>. The two are never sent together.
func WithAuthMode(m AuthMode) Option {
	return func(c *config) error {
		switch m {
		case AuthModeAPIKey, AuthModeBearer:
			c.authMode = m
			return nil
		}
		return fmt.Errorf("nebula: WithAuthMode: unknown mode %q", m)
	}
}

// WithLocale sets the global generated-text locale (X-Nebula-Locale by
// default; see WithLocaleMode for the ?locale= variant).
func WithLocale(l Locale) Option {
	return func(c *config) error {
		c.locale = l
		return nil
	}
}

// WithLocaleMode selects header or query transport for the locale.
func WithLocaleMode(m LocaleMode) Option {
	return func(c *config) error {
		switch m {
		case LocaleModeHeader, LocaleModeQuery:
			c.localeMode = m
			return nil
		}
		return fmt.Errorf("nebula: WithLocaleMode: unknown mode %q", m)
	}
}

// WithUserAgent overrides the User-Agent header. The SDK still defaults to
// nebula-go/<version> when unset.
func WithUserAgent(ua string) Option {
	return func(c *config) error {
		if ua != "" {
			c.userAgent = ua
		}
		return nil
	}
}

// WithLogger installs a debug logger. Transport logs redact API keys, bearer
// tokens and payment headers.
func WithLogger(l Logger) Option {
	return func(c *config) error {
		c.logger = l
		return nil
	}
}

// Logger is the minimal debug sink. nil disables logging.
type Logger interface {
	Printf(format string, args ...any)
}

// LoggerFunc adapts a function to Logger.
type LoggerFunc func(format string, args ...any)

// Printf implements Logger.
func (f LoggerFunc) Printf(format string, args ...any) { f(format, args...) }
