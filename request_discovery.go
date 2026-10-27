// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
)

// ScreenerRequest carries the query parameters of `screener` (GET /screener).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-screener/
type ScreenerRequest struct {
	Kind                 *Kind
	AssetClass           *AssetClass
	Clusters             []string
	Sort                 *ScreenerSort
	Order                *SortOrder
	Limit                *int
	MinShare             *float64
	MaxShare             *float64
	MinShareChange       *float64
	MaxShareChange       *float64
	MinPosts             *int
	MinNetBullishPosts   *int
	MaxNetBullishPosts   *int
	MinPriceChange24hPct *float64
	MaxPriceChange24hPct *float64
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ScreenerRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if r.Kind != nil {
		v.Set("kind", string(*r.Kind))
	}
	if r.AssetClass != nil {
		v.Set("asset_class", string(*r.AssetClass))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	if r.Sort != nil {
		v.Set("sort", string(*r.Sort))
	}
	if r.Order != nil {
		v.Set("order", string(*r.Order))
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	if r.MinShare != nil {
		v.Set("min_share", strconv.FormatFloat(*r.MinShare, 'g', -1, 64))
	}
	if r.MaxShare != nil {
		v.Set("max_share", strconv.FormatFloat(*r.MaxShare, 'g', -1, 64))
	}
	if r.MinShareChange != nil {
		v.Set("min_share_change", strconv.FormatFloat(*r.MinShareChange, 'g', -1, 64))
	}
	if r.MaxShareChange != nil {
		v.Set("max_share_change", strconv.FormatFloat(*r.MaxShareChange, 'g', -1, 64))
	}
	if r.MinPosts != nil {
		v.Set("min_posts", strconv.Itoa(*r.MinPosts))
	}
	if r.MinNetBullishPosts != nil {
		v.Set("min_net_bullish_posts", strconv.Itoa(*r.MinNetBullishPosts))
	}
	if r.MaxNetBullishPosts != nil {
		v.Set("max_net_bullish_posts", strconv.Itoa(*r.MaxNetBullishPosts))
	}
	if r.MinPriceChange24hPct != nil {
		v.Set("min_price_change_24h_pct", strconv.FormatFloat(*r.MinPriceChange24hPct, 'g', -1, 64))
	}
	if r.MaxPriceChange24hPct != nil {
		v.Set("max_price_change_24h_pct", strconv.FormatFloat(*r.MaxPriceChange24hPct, 'g', -1, 64))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `screener`
// before any HTTP call is made.
func (r *ScreenerRequest) validate() error {
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "screener", Field: "limit", Msg: "out of range [1, 100]"}
	}
	if r.MinShare != nil && (*r.MinShare < 0 || *r.MinShare > 1) {
		return &ValidationError{Op: "screener", Field: "min_share", Msg: "out of range [0, 1]"}
	}
	if r.MaxShare != nil && (*r.MaxShare < 0 || *r.MaxShare > 1) {
		return &ValidationError{Op: "screener", Field: "max_share", Msg: "out of range [0, 1]"}
	}
	if r.MinShareChange != nil && (*r.MinShareChange < -1 || *r.MinShareChange > 1) {
		return &ValidationError{Op: "screener", Field: "min_share_change", Msg: "out of range [-1, 1]"}
	}
	if r.MaxShareChange != nil && (*r.MaxShareChange < -1 || *r.MaxShareChange > 1) {
		return &ValidationError{Op: "screener", Field: "max_share_change", Msg: "out of range [-1, 1]"}
	}
	if r.MinPosts != nil && (*r.MinPosts < 1 || *r.MinPosts > 1000000000) {
		return &ValidationError{Op: "screener", Field: "min_posts", Msg: "out of range [1, 1000000000]"}
	}
	if r.MinNetBullishPosts != nil && (*r.MinNetBullishPosts < -1000000000 || *r.MinNetBullishPosts > 1000000000) {
		return &ValidationError{Op: "screener", Field: "min_net_bullish_posts", Msg: "out of range [-1000000000, 1000000000]"}
	}
	if r.MaxNetBullishPosts != nil && (*r.MaxNetBullishPosts < -1000000000 || *r.MaxNetBullishPosts > 1000000000) {
		return &ValidationError{Op: "screener", Field: "max_net_bullish_posts", Msg: "out of range [-1000000000, 1000000000]"}
	}
	if r.MinPriceChange24hPct != nil && (*r.MinPriceChange24hPct < -100 || *r.MinPriceChange24hPct > 100000) {
		return &ValidationError{Op: "screener", Field: "min_price_change_24h_pct", Msg: "out of range [-100, 100000]"}
	}
	if r.MaxPriceChange24hPct != nil && (*r.MaxPriceChange24hPct < -100 || *r.MaxPriceChange24hPct > 100000) {
		return &ValidationError{Op: "screener", Field: "max_price_change_24h_pct", Msg: "out of range [-100, 100000]"}
	}
	kind := KindAsset
	if r.Kind != nil {
		kind = *r.Kind
	}
	if kind == KindSubject && r.AssetClass != nil {
		return &ValidationError{Op: "screener", Field: "asset_class", Msg: "asset_class applies only to kind=asset"}
	}
	if kind == KindSubject {
		if r.Sort != nil && *r.Sort == ScreenerSortPriceChange24hPct {
			return &ValidationError{Op: "screener", Field: "sort", Msg: "price sorting requires kind=asset"}
		}
		if r.MinPriceChange24hPct != nil || r.MaxPriceChange24hPct != nil {
			return &ValidationError{Op: "screener", Field: "price_change_24h_pct", Msg: "price filters require kind=asset"}
		}
	}
	if err := minLTEPtr(r.MinShare, r.MaxShare, "screener", "min_share", "max_share"); err != nil {
		return err
	}
	if err := minLTEPtr(r.MinShareChange, r.MaxShareChange, "screener", "min_share_change", "max_share_change"); err != nil {
		return err
	}
	if err := minLTEPtr(r.MinNetBullishPosts, r.MaxNetBullishPosts, "screener", "min_net_bullish_posts", "max_net_bullish_posts"); err != nil {
		return err
	}
	if err := minLTEPtr(r.MinPriceChange24hPct, r.MaxPriceChange24hPct, "screener", "min_price_change_24h_pct", "max_price_change_24h_pct"); err != nil {
		return err
	}
	if hasCommunity(r.Clusters) {
		class := AssetClassEquity
		if r.AssetClass != nil {
			class = *r.AssetClass
		}
		if kind != KindAsset || class == AssetClassCrypto {
			return &ValidationError{Op: "screener", Field: "clusters", Msg: "posting communities (wallstreetbets, reddit) require kind=asset and asset_class equity, commodity or index"}
		}
	}
	return nil
}

// TopMomentumRequest carries the query parameters of `get_top_momentum` (GET /top-momentum).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-top-momentum/
type TopMomentumRequest struct {
	Kind       *Kind
	AssetClass *AssetClass
	Clusters   []string
	Limit      *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *TopMomentumRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if r.Kind != nil {
		v.Set("kind", string(*r.Kind))
	}
	if r.AssetClass != nil {
		v.Set("asset_class", string(*r.AssetClass))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_top_momentum`
// before any HTTP call is made.
func (r *TopMomentumRequest) validate() error {
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "get_top_momentum", Field: "limit", Msg: "out of range [1, 100]"}
	}
	if r.Kind != nil && *r.Kind == KindSubject && r.AssetClass != nil {
		return &ValidationError{Op: "get_top_momentum", Field: "asset_class", Msg: "asset_class applies only to kind=asset"}
	}
	if hasCommunity(r.Clusters) {
		if r.Kind != nil && *r.Kind == KindSubject {
			return &ValidationError{Op: "get_top_momentum", Field: "clusters", Msg: "posting communities apply to kind=asset"}
		}
		if r.AssetClass != nil && *r.AssetClass == AssetClassCrypto {
			return &ValidationError{Op: "get_top_momentum", Field: "clusters", Msg: "posting communities require asset_class equity, commodity or index"}
		}
	}
	return nil
}

// TopTrendingRequest carries the query parameters of `get_top_trending` (GET /top-trending).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-top-trending/
type TopTrendingRequest struct {
	Kind       *Kind
	AssetClass *AssetClass
	Clusters   []string
	Limit      *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *TopTrendingRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if r.Kind != nil {
		v.Set("kind", string(*r.Kind))
	}
	if r.AssetClass != nil {
		v.Set("asset_class", string(*r.AssetClass))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_top_trending`
// before any HTTP call is made.
func (r *TopTrendingRequest) validate() error {
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "get_top_trending", Field: "limit", Msg: "out of range [1, 100]"}
	}
	if r.Kind != nil && *r.Kind == KindSubject && r.AssetClass != nil {
		return &ValidationError{Op: "get_top_trending", Field: "asset_class", Msg: "asset_class applies only to kind=asset"}
	}
	if hasCommunity(r.Clusters) {
		if r.Kind != nil && *r.Kind == KindSubject {
			return &ValidationError{Op: "get_top_trending", Field: "clusters", Msg: "posting communities apply to kind=asset"}
		}
		if r.AssetClass != nil && *r.AssetClass == AssetClassCrypto {
			return &ValidationError{Op: "get_top_trending", Field: "clusters", Msg: "posting communities require asset_class equity, commodity or index"}
		}
	}
	return nil
}
