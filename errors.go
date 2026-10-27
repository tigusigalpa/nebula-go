package nebula

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Sentinel errors classify failures for errors.Is. Every API error is also an
// *APIError for errors.As inspection.
var (
	// ErrValidation marks client-side request validation failures raised
	// before any HTTP call (also matches *ValidationError).
	ErrValidation = errors.New("nebula: validation failed")
	// ErrTransport marks network/transport failures without an HTTP status.
	ErrTransport = errors.New("nebula: transport error")
	// ErrBadRequest maps HTTP 400.
	ErrBadRequest = errors.New("nebula: bad request")
	// ErrUnauthorized maps HTTP 401 (missing, invalid, or revoked API key).
	ErrUnauthorized = errors.New("nebula: unauthorized")
	// ErrOutOfCredits maps HTTP 402 (account credits used up).
	ErrOutOfCredits = errors.New("nebula: out of credits")
	// ErrForbidden maps HTTP 403 (history beyond the plan's window).
	ErrForbidden = errors.New("nebula: forbidden")
	// ErrNotFound maps HTTP 404.
	ErrNotFound = errors.New("nebula: not found")
	// ErrMethodNotAllowed maps HTTP 405.
	ErrMethodNotAllowed = errors.New("nebula: method not allowed")
	// ErrRateLimited maps HTTP 429.
	ErrRateLimited = errors.New("nebula: rate limited")
	// ErrServer maps HTTP 500.
	ErrServer = errors.New("nebula: internal server error")
	// ErrServiceUnavailable maps HTTP 503.
	ErrServiceUnavailable = errors.New("nebula: service unavailable")
	// ErrPaymentRequired maps x402 HTTP 402 payment offers and settlement
	// failures on the x402 surface.
	ErrPaymentRequired = errors.New("nebula: x402 payment required")
	// ErrPaymentUncertain marks a request whose payment may have settled but
	// whose response was lost (503 after signing, dropped connection).
	// Reconcile the original payment before authorizing another.
	ErrPaymentUncertain = errors.New("nebula: x402 payment state uncertain")
	// ErrNotX402Eligible marks calls to endpoints the x402 surface does not
	// serve (/account, both forecasts, /chat/completions).
	ErrNotX402Eligible = errors.New("nebula: operation not eligible for x402")
)

const maxErrorBodyBytes = 8192

// ValidationError is raised before the network call when a request breaks a
// documented client-side rule (required pair, range, enum, exclusivity).
type ValidationError struct {
	Op    string
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("nebula: %s: invalid %s: %s", e.Op, e.Field, e.Msg)
}

// Is supports errors.Is(err, ErrValidation).
func (e *ValidationError) Is(target error) bool { return target == ErrValidation }

// APIError is an HTTP error response from the API. It keeps the status, the
// API message, headers, parsed metadata, and a bounded copy of the raw body.
type APIError struct {
	// Op is the operation id, e.g. "list_assets".
	Op string
	// StatusCode is the HTTP status.
	StatusCode int
	// Message is the API's human-readable error message.
	Message string
	// Headers of the error response.
	Headers http.Header
	// Metadata parses rate-limit/credit/retry headers.
	Metadata ResponseMetadata
	// Body is the raw response body truncated to maxErrorBodyBytes.
	Body []byte
	// PlanRenewsAt is the UTC credit-renewal time reported on 402 bodies.
	PlanRenewsAt string
	// PurchasedCredits is the purchased-credit balance reported on 402 bodies.
	PurchasedCredits *int64
	// RetryAfter is the parsed Retry-After delay when present.
	RetryAfter time.Duration
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("nebula: %s: HTTP %d: %s", e.Op, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("nebula: %s: HTTP %d", e.Op, e.StatusCode)
}

// Is maps the status code onto the sentinel errors.
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrBadRequest:
		return e.StatusCode == 400
	case ErrUnauthorized:
		return e.StatusCode == 401
	case ErrOutOfCredits:
		return e.StatusCode == 402 && !e.PaymentOffer()
	case ErrForbidden:
		return e.StatusCode == 403
	case ErrNotFound:
		return e.StatusCode == 404
	case ErrMethodNotAllowed:
		return e.StatusCode == 405
	case ErrRateLimited:
		return e.StatusCode == 429
	case ErrServer:
		return e.StatusCode == 500
	case ErrServiceUnavailable:
		return e.StatusCode == 503
	case ErrPaymentRequired:
		return e.StatusCode == 402 && e.PaymentOffer()
	}
	return false
}

// PaymentOffer reports whether the 402 carries a PAYMENT-REQUIRED x402 offer
// rather than an account out-of-credits error.
func (e *APIError) PaymentOffer() bool { return e.Headers.Get("PAYMENT-REQUIRED") != "" }

// TransportError wraps network-level failures (DNS, TLS, timeouts, refused
// connections) that produced no HTTP response.
type TransportError struct {
	Op  string
	Err error
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("nebula: %s: transport: %v", e.Op, e.Err)
}

func (e *TransportError) Unwrap() error { return e.Err }

func (e *TransportError) Is(target error) bool { return target == ErrTransport }

// StreamError terminates a chat stream abnormally (malformed event, transport
// drop, or an API error frame).
type StreamError struct {
	Op  string
	Err error
}

func (e *StreamError) Error() string { return fmt.Sprintf("nebula: %s: stream: %v", e.Op, e.Err) }
func (e *StreamError) Unwrap() error { return e.Err }
