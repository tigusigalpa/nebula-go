// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"net/url"
	"strconv"
)

// EmotionsRequest carries the query parameters of `get_emotions` (GET /emotions).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-emotions/
type EmotionsRequest struct {
	Asset       *AssetRef
	Subject     *SubjectRef
	Hours       *int
	Granularity *Granularity
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *EmotionsRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Granularity != nil {
		v.Set("granularity", string(*r.Granularity))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_emotions`
// before any HTTP call is made.
func (r *EmotionsRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "get_emotions", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("get_emotions"); err != nil {
		return err
	}
	if r.Asset == nil && r.Subject == nil {
		return &ValidationError{Op: "get_emotions", Field: "selector", Msg: "asset or subject selector is required"}
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_emotions", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	return nil
}

// MindshareRequest carries the query parameters of `get_mindshare` (GET /mindshare).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-mindshare/
type MindshareRequest struct {
	Asset       *AssetRef
	Subject     *SubjectRef
	Hours       *int
	Granularity *Granularity
	Clusters    []string
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *MindshareRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Granularity != nil {
		v.Set("granularity", string(*r.Granularity))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_mindshare`
// before any HTTP call is made.
func (r *MindshareRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "get_mindshare", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("get_mindshare"); err != nil {
		return err
	}
	if r.Asset == nil && r.Subject == nil {
		return &ValidationError{Op: "get_mindshare", Field: "selector", Msg: "asset or subject selector is required"}
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_mindshare", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	if hasCommunity(r.Clusters) {
		if r.Subject != nil {
			return &ValidationError{Op: "get_mindshare", Field: "clusters", Msg: "posting communities apply to assets only"}
		}
		if r.Asset != nil && r.Asset.Class == AssetClassCrypto {
			return &ValidationError{Op: "get_mindshare", Field: "clusters", Msg: "posting communities require asset_class equity, commodity or index"}
		}
	}
	return nil
}

// RelatedSubjectsRequest carries the query parameters of `get_related_subjects` (GET /related-subjects).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-related-subjects/
type RelatedSubjectsRequest struct {
	Asset   *AssetRef
	Subject *SubjectRef
	Hours   *int
	Limit   *int
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *RelatedSubjectsRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Limit != nil {
		v.Set("limit", strconv.Itoa(*r.Limit))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_related_subjects`
// before any HTTP call is made.
func (r *RelatedSubjectsRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "get_related_subjects", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("get_related_subjects"); err != nil {
		return err
	}
	if r.Asset == nil && r.Subject == nil {
		return &ValidationError{Op: "get_related_subjects", Field: "selector", Msg: "asset or subject selector is required"}
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 2160) {
		return &ValidationError{Op: "get_related_subjects", Field: "hours", Msg: "out of range [1, 2160]"}
	}
	if r.Limit != nil && (*r.Limit < 1 || *r.Limit > 100) {
		return &ValidationError{Op: "get_related_subjects", Field: "limit", Msg: "out of range [1, 100]"}
	}
	return nil
}

// SentimentRequest carries the query parameters of `get_sentiment` (GET /sentiment).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-sentiment/
type SentimentRequest struct {
	Asset       *AssetRef
	Subject     *SubjectRef
	Hours       *int
	Granularity *Granularity
	Clusters    []string
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *SentimentRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Granularity != nil {
		v.Set("granularity", string(*r.Granularity))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_sentiment`
// before any HTTP call is made.
func (r *SentimentRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "get_sentiment", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("get_sentiment"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_sentiment", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	if hasCommunity(r.Clusters) {
		if r.Subject != nil {
			return &ValidationError{Op: "get_sentiment", Field: "clusters", Msg: "posting communities apply to assets only"}
		}
		if r.Asset != nil && r.Asset.Class == AssetClassCrypto {
			return &ValidationError{Op: "get_sentiment", Field: "clusters", Msg: "posting communities require asset_class equity, commodity or index"}
		}
	}
	return nil
}

// SentimentTimeframesRequest carries the query parameters of `get_sentiment_timeframes` (GET /sentiment/timeframes).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-sentiment-timeframes/
type SentimentTimeframesRequest struct {
	Asset       *AssetRef
	Subject     *SubjectRef
	Hours       *int
	Granularity *Granularity
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *SentimentTimeframesRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Granularity != nil {
		v.Set("granularity", string(*r.Granularity))
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_sentiment_timeframes`
// before any HTTP call is made.
func (r *SentimentTimeframesRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "get_sentiment_timeframes", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("get_sentiment_timeframes"); err != nil {
		return err
	}
	if r.Asset == nil && r.Subject == nil {
		return &ValidationError{Op: "get_sentiment_timeframes", Field: "selector", Msg: "asset or subject selector is required"}
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_sentiment_timeframes", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	return nil
}

// SocialMomentumRequest carries the query parameters of `get_social_momentum` (GET /social-momentum).
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-social-momentum/
type SocialMomentumRequest struct {
	Asset       *AssetRef
	Subject     *SubjectRef
	Hours       *int
	Granularity *Granularity
	Clusters    []string
}

// query serializes the request per the OpenAPI parameter contract:
// repeated params for `clusters`, CSV for style=form explode=false arrays.
func (r *SocialMomentumRequest) query() (url.Values, error) {
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
	if r.Hours != nil {
		v.Set("hours", strconv.Itoa(*r.Hours))
	}
	if r.Granularity != nil {
		v.Set("granularity", string(*r.Granularity))
	}
	for _, c := range r.Clusters {
		if c != "" {
			v.Add("clusters", c)
		}
	}
	return v, nil
}

// validate enforces the documented client-side rules of `get_social_momentum`
// before any HTTP call is made.
func (r *SocialMomentumRequest) validate() error {
	if r.Asset != nil && r.Subject != nil {
		return &ValidationError{Op: "get_social_momentum", Field: "subject", Msg: "subject and asset selectors are mutually exclusive"}
	}
	if err := r.Asset.optionalPair("get_social_momentum"); err != nil {
		return err
	}
	if r.Hours != nil && (*r.Hours < 1 || *r.Hours > 8760) {
		return &ValidationError{Op: "get_social_momentum", Field: "hours", Msg: "out of range [1, 8760]"}
	}
	if hasCommunity(r.Clusters) {
		return &ValidationError{Op: "get_social_momentum", Field: "clusters", Msg: "posting communities are not available on this endpoint"}
	}
	return nil
}
