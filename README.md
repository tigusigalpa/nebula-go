# nebula-go

[![CI](https://github.com/tigusigalpa/nebula-go/actions/workflows/ci.yml/badge.svg)](https://github.com/tigusigalpa/nebula-go/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/tigusigalpa/nebula-go.svg)](https://pkg.go.dev/github.com/tigusigalpa/nebula-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Go SDK for the [Nebula Public API](https://docs.hiddensystems.ai/nebula/) — global markets intelligence across equities, commodities, indices, crypto, and AI. Covers **all 41 operations** with typed requests and responses. See [docs/API-COVERAGE.md](docs/API-COVERAGE.md).

> Market data returned by this SDK is **not investment advice**.

## Install

```bash
go get github.com/tigusigalpa/nebula-go
```

Requires Go 1.21+.

## Quick start

```go
client, err := nebula.New(os.Getenv("NEBULA_API_KEY"),
    nebula.WithTimeout(30*time.Second),
    nebula.WithRetry(3, 500*time.Millisecond, 5*time.Second),
)

ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

assets, err := client.Lookup.ListAssets(ctx, &nebula.ListAssetsRequest{Query: nebula.StringPtr("NVIDIA")})
asset := nebula.AssetRef{ID: assets.Assets[0].AssetID, Class: assets.Assets[0].AssetClass}
sentiment, err := client.Social.Sentiment(ctx, &nebula.SentimentRequest{Asset: &asset, Hours: nebula.IntPtr(168)})
fmt.Println(*sentiment.Metadata.RateLimitRemaining)
```

The API key must live server-side (env var, secret manager): Nebula sends no CORS headers, so never embed it in a browser app.

## Services

`Discovery`, `Lookup`, `Social`, `Forecasts`, `MarketSentiment`, `MarketData`, `Intelligence`, `Clusters`, `Agent`, `Account`, `Predictions` — one method per operation; see the coverage matrix.

## Errors and metadata

```go
res, err := client.Account.Get(ctx)
if errors.Is(err, nebula.ErrOutOfCredits) { /* 402 */ }
var ae *nebula.APIError
errors.As(err, &ae) // StatusCode, Message, Headers, PlanRenewsAt, PurchasedCredits
```

Sentinels: `ErrValidation`, `ErrTransport`, `ErrBadRequest`, `ErrUnauthorized`, `ErrOutOfCredits`, `ErrForbidden`, `ErrNotFound`, `ErrMethodNotAllowed`, `ErrRateLimited`, `ErrServer`, `ErrServiceUnavailable`, `ErrNotX402Eligible`, `ErrPaymentUncertain`.

Every response embeds `Metadata` (`RateLimit*`, `Credits*`, `PaymentReceipt`) and `Raw`.

## Agent streaming

```go
stream, _ := client.Agent.StreamChatCompletion(ctx, &nebula.ChatCompletionRequest{
    Messages: []nebula.ChatMessage{{Role: nebula.ChatRoleUser, Content: "NVIDIA outlook?"}},
})
defer stream.Close()
for chunk := range stream.Events() {
    fmt.Print(chunk.Choices[0].Delta.Content)
}
err = stream.Err()
```

## x402 (opt-in, spends real USDC)

> **Warning:** x402 requests pay USDC on Base mainnet (`eip155:8453`) and are **irreversible**. 37 GET endpoints are eligible; `/account`, forecasts and chat are rejected locally.

```go
import "github.com/tigusigalpa/nebula-go/x402"

signer, _ := x402.NewSignerFromPrivateKeyHex(os.Getenv("X402_KEY_HEX")) // never hardcode
policy := x402.SpendPolicy{MaxUSDCPerPayment: "0.05"}
client, _ := x402.NewClient(signer, policy)
// preview without paying:
offer, _ := x402.PreviewOffer(ctx, nil, "GET", nebula.X402BaseURL+"/assets?query=NVDA")
```

Uses the official `github.com/x402-foundation/x402/go/v2` client with spend controls enabled. Mock tests cover offer → signature → receipt without a wallet or network.

## Docs

- [API coverage matrix](docs/API-COVERAGE.md) (41/41)
- [Upstream contract](docs/UPSTREAM.md) — OpenAPI 3.0.3, version 1.0.2, SHA-256 pinned
- [Русская версия](README.ru.md)

## License

MIT © [Igor Sazonov](https://github.com/tigusigalpa)
