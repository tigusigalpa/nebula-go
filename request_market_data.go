// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
)

// PriceHistoryRequest carries the query parameters of `get_price_history` (GET /price-history).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-price-history/
type PriceHistoryRequest struct {
	Asset       *AssetRef
	Hours       *int
	Granularity *Granularity
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *PriceHistoryRequest) query() (url.Values, error) {
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

// validate enforces the documented client-side rules of `get_price_history`
// before any HTTP call is made.
func (r *PriceHistoryRequest) validate() error {
	if err := r.Asset.requirePair("get_price_history"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_price_history", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	return nil
}
