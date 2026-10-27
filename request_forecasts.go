// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
)

// ForecastAssetsRequest carries the query parameters of `list_forecast_assets` (GET /forecasts/assets).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-forecasts-assets/
type ForecastAssetsRequest struct {
	AssetClass *AssetClass
	Limit      *int
	Offset     *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ForecastAssetsRequest) query() (url.Values, error) {
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
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	if r.Offset != nil {
		v.Set("offset", strconv.Itoa(*r.Offset))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `list_forecast_assets`
// before any HTTP call is made.
func (r *ForecastAssetsRequest) validate() error {
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 1000) {
		return &ValidationError{Op: "list_forecast_assets", Field: "limit", Msg: "out of range [1, 1000]"}
	}
	if r.Offset != nil && (*r.Offset < 0 || *r.Offset > 100000) {
		return &ValidationError{Op: "list_forecast_assets", Field: "offset", Msg: "out of range [0, 100000]"}
	}
	return nil
}

// MindshareForecastRequest carries the query parameters of `get_mindshare_forecast` (GET /forecasts/mindshare).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-forecasts-mindshare/
type MindshareForecastRequest struct {
	Asset   *AssetRef
	Refresh *bool
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *MindshareForecastRequest) query() (url.Values, error) {
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
	if r.Refresh != nil {
		v.Set("refresh", strconv.FormatBool(*r.Refresh))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_mindshare_forecast`
// before any HTTP call is made.
func (r *MindshareForecastRequest) validate() error {
	if err := r.Asset.requirePair("get_mindshare_forecast"); err != nil {
		return err
	}
	return nil
}

// SentimentForecastRequest carries the query parameters of `get_sentiment_forecast` (GET /forecasts/sentiment).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-forecasts-sentiment/
type SentimentForecastRequest struct {
	Asset   *AssetRef
	Refresh *bool
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *SentimentForecastRequest) query() (url.Values, error) {
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
	if r.Refresh != nil {
		v.Set("refresh", strconv.FormatBool(*r.Refresh))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_sentiment_forecast`
// before any HTTP call is made.
func (r *SentimentForecastRequest) validate() error {
	if err := r.Asset.requirePair("get_sentiment_forecast"); err != nil {
		return err
	}
	return nil
}
