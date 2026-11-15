package nebula

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Explicit service methods keep this audit independent of generated validation
// rules and derive request types from the actual public method signatures.
func contractQueryMethods(c *Client) map[string]any {
	return map[string]any{
		"screener":                  c.Discovery.Screener,
		"get_top_momentum":          c.Discovery.TopMomentum,
		"get_top_trending":          c.Discovery.TopTrending,
		"list_assets":               c.Lookup.ListAssets,
		"search":                    c.Lookup.Search,
		"list_subjects":             c.Lookup.ListSubjects,
		"get_emotions":              c.Social.Emotions,
		"get_mindshare":             c.Social.Mindshare,
		"get_related_subjects":      c.Social.RelatedSubjects,
		"get_sentiment":             c.Social.Sentiment,
		"get_sentiment_timeframes":  c.Social.SentimentTimeframes,
		"get_social_momentum":       c.Social.SocialMomentum,
		"list_forecast_assets":      c.Forecasts.ListAssets,
		"get_mindshare_forecast":    c.Forecasts.Mindshare,
		"get_sentiment_forecast":    c.Forecasts.Sentiment,
		"get_calls":                 c.MarketSentiment.Calls,
		"get_call_returns":          c.MarketSentiment.CallReturns,
		"get_conviction_index":      c.MarketSentiment.ConvictionIndex,
		"get_fear_greed_index":      c.MarketSentiment.FearGreedIndex,
		"get_price_levels":          c.MarketSentiment.PriceLevels,
		"get_social_volume_index":   c.MarketSentiment.SocialVolumeIndex,
		"get_social_volume_profile": c.MarketSentiment.SocialVolumeProfile,
		"get_price_history":         c.MarketData.PriceHistory,
		"list_calendar_events":      c.Intelligence.CalendarEvents,
		"list_insights":             c.Intelligence.Insights,
		"get_insight_attention":     c.Intelligence.InsightAttention,
		"list_clusters":             c.Clusters.List,
		"compare_clusters":          c.Clusters.Comparison,
		"get_divisive_assets":       c.Clusters.DivisiveAssets,
		"compare_hype_cycles":       c.Clusters.ComparisonHypeCycles,
		"compare_cluster_traits":    c.Clusters.ComparisonTraits,
		"get_cluster_hype_cycles":   c.Clusters.HypeCycles,
		"get_cluster_overlap":       c.Clusters.Overlap,
		"get_cluster_profile":       c.Clusters.Profile,
		"get_cluster_rotation":      c.Clusters.Rotation,
		"get_cluster_stances":       c.Clusters.Stances,
		"get_account":               c.Account.Get,
		"list_prediction_markets":   c.Predictions.ListMarkets,
		"get_prediction_history":    c.Predictions.History,
		"get_prediction_market":     c.Predictions.Market,
	}
}

func contractRequests(t *testing.T) map[string]queryRequest {
	t.Helper()
	c, err := New("contract-key")
	if err != nil {
		t.Fatal(err)
	}
	requests := map[string]queryRequest{}
	for id, method := range contractQueryMethods(c) {
		typ := reflect.TypeOf(method)
		if typ.NumIn() == 1 {
			requests[id] = nil
			continue
		}
		requests[id] = reflect.New(typ.In(1).Elem()).Interface().(queryRequest)
	}
	return requests
}

func TestQueryParametersMatchContract(t *testing.T) {
	ops := loadContract(t)
	requests := contractRequests(t)
	if len(requests)+1 != len(ops) {
		t.Fatal("request registry must cover all 40 GETs plus the chat POST")
	}
	for id, op := range ops {
		if id == "create_chat_completion" {
			continue
		}
		t.Run(id, func(t *testing.T) {
			prototype, exists := requests[id]
			if !exists {
				t.Fatal("request type missing")
			}
			if prototype == nil {
				if len(op.Parameters) != 0 {
					t.Fatal("parameterless service acquired parameters")
				}
				if wire := queryOnHTTPWire(t, id, nil); len(wire) != 0 {
					t.Fatalf("unexpected account query: %v", wire)
				}
				return
			}
			for _, param := range op.Parameters {
				if param.In != "query" {
					t.Fatalf("unhandled parameter location %s", param.In)
				}
				t.Run(param.Name, func(t *testing.T) {
					for i, value := range parameterSamples(param) {
						t.Run(strconv.Itoa(i), func(t *testing.T) {
							req := reflect.New(reflect.TypeOf(prototype).Elem())
							seedRequiredRequest(t, req.Elem(), op)
							field := queryParameterField(t, req.Elem(), param.Name)
							setContractValue(t, field, value)
							got, err := req.Interface().(queryRequest).query()
							if err != nil {
								t.Fatal(err)
							}
							if err := validateQuery(id, got); err != nil {
								t.Fatal(err)
							}
							want := parameterWireValues(param, value)
							if !reflect.DeepEqual(got[param.Name], want) {
								t.Errorf("encoded %v; want %v", got[param.Name], want)
							}
							// Query values must survive the same escaping used by HTTP.
							wire := queryOnHTTPWire(t, id, req.Interface().(queryRequest))
							if !reflect.DeepEqual(wire, got) {
								t.Fatalf("public method changed query values: wire=%v serializer=%v", wire, got)
							}
							for key := range got {
								if !contractHasParameter(op, key) {
									t.Errorf("undocumented query parameter %q", key)
								}
							}
							if param.Name == "asset_class" && standaloneClassFilter(param) && got.Has("asset_id") {
								t.Error("standalone asset_class filter unexpectedly selects an individual asset")
							}
						})
					}
				})
			}
		})
	}
}

func queryOnHTTPWire(t *testing.T, id string, req queryRequest) url.Values {
	t.Helper()
	var wire url.Values
	c, err := New("contract-key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		wire = r.URL.Query()
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: r}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	args := []reflect.Value{reflect.ValueOf(context.Background())}
	if req != nil {
		args = append(args, reflect.ValueOf(req))
	}
	out := reflect.ValueOf(contractQueryMethods(c)[id]).Call(args)
	if !out[1].IsNil() {
		t.Fatal(out[1].Interface())
	}
	if wire == nil {
		t.Fatal("public service method did not issue HTTP")
	}
	return wire
}

func contractHasParameter(op contractOperation, name string) bool {
	for _, p := range op.Parameters {
		if p.Name == name {
			return true
		}
	}
	return false
}

func standaloneClassFilter(param contractParameter) bool {
	return strings.Contains(param.Description, "Restricts results") || strings.Contains(param.Description, "market to aggregate")
}

func seedRequiredRequest(t *testing.T, req reflect.Value, op contractOperation) {
	t.Helper()
	for _, param := range op.Parameters {
		if param.Required {
			setContractValue(t, queryParameterField(t, req, param.Name), parameterSamples(param)[0])
		}
	}
	switch op.ID {
	case "get_emotions", "get_mindshare", "get_related_subjects", "get_sentiment_timeframes":
		req.FieldByName("Asset").Set(reflect.ValueOf(&AssetRef{ID: "NVDA", Class: AssetClassEquity}))
	}
}

func queryParameterField(t *testing.T, req reflect.Value, name string) reflect.Value {
	t.Helper()
	clear := func(field string) {
		if f := req.FieldByName(field); f.IsValid() {
			f.SetZero()
		}
	}
	switch name {
	case "asset_id":
		clear("Subject")
		clear("AssetClass")
		clear("Query")
		asset := req.FieldByName("Asset")
		if !asset.IsValid() {
			t.Fatal("asset selector is missing")
		}
		if asset.IsNil() {
			asset.Set(reflect.ValueOf(&AssetRef{ID: "NVDA", Class: AssetClassEquity}))
		}
		return asset.Elem().FieldByName("ID")
	case "asset_class":
		for _, field := range []string{"AssetClass", "AssetClasses"} {
			if f := req.FieldByName(field); f.IsValid() {
				clear("Asset")
				clear("Subject")
				return f
			}
		}
		asset := req.FieldByName("Asset")
		if !asset.IsValid() {
			t.Fatal("asset class selector is missing")
		}
		if asset.IsNil() {
			asset.Set(reflect.ValueOf(&AssetRef{ID: "NVDA", Class: AssetClassEquity}))
		}
		return asset.Elem().FieldByName("Class")
	case "subject_id":
		clear("Asset")
		clear("AssetClass")
		clear("Query")
		subject := req.FieldByName("Subject")
		if !subject.IsValid() {
			t.Fatal("subject selector is missing")
		}
		subject.Set(reflect.ValueOf(&SubjectRef{ID: "text:contract"}))
		return subject.Elem().FieldByName("ID")
	case "query":
		clear("Asset")
		clear("Subject")
	}
	for i := 0; i < req.NumField(); i++ {
		if strings.EqualFold(req.Type().Field(i).Name, strings.ReplaceAll(name, "_", "")) {
			return req.Field(i)
		}
	}
	t.Fatalf("%s has no public field for %s", req.Type(), name)
	return reflect.Value{}
}

func parameterSamples(p contractParameter) []any {
	s := p.Schema
	if s.Type == "array" {
		if len(s.Items.Enum) > 0 {
			out := []any{append([]any(nil), s.Items.Enum...)}
			for _, item := range s.Items.Enum {
				out = append(out, []any{item})
			}
			return out
		}
		if p.Name == "clusters" {
			return []any{[]any{"kols", "media"}}
		}
		return []any{[]any{"contract-one", "contract-two"}}
	}
	if len(s.Enum) > 0 {
		return s.Enum
	}
	switch s.Type {
	case "boolean":
		return []any{true, false}
	case "integer", "number":
		value := 17.0
		if s.Type == "number" {
			value = 0.125
		}
		if s.Minimum != nil && value < *s.Minimum {
			value = *s.Minimum
		}
		if s.Maximum != nil && value > *s.Maximum {
			value = *s.Maximum
		}
		out := []any{value}
		if s.Minimum != nil {
			out = append(out, *s.Minimum)
		}
		if s.Maximum != nil {
			out = append(out, *s.Maximum)
		}
		if (s.Minimum == nil || *s.Minimum <= 0) && (s.Maximum == nil || *s.Maximum >= 0) {
			out = append(out, float64(0))
		}
		return out
	case "string":
		if s.Format == "date-time" {
			return []any{"2026-10-27T06:04:05.123456789+03:00"}
		}
		value := "contract & value"
		if s.MaxLength > 0 && len(value) > s.MaxLength {
			value = value[:s.MaxLength]
		}
		return []any{value}
	default:
		panic("unsupported parameter type: " + s.Type)
	}
}

func setContractValue(t *testing.T, field reflect.Value, value any) {
	t.Helper()
	if field.Kind() == reflect.Pointer {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}
	if field.Type() == reflect.TypeFor[time.Time]() {
		parsed, err := time.Parse(time.RFC3339Nano, value.(string))
		if err != nil {
			t.Fatal(err)
		}
		field.Set(reflect.ValueOf(parsed))
		return
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(fmt.Sprint(value))
	case reflect.Int, reflect.Int64:
		n, err := strconv.ParseFloat(fmt.Sprint(value), 64)
		if err != nil {
			t.Fatal(err)
		}
		if n != float64(int64(n)) || field.OverflowInt(int64(n)) {
			t.Fatalf("cannot represent %v as %s", value, field.Type())
		}
		field.SetInt(int64(n))
	case reflect.Float64:
		n, err := strconv.ParseFloat(fmt.Sprint(value), 64)
		if err != nil {
			t.Fatal(err)
		}
		field.SetFloat(n)
	case reflect.Bool:
		field.SetBool(value.(bool))
	case reflect.Slice:
		values := value.([]any)
		field.Set(reflect.MakeSlice(field.Type(), len(values), len(values)))
		for i, item := range values {
			setContractValue(t, field.Index(i), item)
		}
	default:
		t.Fatalf("cannot populate query field of type %s", field.Type())
	}
}

func parameterWireValues(p contractParameter, value any) []string {
	if values, ok := value.([]any); ok {
		out := make([]string, len(values))
		for i, item := range values {
			out[i] = fmt.Sprint(item)
		}
		if p.Explode != nil && !*p.Explode {
			return []string{strings.Join(out, ",")}
		}
		return out
	}
	if p.Schema.Format == "date-time" {
		parsed, _ := time.Parse(time.RFC3339Nano, value.(string))
		return []string{parsed.UTC().Format(time.RFC3339Nano)}
	}
	if p.Schema.Type == "integer" {
		n, _ := strconv.ParseFloat(fmt.Sprint(value), 64)
		return []string{strconv.FormatFloat(n, 'f', -1, 64)}
	}
	return []string{fmt.Sprint(value)}
}

func TestInvalidQueryParametersFailValidation(t *testing.T) {
	requests := contractRequests(t)
	for id, op := range loadContract(t) {
		prototype := requests[id]
		if prototype == nil {
			continue
		}
		for _, param := range op.Parameters {
			bad := []any{}
			s := param.Schema
			if s.Type == "integer" || s.Type == "number" {
				if s.Minimum != nil {
					bad = append(bad, *s.Minimum-1)
				}
				if s.Maximum != nil {
					bad = append(bad, *s.Maximum+1)
				}
			}
			if len(s.Enum) > 0 {
				if s.Type == "integer" || s.Type == "number" {
					bad = append(bad, float64(-123))
				} else {
					bad = append(bad, "invalid-contract-enum")
				}
			}
			if s.Type == "array" && len(s.Items.Enum) > 0 {
				bad = append(bad, []any{"invalid-contract-enum"})
			}
			if s.MinLength > 0 {
				bad = append(bad, "")
			}
			if s.MaxLength > 0 {
				bad = append(bad, strings.Repeat("x", s.MaxLength+1))
			}
			for i, value := range bad {
				t.Run(fmt.Sprintf("%s/%s/%d", id, param.Name, i), func(t *testing.T) {
					req := reflect.New(reflect.TypeOf(prototype).Elem())
					seedRequiredRequest(t, req.Elem(), op)
					setContractValue(t, queryParameterField(t, req.Elem(), param.Name), value)
					assertQueryFailsBeforeHTTP(t, id, req.Interface().(queryRequest))
				})
			}
			if param.Required {
				t.Run(id+"/"+param.Name+"/missing", func(t *testing.T) {
					req := reflect.New(reflect.TypeOf(prototype).Elem())
					seedRequiredRequest(t, req.Elem(), op)
					queryParameterField(t, req.Elem(), param.Name).SetZero()
					assertQueryFailsBeforeHTTP(t, id, req.Interface().(queryRequest))
				})
			}
		}
	}
}

func assertQueryFailsBeforeHTTP(t *testing.T, id string, req queryRequest) {
	t.Helper()
	calls := 0
	c, err := New("contract-key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(`{}`)), Request: r}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	out := reflect.ValueOf(contractQueryMethods(c)[id]).Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(req)})
	if out[1].IsNil() || !errors.Is(out[1].Interface().(error), ErrValidation) {
		t.Fatalf("invalid query must produce ErrValidation: %v", out[1])
	}
	if calls != 0 {
		t.Fatalf("invalid query reached HTTP: %d calls", calls)
	}
}

func TestChatBodyParametersMatchContract(t *testing.T) {
	schema := loadContract(t)["create_chat_completion"].RequestBody.Content["application/json"].Schema
	for _, mode := range []ThinkingMode{ThinkingModeAuto, ThinkingModeDeep} {
		for _, streaming := range []bool{false, true} {
			req := &ChatCompletionRequest{Model: StringPtr("compatibility-model"), ThinkingMode: &mode}
			for _, role := range []ChatRole{ChatRoleSystem, ChatRoleDeveloper, ChatRoleUser, ChatRoleAssistant} {
				req.Messages = append(req.Messages, ChatMessage{Role: role, Content: "Message for " + string(role)})
			}
			if err := req.validate(); err != nil {
				t.Fatal(err)
			}
			body := req.body(streaming)
			checkResponseType(t, schema, reflect.TypeOf(body), "chat request")
			encoded, err := json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
			var actual map[string]any
			if err := json.Unmarshal(encoded, &actual); err != nil {
				t.Fatal(err)
			}
			if len(actual) != len(schema.Properties) {
				t.Fatalf("body fields=%v; schema=%v", actual, schema.Properties)
			}
			if actual["stream"] != streaming || actual["model"] != *req.Model || actual["thinking_mode"] != string(mode) {
				t.Fatalf("chat parameters lost: %s", encoded)
			}
			messages := actual["messages"].([]any)
			if len(messages) != len(req.Messages) {
				t.Fatal("message history truncated")
			}
			for i, msg := range req.Messages {
				got := messages[i].(map[string]any)
				if got["role"] != string(msg.Role) || got["content"] != msg.Content {
					t.Fatalf("message %d lost: %s", i, encoded)
				}
			}
			var wire []byte
			client, err := New("contract-key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				wire, err = io.ReadAll(r.Body)
				r.Body.Close()
				if err != nil {
					return nil, err
				}
				if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/chat/completions") {
					t.Errorf("wrong chat route: %s %s", r.Method, r.URL.Path)
				}
				header, response := http.Header{}, `{}`
				if streaming {
					header.Set("Content-Type", "text/event-stream")
					response = "data: [DONE]\n\n"
				}
				return &http.Response{StatusCode: 200, Header: header, Body: io.NopCloser(strings.NewReader(response)), Request: r}, nil
			})}))
			if err != nil {
				t.Fatal(err)
			}
			req.Stream = BoolPtr(!streaming) // The selected method owns streaming mode.
			if streaming {
				stream, err := client.Agent.StreamChatCompletion(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}
				for range stream.Events() {
				}
				if err := stream.Err(); err != nil {
					t.Fatal(err)
				}
				stream.Close()
			} else if _, err := client.Agent.CreateChatCompletion(context.Background(), req); err != nil {
				t.Fatal(err)
			}
			compareJSONLeaves(t, decodeContractJSON(t, encoded), decodeContractJSON(t, wire), "chat HTTP body")
		}
	}
}

func TestPredictionClassFilterSelectors(t *testing.T) {
	for _, req := range []*PredictionsRequest{
		{AssetClass: AssetClassPtr(AssetClassCrypto), Subject: &SubjectRef{ID: "text:example"}},
		{AssetClass: AssetClassPtr(AssetClassEquity), Asset: &AssetRef{ID: "NVDA", Class: AssetClassEquity}},
		{AssetClass: AssetClassPtr("invalid")},
	} {
		if _, err := req.query(); !errors.Is(err, ErrValidation) {
			t.Fatalf("invalid selectors accepted: %+v, %v", req, err)
		}
	}
}
