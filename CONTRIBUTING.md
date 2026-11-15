# Contributing

1. Fork, branch, and open a pull request against `main`.
2. Keep public API stable within a major version; document breaking changes in CHANGELOG.md.
3. Run `go test ./...`, `go test -race ./...`, `go vet ./...` and `staticcheck ./...` before submitting. Use Go 1.24+ and a supported, patched Go release for security checks.
4. Do not commit API keys, wallet keys, or real account data.
5. `resources/openapi/nebula-openapi.json` is the contract source of truth. After an update, record its version and SHA-256 in `docs/UPSTREAM.md`, run `go generate ./...` for query-validation rules, and update request/response models, fixtures and the coverage matrix together. The original DTO generator is not included; its generated-file headers describe the original source, not an available regeneration command.
6. Run `go mod tidy` and `go mod verify`. Check vulnerabilities with `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` on a supported, fully patched Go toolchain.

All HTTP tests use local servers or in-memory transports. Do not add mandatory live API calls or wallet payments. Examples without an `Output` directive are compiled without execution.

When changing API parameters or response models, run `go test . -run 'Test(QueryParametersMatchContract|InvalidQueryParametersFailValidation|ChatBodyParametersMatchContract|FixtureDecodes)$'`. These tests read the pinned OpenAPI independently of generated validation rules, exercise parameters through public methods, and verify fixture values and every response property, including fields absent from the saved fixtures.
