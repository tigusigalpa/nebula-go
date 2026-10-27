// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
)

// ClustersRequest carries the query parameters of `list_clusters` (GET /clusters).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters/
type ClustersRequest struct {
	AssetClass *AssetClass
	Hours      *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ClustersRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if r.AssetClass != nil {
		v.Set("asset_class", string(*r.AssetClass))
	}
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `list_clusters`
// before any HTTP call is made.
func (r *ClustersRequest) validate() error {
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 720) {
		return &ValidationError{Op: "list_clusters", Field: "hours", Msg: "out of range [1, 720]"}
	}
	return nil
}

// ClustersComparisonRequest carries the query parameters of `compare_clusters` (GET /clusters/comparison).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison/
type ClustersComparisonRequest struct {
	AssetClasses []AssetClass
	Days         *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ClustersComparisonRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Days != nil {
		v.Set("days", strconv.Itoa(*r.Days))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `compare_clusters`
// before any HTTP call is made.
func (r *ClustersComparisonRequest) validate() error {
	if r.Days != nil && !intIn(*r.Days, []int{7, 30, 90}) {
		return &ValidationError{Op: "compare_clusters", Field: "days", Msg: "must be one of 7, 30, 90"}
	}
	return nil
}

// DivisiveAssetsRequest carries the query parameters of `get_divisive_assets` (GET /clusters/comparison/divisive-assets).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-divisive-assets/
type DivisiveAssetsRequest struct {
	AssetClasses []AssetClass
	Days         *int
	Limit        *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *DivisiveAssetsRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Days != nil {
		v.Set("days", strconv.Itoa(*r.Days))
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_divisive_assets`
// before any HTTP call is made.
func (r *DivisiveAssetsRequest) validate() error {
	if r.Days != nil && !intIn(*r.Days, []int{7, 30, 90}) {
		return &ValidationError{Op: "get_divisive_assets", Field: "days", Msg: "must be one of 7, 30, 90"}
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 30) {
		return &ValidationError{Op: "get_divisive_assets", Field: "limit", Msg: "out of range [1, 30]"}
	}
	return nil
}

// ComparisonHypeCyclesRequest carries the query parameters of `compare_hype_cycles` (GET /clusters/comparison/hype-cycles).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-hype-cycles/
type ComparisonHypeCyclesRequest struct {
	AssetClasses []AssetClass
	Days         *int
	Limit        *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ComparisonHypeCyclesRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Days != nil {
		v.Set("days", strconv.Itoa(*r.Days))
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `compare_hype_cycles`
// before any HTTP call is made.
func (r *ComparisonHypeCyclesRequest) validate() error {
	if r.Days != nil && !intIn(*r.Days, []int{7, 30, 90}) {
		return &ValidationError{Op: "compare_hype_cycles", Field: "days", Msg: "must be one of 7, 30, 90"}
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "compare_hype_cycles", Field: "limit", Msg: "out of range [1, 100]"}
	}
	return nil
}

// ComparisonTraitsRequest carries the query parameters of `compare_cluster_traits` (GET /clusters/comparison/traits).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-traits/
type ComparisonTraitsRequest struct {
	AssetClasses []AssetClass
	Days         *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ComparisonTraitsRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Days != nil {
		v.Set("days", strconv.Itoa(*r.Days))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `compare_cluster_traits`
// before any HTTP call is made.
func (r *ComparisonTraitsRequest) validate() error {
	if r.Days != nil && !intIn(*r.Days, []int{30, 90}) {
		return &ValidationError{Op: "compare_cluster_traits", Field: "days", Msg: "must be one of 30, 90"}
	}
	return nil
}

// ClusterHypeCyclesRequest carries the query parameters of `get_cluster_hype_cycles` (GET /clusters/hype-cycles).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-hype-cycles/
type ClusterHypeCyclesRequest struct {
	ClusterID    string
	AssetClasses []AssetClass
	Days         *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ClusterHypeCyclesRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("cluster_id", r.ClusterID)
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Days != nil {
		v.Set("days", strconv.Itoa(*r.Days))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_cluster_hype_cycles`
// before any HTTP call is made.
func (r *ClusterHypeCyclesRequest) validate() error {
	if r.ClusterID == "" {
		return &ValidationError{Op: "get_cluster_hype_cycles", Field: "cluster_id", Msg: "required"}
	}
	if r.Days != nil && !intIn(*r.Days, []int{30, 90}) {
		return &ValidationError{Op: "get_cluster_hype_cycles", Field: "days", Msg: "must be one of 30, 90"}
	}
	if isCommunity(r.ClusterID) && len(r.AssetClasses) > 0 {
		for _, c := range r.AssetClasses {
			if c != AssetClassEquity {
				return &ValidationError{Op: "get_cluster_hype_cycles", Field: "asset_class", Msg: "posting communities (wallstreetbets, reddit) apply only to equity"}
			}
		}
	}
	return nil
}

// ClusterOverlapRequest carries the query parameters of `get_cluster_overlap` (GET /clusters/overlap).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-overlap/
type ClusterOverlapRequest struct {
	ClusterID    string
	AssetClasses []AssetClass
	Days         *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ClusterOverlapRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("cluster_id", r.ClusterID)
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Days != nil {
		v.Set("days", strconv.Itoa(*r.Days))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_cluster_overlap`
// before any HTTP call is made.
func (r *ClusterOverlapRequest) validate() error {
	if r.ClusterID == "" {
		return &ValidationError{Op: "get_cluster_overlap", Field: "cluster_id", Msg: "required"}
	}
	if r.Days != nil && !intIn(*r.Days, []int{7, 30, 90}) {
		return &ValidationError{Op: "get_cluster_overlap", Field: "days", Msg: "must be one of 7, 30, 90"}
	}
	if isCommunity(r.ClusterID) && len(r.AssetClasses) > 0 {
		for _, c := range r.AssetClasses {
			if c != AssetClassEquity {
				return &ValidationError{Op: "get_cluster_overlap", Field: "asset_class", Msg: "posting communities (wallstreetbets, reddit) apply only to equity"}
			}
		}
	}
	return nil
}

// ClusterProfileRequest carries the query parameters of `get_cluster_profile` (GET /clusters/profile).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-profile/
type ClusterProfileRequest struct {
	ClusterID    string
	AssetClasses []AssetClass
	Days         *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ClusterProfileRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("cluster_id", r.ClusterID)
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Days != nil {
		v.Set("days", strconv.Itoa(*r.Days))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_cluster_profile`
// before any HTTP call is made.
func (r *ClusterProfileRequest) validate() error {
	if r.ClusterID == "" {
		return &ValidationError{Op: "get_cluster_profile", Field: "cluster_id", Msg: "required"}
	}
	if r.Days != nil && !intIn(*r.Days, []int{30, 90}) {
		return &ValidationError{Op: "get_cluster_profile", Field: "days", Msg: "must be one of 30, 90"}
	}
	if isCommunity(r.ClusterID) && len(r.AssetClasses) > 0 {
		for _, c := range r.AssetClasses {
			if c != AssetClassEquity {
				return &ValidationError{Op: "get_cluster_profile", Field: "asset_class", Msg: "posting communities (wallstreetbets, reddit) apply only to equity"}
			}
		}
	}
	return nil
}

// ClusterRotationRequest carries the query parameters of `get_cluster_rotation` (GET /clusters/rotation).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-rotation/
type ClusterRotationRequest struct {
	ClusterID    string
	AssetClasses []AssetClass
	Hours        *int
	Limit        *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ClusterRotationRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("cluster_id", r.ClusterID)
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_cluster_rotation`
// before any HTTP call is made.
func (r *ClusterRotationRequest) validate() error {
	if r.ClusterID == "" {
		return &ValidationError{Op: "get_cluster_rotation", Field: "cluster_id", Msg: "required"}
	}
	if r.Hours != nil && !intIn(*r.Hours, []int{24, 168, 720}) {
		return &ValidationError{Op: "get_cluster_rotation", Field: "hours", Msg: "must be one of 24, 168, 720"}
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 20) {
		return &ValidationError{Op: "get_cluster_rotation", Field: "limit", Msg: "out of range [1, 20]"}
	}
	if isCommunity(r.ClusterID) && len(r.AssetClasses) > 0 {
		for _, c := range r.AssetClasses {
			if c != AssetClassEquity {
				return &ValidationError{Op: "get_cluster_rotation", Field: "asset_class", Msg: "posting communities (wallstreetbets, reddit) apply only to equity"}
			}
		}
	}
	return nil
}

// ClusterStancesRequest carries the query parameters of `get_cluster_stances` (GET /clusters/stances).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-clusters-stances/
type ClusterStancesRequest struct {
	ClusterID    string
	AssetClasses []AssetClass
	Days         *int
	Limit        *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ClusterStancesRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("cluster_id", r.ClusterID)
	if len(r.AssetClasses) > 0 {
		v.Set("asset_class", csvJoin(r.AssetClasses))
	}
	if r.Days != nil {
		v.Set("days", strconv.Itoa(*r.Days))
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_cluster_stances`
// before any HTTP call is made.
func (r *ClusterStancesRequest) validate() error {
	if r.ClusterID == "" {
		return &ValidationError{Op: "get_cluster_stances", Field: "cluster_id", Msg: "required"}
	}
	if r.Days != nil && !intIn(*r.Days, []int{7, 30, 90}) {
		return &ValidationError{Op: "get_cluster_stances", Field: "days", Msg: "must be one of 7, 30, 90"}
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "get_cluster_stances", Field: "limit", Msg: "out of range [1, 100]"}
	}
	if isCommunity(r.ClusterID) && len(r.AssetClasses) > 0 {
		for _, c := range r.AssetClasses {
			if c != AssetClassEquity {
				return &ValidationError{Op: "get_cluster_stances", Field: "asset_class", Msg: "posting communities (wallstreetbets, reddit) apply only to equity"}
			}
		}
	}
	return nil
}
