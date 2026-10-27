// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "time"

// AccountResponse models an object of the `account` service.
type AccountResponse struct {
	BaseResponse
	// The account's credit position this cycle; absent for keys that are not metered.
	Credits *AccountResponseCredits `json:"credits"`
	// How many days back the plan can read; null when unbounded.
	HistoryWindowDays *int `json:"history_window_days"`
	// The account's API plan, such as free or pro.
	Plan string `json:"plan"`
	// Requests allowed per minute; null when the key is not rate limited.
	RateLimitPerMinute *int `json:"rate_limit_per_minute"`
}

// AccountResponseCredits models an object of the `account` service.
type AccountResponseCredits struct {
	// UTC end of the current credit cycle, when plan credits renew.
	CycleEnd time.Time `json:"cycle_end"`
	// UTC start of the current credit cycle.
	CycleStart time.Time `json:"cycle_start"`
	// Credits the plan includes this cycle.
	PlanCredits int64 `json:"plan_credits"`
	// Plan credits left this cycle.
	PlanCreditsRemaining int64 `json:"plan_credits_remaining"`
	// Plan credits spent this cycle.
	PlanCreditsUsed int64 `json:"plan_credits_used"`
	// Unexpired purchased credits left, spent after plan credits.
	PurchasedCreditsRemaining int64 `json:"purchased_credits_remaining"`
	// Everything left to spend: plan plus purchased credits.
	TotalRemaining int64 `json:"total_remaining"`
}
