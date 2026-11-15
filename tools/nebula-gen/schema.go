package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type schema struct {
	Ref                  string             `json:"$ref"`
	Type                 string             `json:"type"`
	Format               string             `json:"format"`
	Description          string             `json:"description"`
	Nullable             bool               `json:"nullable"`
	Properties           map[string]*schema `json:"properties"`
	Items                *schema            `json:"items"`
	AdditionalProperties json.RawMessage    `json:"additionalProperties"`
	Required             []string           `json:"required"`
	Enum                 []any              `json:"enum"`
	Minimum              *float64           `json:"minimum"`
	Maximum              *float64           `json:"maximum"`
	MinLength            *int               `json:"minLength"`
	MaxLength            *int               `json:"maxLength"`
	MinItems             *int               `json:"minItems"`
	MaxItems             *int               `json:"maxItems"`
}

func (s *schema) UnmarshalJSON(data []byte) error {
	type plain schema
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	// Fail explicitly instead of silently weakening a future contract.
	for _, key := range []string{"allOf", "anyOf", "oneOf", "not", "discriminator", "pattern", "uniqueItems", "exclusiveMinimum", "exclusiveMaximum", "multipleOf", "minProperties", "maxProperties"} {
		if _, ok := fields[key]; ok {
			return fmt.Errorf("unsupported schema keyword %q", key)
		}
	}
	return json.Unmarshal(data, (*plain)(s))
}

type parameter struct {
	Ref         string  `json:"$ref"`
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Style       string  `json:"style"`
	Explode     *bool   `json:"explode"`
	Schema      *schema `json:"schema"`
}

type media struct {
	Schema *schema `json:"schema"`
}

type body struct {
	Ref     string           `json:"$ref"`
	Content map[string]media `json:"content"`
}

type operation struct {
	ID          string          `json:"operationId"`
	Summary     string          `json:"summary"`
	Tags        []string        `json:"tags"`
	Parameters  []parameter     `json:"parameters"`
	RequestBody body            `json:"requestBody"`
	Responses   map[string]body `json:"responses"`
	External    struct {
		URL string `json:"url"`
	} `json:"externalDocs"`
	X402         bool `json:"x-x402-eligible"`
	path, method string
}

type document struct {
	OpenAPI string
	Version string
	ops     []*operation
	raw     json.RawMessage
}

func parseDocument(data []byte) (*document, error) {
	var root struct {
		OpenAPI string `json:"openapi"`
		Info    struct {
			Version string `json:"version"`
		} `json:"info"`
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(root.OpenAPI, "3.0.") {
		return nil, fmt.Errorf("unsupported OpenAPI version %q (expected 3.0.x)", root.OpenAPI)
	}
	d := &document{OpenAPI: root.OpenAPI, Version: root.Info.Version, raw: data}
	seen := map[string]bool{}
	for _, path := range keys(root.Paths) {
		item := root.Paths[path]
		if _, ok := item["$ref"]; ok {
			return nil, fmt.Errorf("path-item references are unsupported: %s", path)
		}
		var inherited []parameter
		if raw, ok := item["parameters"]; ok {
			if err := json.Unmarshal(raw, &inherited); err != nil {
				return nil, err
			}
		}
		for _, method := range keys(item) {
			switch method {
			case "put", "patch", "delete", "head", "options", "trace":
				return nil, fmt.Errorf("unsupported method %s %s", method, path)
			}
			if method != "get" && method != "post" {
				continue
			}
			var op operation
			if err := json.Unmarshal(item[method], &op); err != nil {
				return nil, fmt.Errorf("%s %s: %w", method, path, err)
			}
			if op.ID == "" || seen[op.ID] {
				return nil, fmt.Errorf("missing or duplicate operationId %q", op.ID)
			}
			seen[op.ID] = true
			op.path, op.method = path, method
			// Operation parameters override path parameters by (name, in).
			params := map[string]parameter{}
			for _, p := range append(inherited, op.Parameters...) {
				for visited := map[string]bool{}; p.Ref != ""; {
					if visited[p.Ref] {
						return nil, fmt.Errorf("cyclic parameter reference %s", p.Ref)
					}
					visited[p.Ref] = true
					raw, err := d.lookup(p.Ref)
					if err != nil {
						return nil, err
					}
					p = parameter{}
					if err := json.Unmarshal(raw, &p); err != nil {
						return nil, err
					}
				}
				if p.In == "header" {
					continue
				} // Authentication/localization belong to the transport.
				if p.In != "query" {
					return nil, fmt.Errorf("%s: unsupported parameter location %q", op.ID, p.In)
				}
				if p.Schema == nil {
					return nil, fmt.Errorf("%s: parameter %s has no schema", op.ID, p.Name)
				}
				if p.Style != "" && p.Style != "form" {
					return nil, fmt.Errorf("%s: unsupported query style %q", op.ID, p.Style)
				}
				params[p.Name] = p
			}
			op.Parameters = nil
			for _, name := range keys(params) {
				op.Parameters = append(op.Parameters, params[name])
			}
			d.ops = append(d.ops, &op)
		}
	}
	if len(d.ops) == 0 {
		return nil, fmt.Errorf("contract contains no operations")
	}
	return d, nil
}

func (d *document) lookup(ref string) (json.RawMessage, error) {
	if !strings.HasPrefix(ref, "#/") {
		return nil, fmt.Errorf("only local JSON references are supported: %q", ref)
	}
	raw := d.raw
	for _, part := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(raw, &object); err != nil {
			return nil, fmt.Errorf("invalid reference %s: %w", ref, err)
		}
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		var ok bool
		raw, ok = object[part]
		if !ok {
			return nil, fmt.Errorf("unresolved reference %s", ref)
		}
	}
	return raw, nil
}

func (d *document) resolve(s *schema) (*schema, error) {
	seen := map[string]bool{}
	for s != nil && s.Ref != "" {
		if seen[s.Ref] {
			return nil, fmt.Errorf("cyclic schema alias %s", s.Ref)
		}
		seen[s.Ref] = true
		raw, err := d.lookup(s.Ref)
		if err != nil {
			return nil, err
		}
		s = &schema{}
		if err := json.Unmarshal(raw, s); err != nil {
			return nil, err
		}
	}
	if s == nil {
		return nil, fmt.Errorf("missing schema")
	}
	return s, nil
}

func (d *document) jsonSchema(b body) (*schema, error) {
	for seen := map[string]bool{}; b.Ref != ""; {
		if seen[b.Ref] {
			return nil, fmt.Errorf("cyclic body/response reference %s", b.Ref)
		}
		seen[b.Ref] = true
		raw, err := d.lookup(b.Ref)
		if err != nil {
			return nil, err
		}
		b = body{}
		if err := json.Unmarshal(raw, &b); err != nil {
			return nil, err
		}
	}
	m, ok := b.Content["application/json"]
	if !ok {
		return nil, fmt.Errorf("missing application/json schema")
	}
	return d.resolve(m.Schema)
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
