# Nebula Golang Client/SDK/Library

![Nebula Golang Client/SDK/Library](https://i.postimg.cc/057CJ44Z/nebula-go-hero.jpg)

[![CI](https://github.com/tigusigalpa/nebula-go/actions/workflows/ci.yml/badge.svg)](https://github.com/tigusigalpa/nebula-go/actions/workflows/ci.yml)
[![Tests](https://img.shields.io/badge/tests-go%20test%20--race-brightgreen)](https://github.com/tigusigalpa/nebula-go/actions/workflows/ci.yml)
[![Go vet](https://img.shields.io/badge/code%20analysis-go%20vet-brightgreen)](https://github.com/tigusigalpa/nebula-go/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![CodeQL](https://github.com/tigusigalpa/nebula-go/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/nebula-go/actions/workflows/codeql.yml)
[![codecov](https://codecov.io/gh/tigusigalpa/nebula-go/graph/badge.svg)](https://codecov.io/gh/tigusigalpa/nebula-go)
[![Version](https://img.shields.io/github/v/tag/tigusigalpa/nebula-go?label=version)](https://github.com/tigusigalpa/nebula-go/tags)

**Explore market conversations with Go.** `nebula-go` is a client for the [Nebula Public API](https://docs.hiddensystems.ai/nebula/): find an asset, follow changes in sentiment and attention, compare communities, explore prediction markets, or ask Nebula Agent to put the data into words.

If you're building a research dashboard, a monitoring job or an assistant, you can start with a single API key. Requests are ordinary Go structs, responses have concrete field types, and every call takes a `context.Context`. The SDK handles authentication, query encoding, response decoding and eligible retries.

The client covers all **41 operations** in the [pinned OpenAPI 1.0.2 contract](docs/UPSTREAM.md). The core package uses the Go standard library. The optional `x402` package connects the official payment client for applications that choose to pay per request.

[Русская версия](README.ru.md) · [API coverage](docs/API-COVERAGE.md) · [Pinned contract](docs/UPSTREAM.md)

## Find your way around

- [Install and make your first request](#install)
- [Choose a service](#choose-a-service)
- [Configure the client](#configure-the-client)
- [Build requests and fetch another page](#build-requests)
- [Read response data](#read-response-data)
- [Check account limits and metadata](#account-limits-and-metadata)
- [Handle failures and retries](#failures-and-retries)
- [Work with forecasts](#forecasts) and [prediction markets](#prediction-markets)
- [Ask Nebula Agent or stream an answer](#nebula-agent)
- [Preview x402 payments](#x402-preview-before-authorizing-spending)
- [Troubleshoot an integration](#troubleshooting) or [work on the SDK](#development)

## Install

From your Go module:

```sh
go get github.com/tigusigalpa/nebula-go
```

Requires **Go 1.24 or newer**, including the module's x402 dependency graph. The SDK version declared in this source tree is `0.3.1`. The core package does not import the payment adapter; the adapter's dependencies are still declared in the same module.

Make your API key available to your application. For a local session, set an environment variable:

```sh
# macOS / Linux
export NEBULA_API_KEY='your-api-key'
```

```powershell
# Windows PowerShell
$env:NEBULA_API_KEY = 'your-api-key'
```

Keep the real key out of source control and browser code. For a web application, call Nebula from your Go backend and return the data your frontend needs. In deployment, a secret manager can supply the same value.

## Start with an asset

Let's answer a small question first: “What does the recent conversation around NVIDIA look like?” Save this program as `main.go` and run `go run .` after setting your key. It makes two account-mode API calls, so normal API credits apply.

The first call finds the asset. The second uses the returned identifier and asset class to fetch a week of sentiment observations.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    nebula "github.com/tigusigalpa/nebula-go"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    client, err := nebula.New(os.Getenv("NEBULA_API_KEY"))
    if err != nil {
        return err
    }
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    assets, err := client.Lookup.ListAssets(ctx, &nebula.ListAssetsRequest{
        Query: nebula.StringPtr("NVIDIA"),
    })
    if err != nil {
        return err
    }
    if len(assets.Assets) == 0 {
        fmt.Println("No matching assets")
        return nil
    }
    first := assets.Assets[0]
    fmt.Printf("Using %s (%s, %s)\n", first.Name, first.AssetID, first.AssetClass)
    asset := nebula.AssetRef{ID: first.AssetID, Class: first.AssetClass}
    sentiment, err := client.Social.Sentiment(ctx, &nebula.SentimentRequest{
        Asset: &asset, Hours: nebula.IntPtr(168),
    })
    if err != nil {
        return err
    }
    fmt.Printf("Received %d observations\n", len(sentiment.Series))
    if len(sentiment.Series) > 0 {
        latest := sentiment.Series[len(sentiment.Series)-1]
        fmt.Printf("Latest signal: %.2f at %s\n",
            latest.SentimentSignal, latest.Timestamp.Format(time.RFC3339))
    }
    if remaining := sentiment.HTTPMetadata().CreditsRemaining; remaining != nil {
        fmt.Printf("Credits left: %d\n", *remaining)
    }
    return nil
}
```

This example uses the first match to keep the program short. In an application, let the user choose a result or resolve the intended asset once and retain its identifier and class. A display symbol alone may not identify the asset you mean, especially for crypto.

The sentiment signal is a decay-weighted bullish-minus-bearish measure. Positive values lean bullish; negative values lean bearish, and magnitude also depends on activity. Read it alongside the timestamps and post counts.

Unless shown as a complete function, the remaining Go snippets belong inside a function returning `error` and assume an initialized `client` and `ctx`. Add standard-library imports as needed. There are also [compilable examples](example_test.go) in the repository.

## Choose a service

Services group related questions. Start with `Lookup` if you have a name but don't yet have an API identifier.

| Service on `client` | What you can ask for | A place to start |
| --- | --- | --- |
| `Discovery` | Screen assets and subjects; find momentum and trending conversations | `Screener`, `TopTrending` |
| `Lookup` | Search, list assets, resolve people and topics | `ListAssets`, `Search` |
| `Social` | Sentiment, emotions, mindshare, related subjects and social momentum | `Sentiment`, `Mindshare` |
| `Forecasts` | Forecastable assets and 72-hour mindshare/sentiment forecasts | `ListAssets`, `Mindshare` |
| `MarketSentiment` | Calls and returns, conviction, fear/greed, price levels and social volume | `Calls`, `FearGreedIndex` |
| `MarketData` | Price history | `PriceHistory` |
| `Intelligence` | Calendar events, insights and attention | `CalendarEvents`, `Insights` |
| `Clusters` | Community profiles, comparisons, overlap, rotation and stances | `List`, `Comparison` |
| `Predictions` | Prediction markets, venue listings and outcome history | `ListMarkets`, `Market` |
| `Agent` | Complete or stream an OpenAI-compatible chat request | `CreateChatCompletion` |
| `Account` | Plan, rate limits and credit balance | `Get` |

The [coverage matrix](docs/API-COVERAGE.md) links each operation to its method, types, official documentation and fixture.

## Configure the client

Create a client when your application starts and reuse it across goroutines. Give each call a context that reflects the lifetime of the work: for an HTTP handler, that will usually come from `r.Context()`.

```go
client, err := nebula.New(os.Getenv("NEBULA_API_KEY"),
    nebula.WithTimeout(45*time.Second),
    nebula.WithRetry(3, 500*time.Millisecond, 5*time.Second),
    nebula.WithLocale(nebula.LocaleRussian),
    nebula.WithAuthMode(nebula.AuthModeBearer),
)
if err != nil {
    return err
}
```

Use these options to adjust the defaults for your application:

| Setting | Default | Option |
| --- | --- | --- |
| API base | `nebula.BaseURL` (`/api/v1/public`) | `WithBaseURL` |
| Authentication | `X-API-Key` | `WithAuthMode` |
| Overall timeout per call | 30 seconds | `WithTimeout` |
| Eligible GET attempts | 3 total, including the first | `WithRetry` or `WithRetryPolicy` |
| Backoff | 250 ms initial delay, 5 s cap, with jitter | `WithRetry` or `WithRetryPolicy` |
| Locale | Unset; uses the server default | `WithLocale` |
| Locale transport | `X-Nebula-Locale` header | `WithLocaleMode` |
| User agent | `nebula-go/0.3.1` | `WithUserAgent` |
| Debug logging | Off | `WithLogger` |

Bearer mode sends `Authorization: Bearer …` instead of `X-API-Key`. JSON calls use `Accept: application/json`; Agent streams request `text/event-stream`.

`WithTimeout` covers all attempts, backoff and body reading for one call, including the full lifetime of a stream. A shorter context deadline or custom HTTP client timeout still wins. Zero keeps the default; negative durations are rejected. Raising only the context timeout will not extend the client's timeout, so configure both for a long Agent answer.

### Bring your own HTTP client

Use `WithHTTPClient` to plug in your proxy, TLS settings or instrumentation. For example, if your application already owns an HTTP client, you can pass it alongside a standard Go logger:

```go
httpClient := &http.Client{Timeout: 45 * time.Second}
client, err := nebula.New(os.Getenv("NEBULA_API_KEY"),
    nebula.WithHTTPClient(httpClient),
    nebula.WithTimeout(45*time.Second),
    nebula.WithUserAgent("market-monitor/1.0"),
    nebula.WithLogger(log.New(os.Stderr, "nebula ", log.LstdFlags)),
)
if err != nil {
    return err
}
fmt.Println("Client ready; x402 enabled:", client.X402())
```

The SDK copies HTTP client settings at construction; the transport and logger remain shared and must be safe for concurrent use. It blocks cross-origin redirects to protect credentials. `WithBaseURL` accepts the full API base, including any path prefix, so it also works with a gateway or local test server.

The debug logger records request methods and URLs with credential-like query fields redacted. Query text and asset identifiers can still appear, so choose log destinations and retention accordingly.

### Choose a language

Locale controls generated text where the API supports localization. The global locale uses `X-Nebula-Locale`. Override it for one call without changing the shared client:

```go
localized := nebula.WithRequestLocale(ctx, nebula.LocaleJapanese)
localized = nebula.WithRequestLocaleMode(localized, nebula.LocaleModeQuery)
answer, err := client.Agent.CreateChatCompletion(localized, &nebula.ChatCompletionRequest{
    Messages: []nebula.ChatMessage{
        {Role: nebula.ChatRoleUser, Content: "Summarize the recent NVIDIA conversation."},
    },
})
if err != nil {
    return err
}
for _, choice := range answer.Choices {
    if choice.Message != nil {
        fmt.Println(choice.Message.Content)
    }
}
```

`WithLocaleMode(LocaleModeQuery)` sets query transport globally. Locale is an extensible string type; constants cover Spanish, Italian, Turkish, Russian, Korean, Japanese and Chinese.

## Build requests

### Keep identifiers together

An asset and a subject are different kinds of lookup results. Use the selector that matches your result:

| You want to query… | Use… |
| --- | --- |
| One stock, commodity, index or crypto asset | `Asset: &nebula.AssetRef{ID: id, Class: class}` |
| One person, account or topic | `Subject: &nebula.SubjectRef{ID: id}` |
| All matching assets in one class | `ListAssetsRequest.AssetClass` |
| Prediction markets associated with one class | `PredictionsRequest.AssetClass` |
| Aggregate sentiment or momentum for one class | `SentimentRequest.AssetClass` or `SocialMomentumRequest.AssetClass` |

Keep the `AssetID` and `AssetClass` returned by lookup together. Classes are `AssetClassEquity`, `AssetClassCommodity`, `AssetClassIndex` and `AssetClassCrypto`. Use one selector at a time: an asset, a subject, or a standalone class filter where supported.

For a person or topic, search for subjects first and carry the returned ID into an analytics call:

```go
results, err := client.Lookup.Search(ctx, &nebula.SearchRequest{
    Query: "artificial intelligence",
    Kind:  nebula.KindPtr(nebula.KindSubject),
})
if err != nil {
    return err
}
for _, result := range results.Results {
    if result.Kind != nebula.KindSubject || result.SubjectID == nil {
        continue
    }
    sentiment, err := client.Social.Sentiment(ctx, &nebula.SentimentRequest{
        Subject: &nebula.SubjectRef{ID: *result.SubjectID},
        Hours:   nebula.IntPtr(24),
    })
    if err != nil {
        return err
    }
    fmt.Printf("%s: %d observations\n", result.Name, len(sentiment.Series))
    break // Just the first usable result for this example.
}
```

### Omitted values and explicit zeroes

Optional request fields use pointers so the API can distinguish “use your default” from “send this value”:

| Request value | What it means |
| --- | --- |
| `Limit: nil` | Omit the parameter; let the server choose its default |
| `Limit: nebula.IntPtr(25)` | Ask for 25 results |
| `Offset: nebula.IntPtr(0)` | Explicitly start at the first result |
| `Refresh: nebula.BoolPtr(false)` | Explicitly disable forecast refresh |

Helpers include `StringPtr`, `IntPtr`, `Int64Ptr`, `Float64Ptr`, `BoolPtr`, `TimePtr` and pointers to typed constants such as `AssetClassPtr`. Explicit values still have to satisfy the endpoint's rules: `Limit: IntPtr(0)` is invalid for asset listings.

Required fields, selector conflicts, documented enums and numeric ranges are checked before the HTTP call. A `ValidationError` identifies the operation and field. You can pass `nil` for an all-optional request, such as `Lookup.ListAssets(ctx, nil)`; operations requiring an asset or another identifier still require it.

### Pagination

For paginated lists such as assets, subjects and prediction markets, use `Limit`, `Offset` and the response's `HasMore`. The SDK fetches one page per call, which lets your application control how much work and credit usage to allow.

Here is a bounded walk through up to three pages of equities:

```go
const pageSize = 50
const maxPages = 3
for page := 0; page < maxPages; page++ {
    assets, err := client.Lookup.ListAssets(ctx, &nebula.ListAssetsRequest{
        AssetClass: nebula.AssetClassPtr(nebula.AssetClassEquity),
        Limit:      nebula.IntPtr(pageSize),
        Offset:     nebula.IntPtr(page * pageSize),
    })
    if err != nil {
        return err
    }
    for _, asset := range assets.Assets {
        fmt.Printf("%s: %s\n", asset.AssetID, asset.Name)
    }
    if !assets.HasMore || len(assets.Assets) == 0 {
        break
    }
}
```

Keep filters unchanged between pages. Page-size and offset limits differ by endpoint; use its request type and [API coverage entry](docs/API-COVERAGE.md) to find the relevant contract. A context deadline applies across the whole loop if you reuse the same context.

### Arrays and time windows

`Clusters: []string{"kols", "media"}` becomes repeated `clusters` parameters. Calendar filters, insight filters and cluster asset-class lists use CSV. The SDK handles encoding.

Set `Hours` for a lookback window and `Granularity` for the requested bucket size where the operation exposes them. The available combinations depend on the endpoint. Calendar and insight date filters accept `*time.Time`; the SDK serializes them in UTC and preserves fractional seconds.

## Read response data

Responses are typed all the way down to nested items and time-series points. Nullable scalar values use pointers, so you can distinguish an unavailable price from a real zero. For example:

```go
prices, err := client.MarketData.PriceHistory(ctx, &nebula.PriceHistoryRequest{
    Asset:       &nebula.AssetRef{ID: "NVDA", Class: nebula.AssetClassEquity},
    Hours:       nebula.IntPtr(24),
    Granularity: nebula.GranularityPtr(nebula.Granularity1h),
})
if err != nil {
    return err
}
for _, point := range prices.Series {
    timestamp := point.Timestamp.UTC().Format(time.RFC3339)
    if point.PriceUSD == nil {
        fmt.Println(timestamp, "price unavailable")
        continue
    }
    fmt.Printf("%s: $%.2f\n", timestamp, *point.PriceUSD)
}
```

An empty slice is a valid result; check its length before taking the first or last item. A missing cohort or timeframe key in a map can also be normal. Avoid turning unavailable data into zero just to fill a chart.

The SDK returns API values without rounding or rescaling them. Format them when presenting them:

- Shares and prediction prices use fractions from 0 to 1; `0.42` becomes `42%` when multiplied by 100 for display.
- Fields ending in `_pct` are already percentage points; `2.5` means `2.5%`.
- Sentiment post totals describe the trailing lookback. `IntervalPosts` and the other `Interval…` fields describe a single emitted bucket; summing overlapping trailing totals would count activity more than once.
- Timestamps decode into `time.Time`, preserving fractional seconds. Use `.UTC()` or `.In(location)` explicitly when displaying them.

Every top-level response also has `Raw`, a decoded JSON object that retains fields added upstream before the SDK gains a typed field. Prefer typed fields for normal application code: numbers in `Raw` use Go's usual `float64` JSON decoding, so large integers can lose precision there. The typed `int64` fields preserve it.

## Account limits and metadata

Use account data when you need the plan and balance. Use HTTP metadata when you want to understand the request you just made:

```go
account, err := client.Account.Get(ctx)
if err != nil {
    return err
}
fmt.Println("Plan:", account.Plan)
if account.Credits != nil {
    fmt.Println("Spendable credits:", account.Credits.TotalRemaining)
}
fmt.Println("Request ID:", account.HTTPMetadata().RequestID)
if charged := account.HTTPMetadata().CreditsCharged; charged != nil {
    fmt.Println("Credits charged for this call:", *charged)
}
```

Every response provides `HTTPMetadata()`: status, headers, request ID, optional rate-limit and credit figures, retry delay and any payment receipt. Check pointers before dereferencing them; an absent header does not mean the balance or limit is zero. `RateLimitReset` is seconds until reset. `Retry-After` accepts seconds or an HTTP date.

Some JSON bodies have their own `Metadata` field for calculation details. `HTTPMetadata()` consistently accesses HTTP information; `BaseResponse.Metadata` remains available too. The request ID is useful when investigating an unexpected response. Decide what is safe to log before recording raw bodies or headers.

## Failures and retries

Use `errors.Is` when deciding what your application should do and `errors.As` when you need details for diagnostics:

```go
account, err := client.Account.Get(ctx)
if err != nil {
    switch {
    case errors.Is(err, nebula.ErrValidation):
        var validationErr *nebula.ValidationError
        if errors.As(err, &validationErr) {
            return fmt.Errorf("fix %s: %s: %w", validationErr.Field, validationErr.Msg, err)
        }
    case errors.Is(err, nebula.ErrUnauthorized):
        return fmt.Errorf("check NEBULA_API_KEY: %w", err)
    case errors.Is(err, nebula.ErrOutOfCredits):
        return fmt.Errorf("account credits exhausted: %w", err)
    case errors.Is(err, context.DeadlineExceeded):
        return fmt.Errorf("request deadline reached: %w", err)
    case errors.Is(err, context.Canceled):
        return err
    }
    var apiErr *nebula.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d; retry hint %s; request %s\n",
            apiErr.StatusCode, apiErr.RetryAfter, apiErr.Metadata.RequestID)
    }
    return err
}
fmt.Println(account.Plan)
```

Import the standard `errors` package for this example. `APIError` retains status, message, headers, metadata and up to 8 KiB of response body. A credit-related 402 preserves `PlanRenewsAt` and `PurchasedCredits` when supplied. `TransportError` wraps network failures; `DecodeError` reports invalid successful JSON responses. Wrapped errors preserve context cancellation.

### What the SDK retries

Account-mode GETs retry on transient network errors and HTTP **429, 500, 502, 503 or 504**. Delay priority: `Retry-After`, `RateLimit-Reset`, then exponential backoff with jitter. `MaxDelay` caps computed backoff including jitter; server hints may exceed it, bounded by the overall deadline. `WithRetry(1, 0, 0)` disables retries.

POST chat requests, x402 requests, other 4xx errors, decoding failures and failed body reads are never automatically retried by the SDK. If a call returns an error after eligible retries, your application receives the final failure. Adding another retry loop outside the SDK multiplies attempts, so choose a retry budget deliberately.

## Forecasts

Start with `Forecasts.ListAssets` to discover forecast candidates. Then request a forecast for a selected asset and check its status before using the series:

```go
forecast, err := client.Forecasts.Mindshare(ctx, &nebula.MindshareForecastRequest{
    Asset: &nebula.AssetRef{ID: "NVDA", Class: nebula.AssetClassEquity},
})
if err != nil {
    return err
}
if !forecast.IsForecastable() {
    if forecast.Reason != nil {
        fmt.Println("Forecast unavailable:", *forecast.Reason)
    }
    return nil
}
fmt.Println("Forecast points:", len(forecast.Series))
```

Both forecasts expose `Status`, `Reason` and `IsForecastable()`. A successful HTTP call can return `status=not_forecastable`; show the reason or an unavailable state rather than treating it as a network failure.

`Refresh: BoolPtr(true)` requests recomputation and can cost more credits. The pinned contract describes about 5 seconds for refresh and a cached fallback after 20 seconds; allow time for that behavior and check `Refreshed`. Use response credit metadata to observe charges. The [forecast documentation](https://docs.hiddensystems.ai/nebula/reference/get-forecasts-mindshare/) describes the upstream pricing and limits.

## Prediction markets

Search for a question, open its canonical market, then use the venue listing and outcome IDs if you need history. These IDs refer to different objects, so carry each one from its own response field:

```go
markets, err := client.Predictions.ListMarkets(ctx, &nebula.PredictionsRequest{
    Query: nebula.StringPtr("NVIDIA"), Limit: nebula.IntPtr(10),
    AssetClass: nebula.AssetClassPtr(nebula.AssetClassEquity),
})
if err != nil {
    return err
}
if len(markets.Markets) == 0 {
    fmt.Println("No matching markets")
    return nil
}
market, err := client.Predictions.Market(ctx, &nebula.PredictionMarketRequest{
    MarketID: markets.Markets[0].ID,
})
if err != nil {
    return err
}
fmt.Println(market.Question)
```

`Predictions.History` takes both a listing ID and an outcome ID. Once you have a `market` from the detail call above, you can fetch the first listing's first outcome like this:

```go
if len(market.Listings) == 0 || len(market.Listings[0].Outcomes) == 0 {
    fmt.Println("No outcome history to request")
    return nil
}
listing := market.Listings[0]
outcome := listing.Outcomes[0]
history, err := client.Predictions.History(ctx, &nebula.PredictionHistoryRequest{
    ListingID: listing.ID,
    OutcomeID: outcome.ID,
    Hours:     nebula.IntPtr(24),
})
if err != nil {
    return err
}
fmt.Printf("%s: %d observations\n", outcome.Label, len(history.Series))
```

Use `Series` for new history integrations; `Points` is a deprecated compatibility field. Venue listings retain their own quotes and resolution rules. Check `ObservedAt` for quote freshness and `PriceBasis` for what a price represents. Compare volume only when `VolumeUnit` matches: one venue can report dollars while another reports contracts.

## Nebula Agent

Use the Agent when you want a written answer from Nebula. The service is stateless: your application owns the conversation and sends its full history each time. Start with a single user message:

```go
request := &nebula.ChatCompletionRequest{
    Messages: []nebula.ChatMessage{
        {Role: nebula.ChatRoleUser, Content: "Summarize the recent NVIDIA conversation."},
    },
}
answer, err := client.Agent.CreateChatCompletion(ctx, request)
if err != nil {
    return err
}
for _, choice := range answer.Choices {
    if choice.Message != nil {
        fmt.Println(choice.Message.Content)
    }
}
if answer.Usage != nil {
    fmt.Println("Credits used:", answer.Usage.CreditsUsed)
}
```

To continue a conversation, append the returned assistant text and the next user message to your application's message history before making another call. Roles are `system`, `developer`, `user` and `assistant`; `developer` is interpreted as a server-side system instruction. `Model` is accepted for OpenAI compatibility and ignored by Nebula. Set `ThinkingMode` with `ThinkingModePtr(ThinkingModeAuto)` or `ThinkingModePtr(ThinkingModeDeep)` when needed.

Billing varies by usage, with a 10-credit minimum in the pinned contract. Inspect `Usage` when present. Chat POSTs are not automatically retried.

### Stream an answer

For a UI that displays text as it arrives, pass the same `request` to `StreamChatCompletion`. Configure a timeout long enough for the answer before creating the client; the default 30 seconds covers the entire stream.

```go
stream, err := client.Agent.StreamChatCompletion(ctx, request)
if err != nil {
    return err
}
defer stream.Close()
for chunk := range stream.Events() {
    for _, choice := range chunk.Choices {
        if choice.Delta != nil && choice.Delta.Content != nil {
            fmt.Print(*choice.Delta.Content)
        }
    }
}
if err := stream.Err(); err != nil {
    return err
}
```

Keep the `defer stream.Close()` even if you expect to read the whole answer: it also releases resources when your code leaves early. Cancel the context when the caller no longer needs the answer.

A closed events channel alone does not tell you whether generation completed successfully. Always check `Err()` after the loop: it reports malformed events, API error frames, cancellation and an unexpected end before `[DONE]`. Partial text may already have been displayed, so keep that state separate from a completed answer in your UI.

`stream.Metadata()` exposes HTTP headers. Usage-only, role and finish chunks may contain no text: check choices and pointers as above. The selected method determines streaming mode regardless of the compatibility `Stream` field.

## x402: preview before authorizing spending

> **x402 can irreversibly spend real USDC on Base mainnet.** API-key mode is the default. The payment client is opt-in and requires a per-payment cap and trusted recipient allowlist.

If you use an API key, you can skip this section. The adapter is for applications that explicitly choose wallet-based payment through Nebula's separate x402 API.

This preview needs no wallet or private key and does not sign a payment:

```go
// Additional imports: "net/http" and "github.com/tigusigalpa/nebula-go/x402".
offer, err := x402.PreviewOffer(ctx, nil, http.MethodGet,
    nebula.X402BaseURL+"/assets?query=NVIDIA")
if err != nil {
    return err
}
fmt.Println("Resource:", offer.ResourceURL)
fmt.Println("Network:", offer.Network, "recipient:", offer.PayTo)
fmt.Println("Atomic USDC requested:", offer.Amount)
```

Pass `nil` or a plain, unpaid HTTP client to `PreviewOffer`; a payment-enabled client can itself authorize spending. Preview surfaces the first requirement and retains the full offer in `Raw`. It does not establish trust: verify recipients independently before adding them to `SpendPolicy.AllowedRecipients`.

To enable payment, supply your own `x402.Signer` and explicit `SpendPolicy` to `x402.NewClient`. `MaxUSDCPerPayment` is a decimal string with at most six decimal places; `AllowedRecipients` must be nonempty. This cap applies **per payment**, so repeated calls can spend more in total. The official client's spend controls remain enabled.

Before signing, the adapter checks x402 v2, the exact scheme, Base network, USDC asset, recipient, amount and requested resource URL. Production resources require HTTPS. All 37 eligible GETs use the separate x402 base; `/account`, `/forecasts/mindshare`, `/forecasts/sentiment` and chat are rejected locally. `/forecasts/assets` is eligible.

A lost signed response, failed signed body read, redirect or server error can return `ErrPaymentUncertain`. Inspect `*nebula.PaymentUncertainError` for the original URL, request ID and available `Metadata.PaymentReceipt`, then reconcile before trying again. Signed GETs also disable the standard Go transport's implicit replay after a connection failure. Tests never make live wallet payments.

## Troubleshooting

| What you see | What to check next |
| --- | --- |
| `New` returns “empty API key” | Check that the variable exists in the process running your app. A terminal, IDE and service can have different environments. |
| `ErrValidation` before a request | Inspect `ValidationError.Field` and `Msg`. Check required identifiers, paired asset ID/class, conflicting selectors and endpoint ranges. |
| `ErrUnauthorized` / HTTP 401 | Verify that the configured key is valid and that your gateway forwards the chosen authentication header. |
| `ErrOutOfCredits` / account HTTP 402 | Inspect the account balance and any renewal information on `APIError`. Immediate retries will not restore credits. |
| `ErrForbidden` / HTTP 403 | Read the API message. A requested history window may exceed the plan's allowance. |
| `ErrRateLimited` / HTTP 429 | Reduce request concurrency and inspect retry metadata. The SDK may already have used its retry attempts. |
| `context.DeadlineExceeded` | Check the SDK timeout, context deadline and custom HTTP client timeout. The earliest one ends the call. |
| Empty results or `nil` data fields | Confirm the selected ID/class and time window. Treat unavailable data explicitly in your application. |
| `not_forecastable` | Check `Reason` and use `Forecasts.ListAssets` to discover candidates. This is a normal forecast response. |
| A stream stops after showing some text | Read `stream.Err()` and retain the distinction between a partial and completed answer. |
| `ErrNotX402Eligible` | Use an account client for account details, the two forecast-generation endpoints and Agent chat. |
| `ErrPaymentUncertain` | Reconcile the original signed request and any receipt before authorizing another payment. |

When reporting a problem, include the SDK and Go versions, the method called, a minimal request with secrets removed, and the request ID if available. Use [GitHub issues](https://github.com/tigusigalpa/nebula-go/issues) for reproducible bugs and the [security policy](SECURITY.md) for sensitive reports.

## Development

The test suite runs without a Nebula key or wallet. It uses local HTTP servers, an in-memory transport and fixtures, so you can work on the SDK without making live API calls or payments.

Start with:

```sh
go test ./...
go vet ./...
```

For the broader checks used in development and CI:

```sh
go generate ./...
go test -race ./...
go run honnef.co/go/tools/cmd/staticcheck@latest ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
go mod tidy
go mod verify
```

The race detector needs a platform and C toolchain supported by Go. The analysis commands may download their tools; vulnerability checking also uses the vulnerability database. CI tests the Go 1.24 line and the stable release, and checks formatting, generated output and dependency consistency.

If you're looking for a particular part of the implementation:

| Location | What lives there |
| --- | --- |
| `service_*.go`, `request_*.go`, `models_*.go` | Public methods, request parameters and typed responses |
| [transport.go](transport.go), [stream.go](stream.go) | HTTP execution, retries, decoding and Agent streaming |
| [errors.go](errors.go), [metadata.go](metadata.go) | Error types and response headers |
| [x402/](x402/) | Opt-in payment adapter and its tests |
| [resources/openapi/nebula-openapi.json](resources/openapi/nebula-openapi.json) | Pinned upstream contract |
| [testdata/](testdata/) | Response fixtures |
| [request_contract_test.go](request_contract_test.go), [contract_test.go](contract_test.go) | Contract checks through public service methods |
| [tools/nebula-gen/](tools/nebula-gen/) | OpenAPI generator for requests, responses, enums, service methods and validation |

The coverage guard connects all **41 operations** to exercised service calls, docs, fixtures and x402 eligibility. Independent contract tests check **173 query parameters**, the chat body and **750 response properties**, including nullable fields, empty collections, large integers, fractional timestamps and unknown fields retained in `Raw`. These checks verify the pinned contract; they do not establish live API compatibility.

`go generate ./...` rebuilds the SDK surface from the pinned OpenAPI: request structs and serialization, typed responses, enums, service methods, query-validation rules and x402 eligibility. It works offline and can restore generated files into an empty directory. To check for missing, edited or obsolete generated files without writing anything, run `go run ./tools/nebula-gen -check`.

Public Go naming conventions live in `tools/nebula-gen/compat.json`; business rules involving several parameters live in `request_validation.go`. Keep changes to schema-derived fields in the OpenAPI rather than editing generated files. See the [generator guide](docs/GENERATION.md), [Contributing](CONTRIBUTING.md) and [Changelog](CHANGELOG.md).

[Contributing](CONTRIBUTING.md) · [Security](SECURITY.md) · [Changelog](CHANGELOG.md)

Market data, sentiment, forecasts and Agent output are **not investment advice**.

## License

[MIT](LICENSE) © [Igor Sazonov](https://github.com/tigusigalpa) · sovletig@gmail.com
