// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
	"time"
)

// CalendarEventsRequest carries the query parameters of `list_calendar_events` (GET /calendar-events).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-calendar-events/
type CalendarEventsRequest struct {
	Asset              *AssetRef
	Subject            *SubjectRef
	Limit              *int
	Offset             *int
	StartTimeAfter     *time.Time
	StartTimeBefore    *time.Time
	Domains            []CalendarDomain
	EventTypes         []EventType
	Statuses           []EventStatus
	MinImportanceScore *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *CalendarEventsRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if r.Asset != nil {
		v.Set("asset_id", r.Asset.ID)
		v.Set("asset_class", string(r.Asset.Class))
	}
	if r.Subject != nil {
		v.Set("subject_id", r.Subject.ID)
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	if r.Offset != nil {
		v.Set("offset", strconv.Itoa(*r.Offset))
	}
	if r.StartTimeAfter != nil {
		v.Set("start_time_after", r.StartTimeAfter.UTC().Format(time.RFC3339))
	}
	if r.StartTimeBefore != nil {
		v.Set("start_time_before", r.StartTimeBefore.UTC().Format(time.RFC3339))
	}
	if len(r.Domains) > 0 {
		v.Set("domains", csvJoin(r.Domains))
	}
	if len(r.EventTypes) > 0 {
		v.Set("event_types", csvJoin(r.EventTypes))
	}
	if len(r.Statuses) > 0 {
		v.Set("statuses", csvJoin(r.Statuses))
	}
	if r.MinImportanceScore != nil {
		v.Set("min_importance_score", strconv.Itoa(*r.MinImportanceScore))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `list_calendar_events`
// before any HTTP call is made.
func (r *CalendarEventsRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "list_calendar_events", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("list_calendar_events"); err != nil {
		return err
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 200) {
		return &ValidationError{Op: "list_calendar_events", Field: "limit", Msg: "out of range [1, 200]"}
	}
	if r.Offset != nil && (*r.Offset < 0 || *r.Offset > 100000) {
		return &ValidationError{Op: "list_calendar_events", Field: "offset", Msg: "out of range [0, 100000]"}
	}
	if r.MinImportanceScore != nil && (*r.MinImportanceScore < 1 || *r.MinImportanceScore > 3) {
		return &ValidationError{Op: "list_calendar_events", Field: "min_importance_score", Msg: "out of range [1, 3]"}
	}
	return nil
}

// InsightsRequest carries the query parameters of `list_insights` (GET /insights).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-insights/
type InsightsRequest struct {
	Asset           *AssetRef
	Subject         *SubjectRef
	Limit           *int
	Offset          *int
	StartTimeAfter  *time.Time
	StartTimeBefore *time.Time
	Tags            []string
	ExcludeTags     []string
	Categories      []string
	Domains         []InsightDomain
	OrderBy         *InsightOrderBy
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *InsightsRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if r.Asset != nil {
		v.Set("asset_id", r.Asset.ID)
		v.Set("asset_class", string(r.Asset.Class))
	}
	if r.Subject != nil {
		v.Set("subject_id", r.Subject.ID)
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	if r.Offset != nil {
		v.Set("offset", strconv.Itoa(*r.Offset))
	}
	if r.StartTimeAfter != nil {
		v.Set("start_time_after", r.StartTimeAfter.UTC().Format(time.RFC3339))
	}
	if r.StartTimeBefore != nil {
		v.Set("start_time_before", r.StartTimeBefore.UTC().Format(time.RFC3339))
	}
	if len(r.Tags) > 0 {
		v.Set("tags", csvJoin(r.Tags))
	}
	if len(r.ExcludeTags) > 0 {
		v.Set("exclude_tags", csvJoin(r.ExcludeTags))
	}
	if len(r.Categories) > 0 {
		v.Set("categories", csvJoin(r.Categories))
	}
	if len(r.Domains) > 0 {
		v.Set("domains", csvJoin(r.Domains))
	}
	if r.OrderBy != nil {
		v.Set("order_by", string(*r.OrderBy))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `list_insights`
// before any HTTP call is made.
func (r *InsightsRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "list_insights", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("list_insights"); err != nil {
		return err
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 50) {
		return &ValidationError{Op: "list_insights", Field: "limit", Msg: "out of range [1, 50]"}
	}
	if r.Offset != nil && (*r.Offset < 0 || *r.Offset > 100000) {
		return &ValidationError{Op: "list_insights", Field: "offset", Msg: "out of range [0, 100000]"}
	}
	return nil
}

// InsightAttentionRequest carries the query parameters of `get_insight_attention` (GET /insights/attention).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-insight-attention/
type InsightAttentionRequest struct {
	InsightID int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *InsightAttentionRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("insight_id", strconv.Itoa(r.InsightID))
	return v, nil
}

// validate enforces the documented client-side rules of `get_insight_attention`
// before any HTTP call is made.
func (r *InsightAttentionRequest) validate() error {
	if r.InsightID < 1 {
		return &ValidationError{Op: "get_insight_attention", Field: "insight_id", Msg: "out of range [1, +inf]"}
	}
	return nil
}
