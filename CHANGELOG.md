# Changelog

All notable changes to this project will be documented in this file.
The format is based on Keep a Changelog, and this project adheres to Semantic Versioning.

## [Unreleased]

## [0.3.0] - 2026-11-15

### Added
- Restore the complete offline OpenAPI generator for requests, typed responses, enums, service methods, query validation and x402 eligibility, with read-only drift checking.
- Generator tests for clean reconstruction, deterministic output, v0.2.0 public API compatibility, local references and schema evolution in a separately compiled SDK.

### Changed
- Keep public naming overrides and cross-parameter business rules outside generated files; CI now checks the complete generated SDK.

## [0.2.0] - 2026-11-15

### Added
- OpenAPI-derived query validation and coverage guard tying all 41 operations to service calls, docs, fixtures and x402 eligibility.
- `HTTPMetadata()` for consistent access even when a response body has its own metadata; `IsForecastable()` on both forecast results.
- Standalone asset-class filters for asset lookup, prediction markets, sentiment and social momentum.
- Contract tests exercising all 173 query parameters and the chat body through public methods, including enums, boundaries, required fields and array encoding.
- Typed response checks for all 41 operations: every fixture value, all 750 schema properties, nullability, empty collections, precise timestamps and large int64 values.
- Compilable examples and expanded English/Russian READMEs.

### Fixed
- Preserve fractional seconds in calendar-event and insight time filters.
- Keep request context alive through body decoding and preserve cancellation while waiting for retry backoff.
- Make retry jitter concurrency-safe, bounded and overflow-safe; parse HTTP-date Retry-After before considering rate-limit reset.
- Apply per-request locale and locale transport consistently to JSON and SSE calls.
- Reject malformed/empty successful JSON; recognize SSE error frames, interrupted streams and cancellation with a full event buffer.
- Block cross-origin credential redirects and validate x402 resource URLs structurally against the requested resource.
- Disable SDK retries and the standard transport's implicit replay for signed x402 GETs; retain uncertain payment identity and receipts.
- Upgrade Gorilla WebSocket to v1.5.3 to address GO-2026-6278.

### Changed
- x402 payment clients now require an explicit, nonempty `AllowedRecipients` policy; exact-scheme validation is unconditional. This intentionally tightens previously permissive payment behavior.
- Declare Go 1.24 as the minimum required by dependencies, replacing the unnecessarily restrictive Go 1.26.4 directive and inaccurate Go 1.21 README claim. CI checks Go 1.24 and the current stable release.
- Streams now use the configured overall timeout; configure a longer timeout for long Agent answers. An EOF before `[DONE]` is reported as an error.

## [0.1.0] - 2026-10-26

### Added
- Initial release covering all 41 Nebula Public API operations (OpenAPI 3.0.3, info.version 1.0.2).
- Typed request/response DTOs generated from the vendored spec, response metadata (rate limit, credits, request id, payment receipt), inspectable error hierarchy, retry policy with Retry-After/RateLimit-Reset aware backoff, X-Nebula-Locale/`locale` support.
- x402 opt-in mode: eligibility enforcement, offer parsing, spend policy checks.
