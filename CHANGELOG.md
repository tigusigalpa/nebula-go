# Changelog

All notable changes to this project will be documented in this file.
The format is based on Keep a Changelog, and this project adheres to Semantic Versioning.

## [0.1.0] - 2026-10-26

### Added
- Initial release covering all 41 Nebula Public API operations (OpenAPI 3.0.3, info.version 1.0.2).
- Typed request/response DTOs generated from the vendored spec, response metadata (rate limit, credits, request id, payment receipt), inspectable error hierarchy, retry policy with Retry-After/RateLimit-Reset aware backoff, X-Nebula-Locale/`locale` support.
- x402 opt-in mode: eligibility enforcement, offer parsing, spend policy checks.
