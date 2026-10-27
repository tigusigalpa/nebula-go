// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// ListAssets - Text to asset ID. GET /assets
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-assets/
func (s *LookupService) ListAssets(ctx context.Context, req *ListAssetsRequest) (*ListAssetsResponse, error) {
	return get[ListAssetsResponse](ctx, s.client, "/assets", req, "list_assets")
}

// Search - Search. GET /search
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-search/
func (s *LookupService) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	return get[SearchResponse](ctx, s.client, "/search", req, "search")
}

// ListSubjects - Text to subject ID. GET /subjects
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-subjects/
func (s *LookupService) ListSubjects(ctx context.Context, req *ListSubjectsRequest) (*ListSubjectsResponse, error) {
	return get[ListSubjectsResponse](ctx, s.client, "/subjects", req, "list_subjects")
}
