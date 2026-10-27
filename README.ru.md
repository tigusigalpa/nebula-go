# nebula-go

Go SDK для [Nebula Public API](https://docs.hiddensystems.ai/nebula/) — глобальная разведка рынков: акции, сырьё, индексы, крипто и AI. Покрывает **все 41 операций** с типизированными запросами и ответами. См. [docs/API-COVERAGE.md](docs/API-COVERAGE.md).

> Данные рынка **не являются инвестиционной рекомендацией**.

## Установка

```bash
go get github.com/tigusigalpa/nebula-go
```

Go 1.21+.

## Быстрый старт

```go
client, _ := nebula.New(os.Getenv("NEBULA_API_KEY"))
assets, _ := client.Lookup.ListAssets(ctx, &nebula.ListAssetsRequest{Query: nebula.StringPtr("NVIDIA")})
asset := nebula.AssetRef{ID: "NVDA", Class: nebula.AssetClassEquity}
s, _ := client.Social.Sentiment(ctx, &nebula.SentimentRequest{Asset: &asset, Hours: nebula.IntPtr(168)})
fmt.Println(s.Metadata.RateLimitRemaining)
```

Ключ — только на сервере: API не отдаёт CORS-заголовки.

## Сервисы

`Discovery`, `Lookup`, `Social`, `Forecasts`, `MarketSentiment`, `MarketData`, `Intelligence`, `Clusters`, `Agent`, `Account`, `Predictions`.

## Ошибки и метаданные

`errors.Is` — сентинелы `ErrUnauthorized`, `ErrOutOfCredits`, `ErrRateLimited` и др.; `errors.As` — `*APIError` (StatusCode, Headers, PurchasedCredits, RetryAfter). Каждый ответ несёт `Metadata` и `Raw`.

## Agent streaming (SSE)

```go
stream, _ := client.Agent.StreamChatCompletion(ctx, req)
defer stream.Close()
for chunk := range stream.Events() { /* ... */ }
err := stream.Err()
```

## x402 (opt-in)

**Внимание:** запросы x402 тратят настоящие USDC в Base mainnet и необратимы. Адаптер `x402/` использует официальный `x402-foundation/x402/go/v2` со spend controls; 37 GET-операций доступны, остальные отклоняются локально.

## Документация

- [Матрица покрытия API](docs/API-COVERAGE.md) (41/41) · [Upstream-контракт](docs/UPSTREAM.md)
- [English README](README.md)

## Лицензия

MIT © [Igor Sazonov](https://github.com/tigusigalpa)
