package nebula

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
)

// ResponseMetadata carries the HTTP context of an API response: status,
// headers, and the parsed rate-limit / credit / payment figures Nebula
// reports on every call.
type ResponseMetadata struct {
	// StatusCode is the HTTP status of the response.
	StatusCode int
	// Headers is the full response header set.
	Headers http.Header
	// RequestID is the request id when the API sends one.
	RequestID string

	// RateLimitLimit is the requests-per-minute allowance on the key's plan.
	RateLimitLimit *int
	// RateLimitRemaining is the requests left in the current window.
	RateLimitRemaining *int
	// RateLimitReset is the seconds until the window resets.
	RateLimitReset *int
	// RetryAfterSeconds is the server hint on 429/503 responses.
	RetryAfterSeconds *int

	// CreditsCharged is the number of credits the request cost (0 on errors).
	CreditsCharged *int64
	// CreditsRemaining is the spendable balance after the call.
	CreditsRemaining *int64

	// PaymentReceipt is the raw PAYMENT-RESPONSE receipt header returned by
	// the x402 surface after settlement; nil on account-mode calls.
	PaymentReceipt []byte
}

// BaseResponse is embedded in every typed response so transport metadata and
// the raw decoded body stay reachable without weakening the field types.
type BaseResponse struct {
	Metadata ResponseMetadata `json:"-"`
	// Raw keeps the decoded JSON object for forward-compatible access to
	// fields added upstream after this SDK release.
	Raw map[string]any `json:"-"`
}

type metaCarrier interface{ setMetadata(ResponseMetadata) }

func (b *BaseResponse) setMetadata(m ResponseMetadata) { b.Metadata = m }

func (b *BaseResponse) setRaw(raw map[string]any) { b.Raw = raw }

func parseMetadata(status int, h http.Header) ResponseMetadata {
	m := ResponseMetadata{
		StatusCode:     status,
		Headers:        h.Clone(),
		RequestID:      firstNonEmpty(h.Get("X-Request-Id"), h.Get("X-Request-ID"), h.Get("Request-Id")),
		PaymentReceipt: paymentReceipt(h),
	}
	if v := h.Get("RateLimit-Limit"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			m.RateLimitLimit = &n
		}
	}
	if v := h.Get("RateLimit-Remaining"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			m.RateLimitRemaining = &n
		}
	}
	if v := h.Get("RateLimit-Reset"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			m.RateLimitReset = &n
		}
	}
	if v := h.Get("Retry-After"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			m.RetryAfterSeconds = &n
		}
	}
	if v := h.Get("X-Credits-Charged"); v != "" {
		if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			m.CreditsCharged = &n
		}
	}
	if v := h.Get("X-Credits-Remaining"); v != "" {
		if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			m.CreditsRemaining = &n
		}
	}
	return m
}

func paymentReceipt(h http.Header) []byte {
	for _, k := range []string{"PAYMENT-RESPONSE", "Payment-Response", "X-Payment-Response"} {
		if v := h.Get(k); v != "" {
			if decoded, err := base64.StdEncoding.DecodeString(v); err == nil {
				return decoded
			}
			return []byte(v)
		}
	}
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
