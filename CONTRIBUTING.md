# Contributing

1. Fork, branch, and open a pull request against `main`.
2. Keep public API stable within a major version; document breaking changes in CHANGELOG.md.
3. Run `go test ./...`, `go test -race ./...`, `go vet ./...` and `staticcheck ./...` before submitting. Use Go 1.24+ and a supported, patched Go release for security checks.
4. Do not commit API keys, wallet keys, or real account data.
5. `resources/openapi/nebula-openapi.json` is the contract source of truth. After an update, record its version and SHA-256 in `docs/UPSTREAM.md`, run `go generate ./...`, and update fixtures and the coverage matrix together. The generator rebuilds requests, responses, enums, service methods, validation rules and x402 eligibility. Keep Go naming overrides in `tools/nebula-gen/compat.json` and cross-parameter business rules in `request_validation.go`. See [Generation](docs/GENERATION.md) for supported schemas and the update workflow.
6. Run `go mod tidy` and `go mod verify`. Check vulnerabilities with `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` on a supported, fully patched Go toolchain.

All HTTP tests use local servers or in-memory transports. Do not add mandatory live API calls or wallet payments. Examples without an `Output` directive are compiled without execution.

When changing API parameters or response models, run `go test . -run 'Test(QueryParametersMatchContract|InvalidQueryParametersFailValidation|ChatBodyParametersMatchContract|FixtureDecodes)$'`. These tests read the pinned OpenAPI independently of generated validation rules, exercise parameters through public methods, and verify fixture values and every response property, including fields absent from the saved fixtures.

Run `go test ./tools/nebula-gen` when changing the generator. Its tests regenerate from an empty directory, check the released public API, and compile and exercise a separate SDK after adding parameters, nested response fields, an enum value and a request-body property to a temporary contract. `go run ./tools/nebula-gen -check` verifies every generated file without modifying it; CI also verifies that regeneration leaves the checkout unchanged.
