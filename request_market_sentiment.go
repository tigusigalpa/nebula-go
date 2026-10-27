// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
)

// CallsRequest carries the query parameters of `get_calls` (GET /calls).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-calls/
type CallsRequest struct {
	Asset    *AssetRef
	Hours    *int
	Clusters []string
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *CallsRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_calls`
// before any HTTP call is made.
func (r *CallsRequest) validate() error {
	if err := r.Asset.optionalPair("get_calls"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 4320) {
		return &ValidationError{Op: "get_calls", Field: "hours", Msg: "out of range [1, 4320]"}
	}
	if hasCommunity(r.Clusters) && r.Asset != nil && r.Asset.Class == AssetClassCrypto {
		return &ValidationError{Op: "get_calls", Field: "clusters", Msg: "posting communities require asset_class equity, commodity or index"}
	}
	return nil
}

// CallReturnsRequest carries the query parameters of `get_call_returns` (GET /calls/returns).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-calls-returns/
type CallReturnsRequest struct {
	Asset        *AssetRef
	ResolvedDays *int
	Clusters     []string
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *CallReturnsRequest) query() (url.Values, error) {
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
	if r.ResolvedDays != nil {
		v.Set("resolved_days", strconv.Itoa(*r.ResolvedDays))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_call_returns`
// before any HTTP call is made.
func (r *CallReturnsRequest) validate() error {
	if err := r.Asset.optionalPair("get_call_returns"); err != nil {
		return err
	}
	if r.ResolvedDays != nil && (*r.ResolvedDays < 1 || *r.ResolvedDays > 90) {
		return &ValidationError{Op: "get_call_returns", Field: "resolved_days", Msg: "out of range [1, 90]"}
	}
	if hasCommunity(r.Clusters) && r.Asset != nil && r.Asset.Class == AssetClassCrypto {
		return &ValidationError{Op: "get_call_returns", Field: "clusters", Msg: "posting communities require asset_class equity, commodity or index"}
	}
	return nil
}

// ConvictionIndexRequest carries the query parameters of `get_conviction_index` (GET /conviction-index).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-conviction-index/
type ConvictionIndexRequest struct {
	Asset       *AssetRef
	Hours       *int
	Granularity *Granularity
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ConvictionIndexRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Granularity != nil {
		v.Set("granularity", string(*r.Granularity))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_conviction_index`
// before any HTTP call is made.
func (r *ConvictionIndexRequest) validate() error {
	if err := r.Asset.requirePair("get_conviction_index"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_conviction_index", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	return nil
}

// FearGreedIndexRequest carries the query parameters of `get_fear_greed_index` (GET /fear-greed-index).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-fear-greed-index/
type FearGreedIndexRequest struct {
	Asset       *AssetRef
	Hours       *int
	Granularity *Granularity
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *FearGreedIndexRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Granularity != nil {
		v.Set("granularity", string(*r.Granularity))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_fear_greed_index`
// before any HTTP call is made.
func (r *FearGreedIndexRequest) validate() error {
	if err := r.Asset.requirePair("get_fear_greed_index"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_fear_greed_index", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	return nil
}

// PriceLevelsRequest carries the query parameters of `get_price_levels` (GET /price-levels).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-price-levels/
type PriceLevelsRequest struct {
	Asset *AssetRef
	Hours *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *PriceLevelsRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_price_levels`
// before any HTTP call is made.
func (r *PriceLevelsRequest) validate() error {
	if err := r.Asset.requirePair("get_price_levels"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_price_levels", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	return nil
}

// SocialVolumeIndexRequest carries the query parameters of `get_social_volume_index` (GET /social-volume-index).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-social-volume-index/
type SocialVolumeIndexRequest struct {
	Asset *AssetRef
	Hours *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *SocialVolumeIndexRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_social_volume_index`
// before any HTTP call is made.
func (r *SocialVolumeIndexRequest) validate() error {
	if err := r.Asset.requirePair("get_social_volume_index"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 2160) {
		return &ValidationError{Op: "get_social_volume_index", Field: "hours", Msg: "out of range [1, 2160]"}
	}
	return nil
}

// SocialVolumeProfileRequest carries the query parameters of `get_social_volume_profile` (GET /social-volume-profile).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-social-volume-profile/
type SocialVolumeProfileRequest struct {
	Asset    *AssetRef
	Hours    *int
	Bins     *int
	Clusters []string
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *SocialVolumeProfileRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Bins != nil {
		v.Set("bins", strconv.Itoa(*r.Bins))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_social_volume_profile`
// before any HTTP call is made.
func (r *SocialVolumeProfileRequest) validate() error {
	if err := r.Asset.requirePair("get_social_volume_profile"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_social_volume_profile", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	if r.Bins != nil && (*r.Bins < 10 || *r.Bins > 100) {
		return &ValidationError{Op: "get_social_volume_profile", Field: "bins", Msg: "out of range [10, 100]"}
	}
	if hasCommunity(r.Clusters) && r.Asset.Class == AssetClassCrypto {
		return &ValidationError{Op: "get_social_volume_profile", Field: "clusters", Msg: "posting communities require asset_class equity, commodity or index"}
	}
	return nil
}
