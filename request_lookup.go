// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
)

// ListAssetsRequest carries the query parameters of `list_assets` (GET /assets).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-assets/
type ListAssetsRequest struct {
	Asset  *AssetRef
	Query  *string
	Limit  *int
	Offset *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ListAssetsRequest) query() (url.Values, error) {
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
	if r.Query != nil {
		v.Set("query", *r.Query)
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	if r.Offset != nil {
		v.Set("offset", strconv.Itoa(*r.Offset))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `list_assets`
// before any HTTP call is made.
func (r *ListAssetsRequest) validate() error {
	if err := r.Asset.optionalPair("list_assets"); err != nil {
		return err
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "list_assets", Field: "limit", Msg: "out of range [1, 100]"}
	}
	if r.Offset != nil && (*r.Offset < 0 || *r.Offset > 100000) {
		return &ValidationError{Op: "list_assets", Field: "offset", Msg: "out of range [0, 100000]"}
	}
	if r.Asset != nil && r.Query != nil {
		return &ValidationError{Op: "list_assets", Field: "query", Msg: "asset_id cannot be combined with query"}
	}
	return nil
}

// SearchRequest carries the query parameters of `search` (GET /search).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-search/
type SearchRequest struct {
	Query string
	Kind  *Kind
	Limit *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *SearchRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("query", r.Query)
	if r.Kind != nil {
		v.Set("kind", string(*r.Kind))
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `search`
// before any HTTP call is made.
func (r *SearchRequest) validate() error {
	if r.Query == "" {
		return &ValidationError{Op: "search", Field: "query", Msg: "required"}
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 25) {
		return &ValidationError{Op: "search", Field: "limit", Msg: "out of range [1, 25]"}
	}
	return nil
}

// ListSubjectsRequest carries the query parameters of `list_subjects` (GET /subjects).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-subjects/
type ListSubjectsRequest struct {
	Subject *SubjectRef
	Query   *string
	Limit   *int
	Offset  *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *ListSubjectsRequest) query() (url.Values, error) {
	if r == nil {
		return url.Values{}, nil
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	if r.Subject != nil {
		v.Set("subject_id", r.Subject.ID)
	}
	if r.Query != nil {
		v.Set("query", *r.Query)
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	if r.Offset != nil {
		v.Set("offset", strconv.Itoa(*r.Offset))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `list_subjects`
// before any HTTP call is made.
func (r *ListSubjectsRequest) validate() error {
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "list_subjects", Field: "limit", Msg: "out of range [1, 100]"}
	}
	if r.Offset != nil && (*r.Offset < 0 || *r.Offset > 100000) {
		return &ValidationError{Op: "list_subjects", Field: "offset", Msg: "out of range [0, 100000]"}
	}
	if r.Subject != nil && r.Query != nil {
		return &ValidationError{Op: "list_subjects", Field: "query", Msg: "subject_id cannot be combined with query"}
	}
	return nil
}
