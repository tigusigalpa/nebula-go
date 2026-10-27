// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
)

// PredictionsRequest carries the query parameters of `list_prediction_markets` (GET /predictions).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-predictions/
type PredictionsRequest struct {
	Asset    *AssetRef
	Subject  *SubjectRef
	Query    *string
	Provider *PredictionProvider
	Status   *PredictionStatus
	Limit    *int
	Offset   *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *PredictionsRequest) query() (url.Values, error) {
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
	if r.Query != nil {
		v.Set("query", *r.Query)
	}
	if r.Provider != nil {
		v.Set("provider", string(*r.Provider))
	}
	if r.Status != nil {
		v.Set("status", string(*r.Status))
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	if r.Offset != nil {
		v.Set("offset", strconv.Itoa(*r.Offset))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `list_prediction_markets`
// before any HTTP call is made.
func (r *PredictionsRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "list_prediction_markets", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("list_prediction_markets"); err != nil {
		return err
	}
	if r.Query != nil && len(*r.Query) < 1 {
		return &ValidationError{Op: "list_prediction_markets", Field: "query", Msg: "min length 1"}
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "list_prediction_markets", Field: "limit", Msg: "out of range [1, 100]"}
	}
	if r.Offset != nil && (*r.Offset < 0 || *r.Offset > 10000) {
		return &ValidationError{Op: "list_prediction_markets", Field: "offset", Msg: "out of range [0, 10000]"}
	}
	return nil
}

// PredictionHistoryRequest carries the query parameters of `get_prediction_history` (GET /predictions/history).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-prediction-history/
type PredictionHistoryRequest struct {
	ListingID string
	OutcomeID string
	Hours     *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *PredictionHistoryRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("listing_id", r.ListingID)
	v.Set("outcome_id", r.OutcomeID)
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_prediction_history`
// before any HTTP call is made.
func (r *PredictionHistoryRequest) validate() error {
	if len(r.ListingID) < 1 {
		return &ValidationError{Op: "get_prediction_history", Field: "listing_id", Msg: "min length 1"}
	}
	if r.ListingID == "" {
		return &ValidationError{Op: "get_prediction_history", Field: "listing_id", Msg: "required"}
	}
	if len(r.OutcomeID) < 1 {
		return &ValidationError{Op: "get_prediction_history", Field: "outcome_id", Msg: "min length 1"}
	}
	if r.OutcomeID == "" {
		return &ValidationError{Op: "get_prediction_history", Field: "outcome_id", Msg: "required"}
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8784) {
		return &ValidationError{Op: "get_prediction_history", Field: "hours", Msg: "out of range [1, 8784]"}
	}
	return nil
}

// PredictionMarketRequest carries the query parameters of `get_prediction_market` (GET /predictions/market).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-prediction-market/
type PredictionMarketRequest struct {
	MarketID string
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *PredictionMarketRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("market_id", r.MarketID)
	return v, nil
}

// validate enforces the documented client-side rules of `get_prediction_market`
// before any HTTP call is made.
func (r *PredictionMarketRequest) validate() error {
	if len(r.MarketID) < 1 {
		return &ValidationError{Op: "get_prediction_market", Field: "market_id", Msg: "min length 1"}
	}
	if r.MarketID == "" {
		return &ValidationError{Op: "get_prediction_market", Field: "market_id", Msg: "required"}
	}
	return nil
}
