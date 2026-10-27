// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// List - Author clusters. GET /clusters
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters/
func (s *ClustersService) List(ctx context.Context, req *ClustersRequest) (*ClustersResponse, error) {
	return get[ClustersResponse](ctx, s.client, "/clusters", req, "list_clusters")
}

// Comparison - Cluster comparison. GET /clusters/comparison
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison/
func (s *ClustersService) Comparison(ctx context.Context, req *ClustersComparisonRequest) (*ClustersComparisonResponse, error) {
	return get[ClustersComparisonResponse](ctx, s.client, "/clusters/comparison", req, "compare_clusters")
}

// DivisiveAssets - Divisive assets. GET /clusters/comparison/divisive-assets
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-divisive-assets/
func (s *ClustersService) DivisiveAssets(ctx context.Context, req *DivisiveAssetsRequest) (*DivisiveAssetsResponse, error) {
	return get[DivisiveAssetsResponse](ctx, s.client, "/clusters/comparison/divisive-assets", req, "get_divisive_assets")
}

// ComparisonHypeCycles - Hype cycle arrivals. GET /clusters/comparison/hype-cycles
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-hype-cycles/
func (s *ClustersService) ComparisonHypeCycles(ctx context.Context, req *ComparisonHypeCyclesRequest) (*ComparisonHypeCyclesResponse, error) {
	return get[ComparisonHypeCyclesResponse](ctx, s.client, "/clusters/comparison/hype-cycles", req, "compare_hype_cycles")
}

// ComparisonTraits - Cluster trait comparison. GET /clusters/comparison/traits
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-traits/
func (s *ClustersService) ComparisonTraits(ctx context.Context, req *ComparisonTraitsRequest) (*ComparisonTraitsResponse, error) {
	return get[ComparisonTraitsResponse](ctx, s.client, "/clusters/comparison/traits", req, "compare_cluster_traits")
}

// HypeCycles - Cluster hype cycles. GET /clusters/hype-cycles
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-hype-cycles/
func (s *ClustersService) HypeCycles(ctx context.Context, req *ClusterHypeCyclesRequest) (*ClusterHypeCyclesResponse, error) {
	return get[ClusterHypeCyclesResponse](ctx, s.client, "/clusters/hype-cycles", req, "get_cluster_hype_cycles")
}

// Overlap - Cluster overlap. GET /clusters/overlap
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-overlap/
func (s *ClustersService) Overlap(ctx context.Context, req *ClusterOverlapRequest) (*ClusterOverlapResponse, error) {
	return get[ClusterOverlapResponse](ctx, s.client, "/clusters/overlap", req, "get_cluster_overlap")
}

// Profile - Cluster profile. GET /clusters/profile
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-profile/
func (s *ClustersService) Profile(ctx context.Context, req *ClusterProfileRequest) (*ClusterProfileResponse, error) {
	return get[ClusterProfileResponse](ctx, s.client, "/clusters/profile", req, "get_cluster_profile")
}

// Rotation - Cluster rotation. GET /clusters/rotation
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-rotation/
func (s *ClustersService) Rotation(ctx context.Context, req *ClusterRotationRequest) (*ClusterRotationResponse, error) {
	return get[ClusterRotationResponse](ctx, s.client, "/clusters/rotation", req, "get_cluster_rotation")
}

// Stances - Cluster stances. GET /clusters/stances
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-stances/
func (s *ClustersService) Stances(ctx context.Context, req *ClusterStancesRequest) (*ClusterStancesResponse, error) {
	return get[ClusterStancesResponse](ctx, s.client, "/clusters/stances", req, "get_cluster_stances")
}
