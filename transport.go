package nebula

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"
)

// queryRequest is implemented by every generated request type.
type queryRequest interface {
	query() (url.Values, error)
}

type transport struct {
	cfg    *config
	apiKey string
}

func newTransport(cfg *config, apiKey string) *transport {
	return &transport{cfg: cfg, apiKey: apiKey}
}

// get executes a GET for operation op and decodes the body into T.
func get[T any](ctx context.Context, c *Client, path string, req queryRequest, op string) (*T, error) {
	var q url.Values
	var err error
	if req != nil {
		q, err = req.query()
		if err != nil {
			return nil, err
		}
	}
	if c.cfg.x402Mode && !x402Eligible(op) {
		return nil, fmt.Errorf("nebula: %s: %w", op, ErrNotX402Eligible)
	}
	if err := validateQuery(op, q); err != nil {
		return nil, err
	}
	return doJSON[T](ctx, c.tr, http.MethodGet, path, q, nil, op)
}

// postJSON executes a non-streaming POST with a JSON body.
func postJSON[T any](ctx context.Context, c *Client, path string, body any, op string) (*T, error) {
	if c.cfg.x402Mode {
		return nil, fmt.Errorf("nebula: %s: %w", op, ErrNotX402Eligible)
	}
	return doJSON[T](ctx, c.tr, http.MethodPost, path, nil, body, op)
}

// doJSON performs the request with retry, decodes the JSON body into a new T,
// attaches metadata, and keeps the raw decoded object.
func doJSON[T any](ctx context.Context, t *transport, method, path string, q url.Values, body any, op string) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, t.cfg.timeout)
	defer cancel()
	var lastErr error
	attempts := t.retryableAttempts(method)
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, &TransportError{Op: op, Err: err}
		}
		if attempt > 0 {
			delay := t.backoff(attempt, lastErr)
			if sleepErr := sleepCtx(ctx, delay); sleepErr != nil {
				return nil, &TransportError{Op: op, Err: sleepErr}
			}
		}
		resp, err := t.doOnce(ctx, method, path, q, body, op)
		if err == nil {
			out, derr := decodeResponse[T](resp, op)
			resp.Body.Close()
			if derr != nil {
				return nil, derr
			}
			return out, nil
		}
		lastErr = err
		if !t.shouldRetry(method, err) {
			return nil, err
		}
	}
	return nil, lastErr
}

func (t *transport) retryableAttempts(method string) int {
	if method != http.MethodGet || t.cfg.x402Mode {
		return 1
	}
	return t.cfg.retry.normalized().MaxAttempts
}

// shouldRetry reports whether err is a transient failure worth retrying.
// Only safe idempotent GET retries: transport errors, 429, and 5xx.
// Deterministic 4xx (400, 401, 402, 403, 404, 405) and POSTs never retry.
func (t *transport) shouldRetry(method string, err error) bool {
	if method != http.MethodGet || t.cfg.x402Mode || errors.Is(err, ErrPaymentUncertain) {
		return false
	}
	var te *TransportError
	if errors.As(err, &te) {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return false
		}
		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() {
			return true
		}
		var dns *net.DNSError
		if errors.As(err, &dns) && dns.IsTemporary {
			return true
		}
		for _, transient := range []error{io.EOF, io.ErrUnexpectedEOF, syscall.ECONNRESET, syscall.ECONNREFUSED, syscall.ECONNABORTED, syscall.EPIPE} {
			if errors.Is(err, transient) {
				return true
			}
		}
		return false
	}
	var ae *APIError
	if errors.As(err, &ae) {
		switch ae.StatusCode {
		case 429, 500, 502, 503, 504:
			return true
		}
		return false
	}
	return false
}

// backoff returns the delay before the next attempt: Retry-After wins, then
// RateLimit-Reset, then exponential backoff with +/-25% jitter.
func (t *transport) backoff(attempt int, lastErr error) time.Duration {
	var ae *APIError
	if errors.As(lastErr, &ae) {
		if ae.Metadata.RetryAfterSeconds != nil {
			return secondsDelay(*ae.Metadata.RetryAfterSeconds)
		}
		if ae.RetryAfter > 0 {
			return ae.RetryAfter
		}
		if ae.Metadata.RateLimitReset != nil {
			return secondsDelay(*ae.Metadata.RateLimitReset)
		}
	}
	p := t.cfg.retry.normalized()
	d := min(p.BaseDelay, p.MaxDelay)
	for i := 1; i < attempt && d < p.MaxDelay; i++ {
		if d > p.MaxDelay/2 {
			d = p.MaxDelay
		} else {
			d *= 2
		}
	}
	// The package-level generator is safe for concurrent use. Cap after jitter.
	jitter := time.Duration(float64(d) * 0.25 * (rand.Float64()*2 - 1))
	if jitter > 0 && d > p.MaxDelay-jitter {
		return p.MaxDelay
	}
	return d + jitter
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// doOnce performs a single HTTP request and returns either the response
// (caller closes Body) or a classified error.
func (t *transport) doOnce(ctx context.Context, method, path string, q url.Values, body any, op string) (*http.Response, error) {
	u, err := t.buildURL(ctx, path, q)
	if err != nil {
		return nil, err
	}
	var rdr io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, &ValidationError{Op: op, Field: "body", Msg: "cannot serialize request body"}
		}
		rdr = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", t.cfg.userAgent)
	t.applyAuth(req)
	locale, localeMode := localeFor(ctx, t.cfg)
	if locale != "" && localeMode == LocaleModeHeader {
		req.Header.Set("X-Nebula-Locale", string(locale))
	}
	if t.cfg.logger != nil {
		t.cfg.logger.Printf("nebula: %s %s", method, redactURL(u))
	}
	resp, err := t.cfg.httpClient.Do(req)
	if err != nil {
		return nil, &TransportError{Op: op, Err: err}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, decodeAPIError(resp, op)
	}
	return resp, nil
}

func (t *transport) applyAuth(req *http.Request) {
	if t.apiKey == "" {
		return
	}
	switch t.cfg.authMode {
	case AuthModeBearer:
		req.Header.Set("Authorization", "Bearer "+t.apiKey)
	default:
		req.Header.Set("X-API-Key", t.apiKey)
	}
}

// buildURL joins base + path + query without producing double slashes.
func (t *transport) buildURL(ctx context.Context, path string, q url.Values) (string, error) {
	base, err := url.Parse(t.cfg.baseURL)
	if err != nil {
		return "", err
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/" + strings.TrimLeft(path, "/")
	vals := url.Values{}
	for k, vs := range q {
		for _, v := range vs {
			vals.Add(k, v)
		}
	}
	if locale, mode := localeFor(ctx, t.cfg); locale != "" && mode == LocaleModeQuery {
		vals.Set("locale", string(locale))
	}
	base.RawQuery = vals.Encode()
	return base.String(), nil
}

// decodeResponse unmarshals the JSON body into T, attaches metadata and keeps
// the decoded raw object reachable for forward-compatible field access.
func decodeResponse[T any](resp *http.Response, op string) (*T, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &TransportError{Op: op, Err: err}
	}
	out := new(T)
	if trimmed := bytes.TrimSpace(body); len(trimmed) == 0 || trimmed[0] != '{' {
		return nil, &DecodeError{Op: op, Err: errors.New("expected a JSON object")}
	}
	if err := json.Unmarshal(body, out); err != nil {
		return nil, &DecodeError{Op: op, Err: fmt.Errorf("decode response: %w", err)}
	}
	if c, ok := any(out).(metaCarrier); ok {
		c.setMetadata(parseMetadata(resp.StatusCode, resp.Header))
	}
	var raw map[string]any
	if p, ok := any(out).(interface{ setRaw(map[string]any) }); ok && json.Unmarshal(body, &raw) == nil {
		p.setRaw(raw)
	}
	return out, nil
}

func decodeAPIError(resp *http.Response, op string) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	ae := &APIError{
		Op:         op,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Metadata:   parseMetadata(resp.StatusCode, resp.Header),
		Body:       body,
	}
	var payload struct {
		Message          *string `json:"message"`
		PlanRenewsAt     *string `json:"plan_renews_at"`
		PurchasedCredits *int64  `json:"purchased_credits"`
	}
	if len(body) > 0 && json.Unmarshal(body, &payload) == nil {
		if payload.Message != nil {
			ae.Message = *payload.Message
		}
		if payload.PlanRenewsAt != nil {
			ae.PlanRenewsAt = *payload.PlanRenewsAt
		}
		if payload.PurchasedCredits != nil {
			ae.PurchasedCredits = payload.PurchasedCredits
		}
	}
	if ae.Message == "" {
		ae.Message = http.StatusText(resp.StatusCode)
	}
	if ae.Metadata.RetryAfterSeconds != nil {
		ae.RetryAfter = secondsDelay(*ae.Metadata.RetryAfterSeconds)
	}
	return ae
}

// redactURL removes credential-bearing query params from logged URLs.
func redactURL(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return "<unparseable url>"
	}
	q := parsed.Query()
	parsed.User = nil
	for k := range q {
		switch strings.ToLower(k) {
		case "api_key", "key", "token", "payment-signature", "signature", "authorization", "x-api-key", "x-payment":
			q.Set(k, "[redacted]")
		}
	}
	parsed.RawQuery = q.Encode()
	return parsed.String()
}
