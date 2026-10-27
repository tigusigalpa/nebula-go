// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// CalendarEvents - Calendar events. GET /calendar-events
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-calendar-events/
func (s *IntelligenceService) CalendarEvents(ctx context.Context, req *CalendarEventsRequest) (*CalendarEventsResponse, error) {
	return get[CalendarEventsResponse](ctx, s.client, "/calendar-events", req, "list_calendar_events")
}

// Insights - Intelligence. GET /insights
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-insights/
func (s *IntelligenceService) Insights(ctx context.Context, req *InsightsRequest) (*InsightsResponse, error) {
	return get[InsightsResponse](ctx, s.client, "/insights", req, "list_insights")
}

// InsightAttention - Insight attention. GET /insights/attention
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-insight-attention/
func (s *IntelligenceService) InsightAttention(ctx context.Context, req *InsightAttentionRequest) (*InsightAttentionResponse, error) {
	return get[InsightAttentionResponse](ctx, s.client, "/insights/attention", req, "get_insight_attention")
}
