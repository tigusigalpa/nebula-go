# Generating the SDK

The source contract is [the vendored OpenAPI JSON](../resources/openapi/nebula-openapi.json). Generation uses the Go standard library and never downloads the contract or calls Nebula. It does not need an API key or wallet.

From the repository root:

```sh
go generate ./...
go run ./tools/nebula-gen -check
go test ./...
```

The first command writes generated files. The second checks them without writing: missing, edited and obsolete files are errors. Unchanged files keep their contents and modification time, including a checkout with CRLF line endings. When a whole generated file becomes obsolete, regeneration removes it only if its header identifies it as owned by this generator.

## What comes from the contract

| Output | Contents |
| --- | --- |
| `request_*.go` | Query/body structs, parameter encoding, Agent body validation |
| `models_*.go` | Response structs, nested objects and typed dynamic maps |
| `enums_generated.go` | Shared string enum types and constants |
| `service_*.go` | Public service methods, paths and operation IDs |
| `validation_rules.go` | Query required fields, enum values, bounds and lengths |
| `x402eligible.go` | Eligibility from `x-x402-eligible` |

The generator can reconstruct these files without reading an existing generated file:

```sh
go run ./tools/nebula-gen -out ./generated-preview
```

That directory contains the generated portion of the SDK. A standalone client also needs the handwritten transport, client configuration, metadata, errors, helpers and business rules. Remove a preview directory after inspecting it; generated files there use the same package as the SDK.

Available flags are `-spec`, `-compat`, `-out` and `-check`. Paths are relative to the working directory. The defaults assume the repository root.

## Preserve public Go names

[compat.json](../tools/nebula-gen/compat.json) records naming choices that cannot be inferred reliably from OpenAPI: service/method names, shared response types, nested type names, enum aliases and constant spellings, field-name exceptions and the order of public request fields. It contains no response fields, parameter constraints or generated Go templates.

For existing operations, keep those names stable. For a new operation, add an operation entry to choose its public service, method and request/response names. When introducing a new service group, also register the group on `Client`. New fields within known schemas receive deterministic names automatically; add an override when a different public name is needed.

The generator verifies that schemas sharing a Go type remain compatible. If upstream changes one of them independently, give it a separate name before regenerating. The released v0.2.0 type layouts, field order, JSON tags and enum constants are checked by a compatibility test. The six empty legacy helper types remain in `model_compat.go` for source compatibility; actual dynamic response fields use typed maps.

## Types and validation

- Objects become structs; arrays become slices; `additionalProperties` becomes a typed map. A free-form dynamic object becomes `map[string]any`.
- Object properties use pointers, preserving the SDK's released layouts. Nullable or optional scalar properties also use pointers. Slices and maps already represent a missing or null value with `nil`.
- Integer `format: int64` becomes `int64`; other integers use `int`. Numbers use `float64`; `date-time` uses `time.Time`.
- Query arrays follow `style: form` and `explode`: repeated values by default, CSV when `explode: false`. Dates are sent in UTC with fractional seconds preserved.
- `asset_id`/`asset_class` and `subject_id` retain the SDK's `AssetRef` and `SubjectRef` selectors. Standalone class filters are enabled by an explicit naming/selector override.
- Query constraints and Agent body constraints are generated from schema values. The Agent methods control the wire `stream` flag.

[request_validation.go](../request_validation.go) holds prose rules involving several parameters: mutually exclusive selectors, posting-community restrictions, and minimum/maximum filter relationships. These cannot be expressed by the current OpenAPI parameter schemas. Transport behavior, payment policy, forecast helpers and locale/cluster constants also remain handwritten.

The supported input is OpenAPI **3.0.x JSON**, covering the complete pinned Nebula contract. Local schema, parameter, request-body and response references are supported, including recursive response objects. Query parameters can be inherited from a path item. Success responses must provide a `200` JSON object; the supported POST adapter is Nebula Agent chat.

External references, path-item references, other HTTP methods, schema composition (`allOf`/`anyOf`/`oneOf`), mixed fixed/dynamic objects and unimplemented constraint keywords are rejected with an error. Recursive request-body validation is also rejected. Extend the generator and its tests before adopting such a contract change; it does not silently turn unsupported schemas into untyped values. All output is rendered and formatted before any file is changed.

## Update the contract

1. Replace the local contract and record its version, download date and SHA-256 in [UPSTREAM.md](UPSTREAM.md).
2. Review added operations and naming overrides, then run `go generate ./...`.
3. Review the diff. Update fixtures, the coverage matrix and any changed prose rules. Treat public type changes as API changes.
4. Run `go run ./tools/nebula-gen -check`, `go test ./...`, `go test -race ./...` and `go vet ./...`.

Generator tests cover a clean output directory, repeatable generation, read-only checks, obsolete output cleanup and public API compatibility. A schema-evolution test adds parameters, a nested response, an enum value and a nested Agent body field to a temporary contract, then builds and tests the generated SDK in an empty module with network access disabled. Existing independent contract tests continue to check all 41 operations against the pinned specification.
