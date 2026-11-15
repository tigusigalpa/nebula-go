package nebula

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// This is a test-only view of the pinned contract, independent of generated
// validation rules. Unknown schema constructs must fail rather than be ignored.
type contractSchema struct {
	Type                 string                     `json:"type"`
	Format               string                     `json:"format"`
	Nullable             bool                       `json:"nullable"`
	Properties           map[string]*contractSchema `json:"properties"`
	Items                *contractSchema            `json:"items"`
	Required             []string                   `json:"required"`
	Enum                 []any                      `json:"enum"`
	Minimum              *float64                   `json:"minimum"`
	Maximum              *float64                   `json:"maximum"`
	MinLength            int                        `json:"minLength"`
	MaxLength            int                        `json:"maxLength"`
	Example              any                        `json:"example"`
	MinItems             int                        `json:"minItems"`
	AdditionalProperties json.RawMessage            `json:"additionalProperties"`
}

type contractParameter struct {
	Name        string          `json:"name"`
	In          string          `json:"in"`
	Description string          `json:"description"`
	Required    bool            `json:"required"`
	Style       string          `json:"style"`
	Explode     *bool           `json:"explode"`
	Example     any             `json:"example"`
	Schema      *contractSchema `json:"schema"`
}

type contractMedia struct {
	Schema *contractSchema `json:"schema"`
}
type contractBody struct {
	Content map[string]contractMedia `json:"content"`
}
type contractOperation struct {
	ID          string                  `json:"operationId"`
	Parameters  []contractParameter     `json:"parameters"`
	RequestBody contractBody            `json:"requestBody"`
	Responses   map[string]contractBody `json:"responses"`
}

func loadContract(t *testing.T) map[string]contractOperation {
	t.Helper()
	b, err := os.ReadFile("resources/openapi/nebula-openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		Paths map[string]map[string]contractOperation `json:"paths"`
	}
	if err := json.Unmarshal(b, &spec); err != nil {
		t.Fatal(err)
	}
	ops := map[string]contractOperation{}
	for _, methods := range spec.Paths {
		for method, op := range methods {
			if method == "get" || method == "post" {
				if _, exists := ops[op.ID]; exists {
					t.Fatalf("duplicate operation %s", op.ID)
				}
				ops[op.ID] = op
			}
		}
	}
	if len(ops) != 41 {
		t.Fatalf("expected 41 contract operations, got %d", len(ops))
	}
	return ops
}

// TestFixtureDecodes checks actual service return types, including every JSON
// leaf in the saved fixture. Complete probes cover fields missing from fixtures;
// nullable probes distinguish null from an incorrectly decoded scalar zero.
func TestFixtureDecodes(t *testing.T) {
	ops := loadContract(t)
	steps := coverageSteps(context.Background())
	if len(steps) != len(ops) {
		t.Fatal("fixture/service matrix is incomplete")
	}
	for _, step := range steps {
		t.Run(step.op, func(t *testing.T) {
			op, ok := ops[step.op]
			if !ok {
				t.Fatal("operation missing from OpenAPI")
			}
			schema := op.Responses["200"].Content["application/json"].Schema
			if schema == nil {
				t.Fatal("missing response schema")
			}
			fixture, err := os.ReadFile("testdata/" + step.op + ".json")
			if err != nil {
				t.Fatal(err)
			}
			checkFixtureContract(t, schema, decodeContractJSON(t, fixture), step.op)
			payloads := map[string][]byte{"fixture": fixture}
			for _, mode := range []string{"complete", "nullable", "nullable-leaves", "empty-collections"} {
				payloads[mode], err = json.Marshal(responseProbe(schema, mode))
				if err != nil {
					t.Fatal(err)
				}
			}
			for name, body := range payloads {
				t.Run(name, func(t *testing.T) {
					client, err := New("test-key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{StatusCode: 200, Header: http.Header{"X-Request-Id": []string{"contract-request"}}, Body: io.NopCloser(bytes.NewReader(body)), Request: r}, nil
					})}))
					if err != nil {
						t.Fatal(err)
					}
					result, err := step.call(t, client)
					if err != nil {
						t.Fatal(err)
					}
					if result == nil {
						t.Fatal("nil typed response")
					}
					checkResponseType(t, schema, reflect.TypeOf(result), step.op)
					encoded, err := json.Marshal(result)
					if err != nil {
						t.Fatal(err)
					}
					compareJSONLeaves(t, decodeContractJSON(t, body), decodeContractJSON(t, encoded), step.op)
					carrier, ok := result.(interface{ HTTPMetadata() ResponseMetadata })
					if !ok || carrier.HTTPMetadata().RequestID != "contract-request" {
						t.Fatal("typed response lost HTTP metadata")
					}
				})
			}
			// Unknown upstream fields remain available without weakening the DTO.
			var future map[string]any
			if err := json.Unmarshal(fixture, &future); err != nil {
				t.Fatal(err)
			}
			future["future_field"] = "forward-compatible"
			body, err := json.Marshal(future)
			if err != nil {
				t.Fatal(err)
			}
			client, err := New("test-key", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(bytes.NewReader(body)), Request: r}, nil
			})}))
			if err != nil {
				t.Fatal(err)
			}
			result, err := step.call(t, client)
			if err != nil {
				t.Fatal(err)
			}
			raw := reflect.ValueOf(result).Elem().FieldByName("BaseResponse").Interface().(BaseResponse).Raw
			if raw["future_field"] != "forward-compatible" {
				t.Fatal("unknown upstream field lost from Raw")
			}
		})
	}
}

func responseProbe(s *contractSchema, mode string) any {
	if s.Nullable && (mode == "nullable" || mode == "nullable-leaves" && s.Type != "object" && s.Type != "array") {
		return nil
	}
	if len(s.Enum) > 0 {
		return s.Enum[0]
	}
	switch s.Type {
	case "object":
		out := map[string]any{}
		for name, child := range s.Properties {
			out[name] = responseProbe(child, mode)
		}
		if mode == "empty-collections" && len(s.Properties) == 0 {
			return out
		}
		if len(s.AdditionalProperties) > 0 && string(s.AdditionalProperties) != "false" {
			if string(s.AdditionalProperties) == "true" {
				out["contract_key"] = map[string]any{"nested": []any{true, "value", 12.5}}
			} else {
				var child contractSchema
				if err := json.Unmarshal(s.AdditionalProperties, &child); err != nil {
					panic(err)
				}
				out["contract_key"] = responseProbe(&child, mode)
			}
		}
		return out
	case "array":
		if mode == "empty-collections" {
			return []any{}
		}
		return []any{responseProbe(s.Items, mode)}
	case "boolean":
		return true
	case "integer":
		if s.Format == "int64" {
			return json.Number("9007199254740993")
		}
		return json.Number("17")
	case "number":
		return json.Number("12.3456789012345")
	case "string":
		if s.Format == "date-time" {
			return "2026-10-27T03:04:05.123456789Z"
		}
		if s.Format == "date" {
			return "2026-10-27"
		}
		if example, ok := s.Example.(string); ok {
			return example
		}
		return "contract-value"
	default:
		panic("unsupported response schema type: " + s.Type)
	}
}

func checkResponseType(t *testing.T, s *contractSchema, typ reflect.Type, path string) {
	t.Helper()
	if s.Nullable && typ.Kind() != reflect.Pointer && typ.Kind() != reflect.Map && typ.Kind() != reflect.Slice && typ.Kind() != reflect.Interface {
		t.Errorf("%s: nullable schema represented by non-nullable %s", path, typ)
	}
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	switch s.Type {
	case "object":
		if len(s.Properties) > 0 {
			if typ.Kind() != reflect.Struct {
				t.Errorf("%s: expected typed object, got %s", path, typ)
				return
			}
			fields := responseJSONFields(typ)
			for name, child := range s.Properties {
				fieldType, ok := fields[name]
				if !ok {
					t.Errorf("%s.%s: no typed field", path, name)
					continue
				}
				checkResponseType(t, child, fieldType, path+"."+name)
			}
		} else if typ.Kind() != reflect.Map || typ.Key().Kind() != reflect.String {
			t.Errorf("%s: expected string-keyed object map, got %s", path, typ)
		} else if len(s.AdditionalProperties) > 0 && string(s.AdditionalProperties) != "true" && string(s.AdditionalProperties) != "false" {
			var child contractSchema
			if err := json.Unmarshal(s.AdditionalProperties, &child); err != nil {
				t.Fatal(err)
			}
			checkResponseType(t, &child, typ.Elem(), path+".*")
		}
	case "array":
		if typ.Kind() != reflect.Slice {
			t.Errorf("%s: expected slice, got %s", path, typ)
			return
		}
		checkResponseType(t, s.Items, typ.Elem(), path+"[]")
	case "string":
		if s.Format == "date-time" {
			if typ != reflect.TypeFor[time.Time]() {
				t.Errorf("%s: date-time must use time.Time, got %s", path, typ)
			}
		} else if typ.Kind() != reflect.String && typ != reflect.TypeFor[time.Time]() {
			t.Errorf("%s: expected string, got %s", path, typ)
		}
	case "integer":
		if typ.Kind() != reflect.Int && typ.Kind() != reflect.Int64 {
			t.Errorf("%s: expected integer, got %s", path, typ)
		}
		if s.Format == "int64" && typ.Kind() != reflect.Int64 {
			t.Errorf("%s: int64 schema uses %s", path, typ)
		}
	case "number":
		if typ.Kind() != reflect.Float64 {
			t.Errorf("%s: expected float64, got %s", path, typ)
		}
	case "boolean":
		if typ.Kind() != reflect.Bool {
			t.Errorf("%s: expected bool, got %s", path, typ)
		}
	default:
		t.Fatalf("%s: unsupported schema type %q", path, s.Type)
	}
}

func responseJSONFields(typ reflect.Type) map[string]reflect.Type {
	fields := map[string]reflect.Type{}
	for _, field := range reflect.VisibleFields(typ) {
		if !field.IsExported() || field.Anonymous {
			continue
		}
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		fields[name] = field.Type
	}
	return fields
}

func decodeContractJSON(t *testing.T, b []byte) any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		t.Fatal(err)
	}
	return value
}

func checkFixtureContract(t *testing.T, s *contractSchema, value any, path string) {
	t.Helper()
	if value == nil {
		if !s.Nullable {
			t.Errorf("%s: null is not allowed by the response schema", path)
		}
		return
	}
	if len(s.Enum) > 0 {
		found := false
		for _, allowed := range s.Enum {
			if fmt.Sprint(value) == fmt.Sprint(allowed) {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: fixture value %v is not in enum %v", path, value, s.Enum)
		}
	}
	switch s.Type {
	case "object":
		object, ok := value.(map[string]any)
		if !ok {
			t.Errorf("%s: fixture is not an object", path)
			return
		}
		for _, required := range s.Required {
			if _, ok := object[required]; !ok {
				t.Errorf("%s.%s: required fixture field is missing", path, required)
			}
		}
		for name, child := range s.Properties {
			if item, exists := object[name]; exists {
				checkFixtureContract(t, child, item, path+"."+name)
			}
		}
		if len(s.AdditionalProperties) > 0 && string(s.AdditionalProperties) != "true" && string(s.AdditionalProperties) != "false" {
			var child contractSchema
			if err := json.Unmarshal(s.AdditionalProperties, &child); err != nil {
				t.Fatal(err)
			}
			for key, item := range object {
				if _, fixed := s.Properties[key]; !fixed {
					checkFixtureContract(t, &child, item, path+"."+key)
				}
			}
		}
	case "array":
		items, ok := value.([]any)
		if !ok {
			t.Errorf("%s: fixture is not an array", path)
			return
		}
		if len(items) < s.MinItems {
			t.Errorf("%s: fixture has fewer than %d items", path, s.MinItems)
		}
		for i, item := range items {
			checkFixtureContract(t, s.Items, item, fmt.Sprintf("%s[%d]", path, i))
		}
	case "integer", "number":
		number, ok := value.(json.Number)
		if !ok {
			t.Errorf("%s: fixture is not a number", path)
			return
		}
		rat, ok := new(big.Rat).SetString(string(number))
		if !ok || s.Type == "integer" && !rat.IsInt() {
			t.Errorf("%s: invalid fixture number %s", path, number)
			return
		}
		n, _ := number.Float64()
		if s.Minimum != nil && n < *s.Minimum || s.Maximum != nil && n > *s.Maximum {
			t.Errorf("%s: fixture number is outside its schema range", path)
		}
	case "string":
		text, ok := value.(string)
		if !ok {
			t.Errorf("%s: fixture is not a string", path)
			return
		}
		length := utf8.RuneCountInString(text)
		if length < s.MinLength || s.MaxLength > 0 && length > s.MaxLength {
			t.Errorf("%s: invalid fixture string length", path)
		}
		if s.Format == "date-time" {
			if _, err := time.Parse(time.RFC3339Nano, text); err != nil {
				t.Errorf("%s: invalid RFC 3339 fixture timestamp", path)
			}
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			t.Errorf("%s: fixture is not a boolean", path)
		}
	default:
		t.Fatalf("%s: unsupported fixture schema type %q", path, s.Type)
	}
}

func compareJSONLeaves(t *testing.T, expected, actual any, path string) {
	t.Helper()
	switch want := expected.(type) {
	case map[string]any:
		got, ok := actual.(map[string]any)
		if !ok {
			t.Errorf("%s: expected object, got %T", path, actual)
			return
		}
		keys := make([]string, 0, len(want))
		for key := range want {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value, exists := got[key]
			if !exists {
				t.Errorf("%s.%s: field lost during typed decode", path, key)
				continue
			}
			compareJSONLeaves(t, want[key], value, path+"."+key)
		}
	case []any:
		got, ok := actual.([]any)
		if !ok || len(got) != len(want) {
			t.Errorf("%s: expected %d array entries, got %v", path, len(want), actual)
			return
		}
		for i, value := range want {
			compareJSONLeaves(t, value, got[i], fmt.Sprintf("%s[%d]", path, i))
		}
	case json.Number:
		got, ok := actual.(json.Number)
		if !ok {
			t.Errorf("%s: expected number, got %T", path, actual)
			return
		}
		left, leftOK := new(big.Rat).SetString(string(want))
		right, rightOK := new(big.Rat).SetString(string(got))
		if !leftOK || !rightOK || left.Cmp(right) != 0 {
			t.Errorf("%s: number changed from %s to %s", path, want, got)
		}
	case string:
		got, ok := actual.(string)
		if !ok {
			t.Errorf("%s: expected string, got %T", path, actual)
			return
		}
		if before, err := time.Parse(time.RFC3339Nano, want); err == nil {
			after, err := time.Parse(time.RFC3339Nano, got)
			if err != nil || !before.Equal(after) {
				t.Errorf("%s: timestamp changed from %s to %s", path, want, got)
			}
		} else if want != got {
			t.Errorf("%s: string changed from %q to %q", path, want, got)
		}
	default:
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("%s: value changed from %v to %v", path, expected, actual)
		}
	}
}
