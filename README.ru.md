# nebula-go

[![CI](https://github.com/tigusigalpa/nebula-go/actions/workflows/ci.yml/badge.svg)](https://github.com/tigusigalpa/nebula-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tigusigalpa/nebula-go.svg)](https://pkg.go.dev/github.com/tigusigalpa/nebula-go)
[![Version](https://img.shields.io/github/v/tag/tigusigalpa/nebula-go?label=version)](https://github.com/tigusigalpa/nebula-go/tags)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**Go-клиент для [Nebula Public API](https://docs.hiddensystems.ai/nebula/).** Помогает найти актив, посмотреть, что о нём говорят, получить прогноз внимания и настроений, изучить рынки предсказаний или задать вопрос Nebula Agent.

Все **41 операция** доступны через типизированные запросы и ответы. Основной клиент использует стандартную библиотеку Go; пакет `x402` подключает официальный платёжный клиент. Один экземпляр можно использовать из нескольких горутин и настроить через свой `*http.Client`.

[English README](README.md) · [Все методы API](docs/API-COVERAGE.md) · [Сохранённый контракт](docs/UPSTREAM.md)

## Установка

```sh
go get github.com/tigusigalpa/nebula-go
```

Нужен **Go 1.24 или новее**, с учётом зависимостей x402 в модуле. Версия SDK — `0.2.0`. Для production используйте поддерживаемый Go с последними исправлениями безопасности.

Храните ключ на сервере: в `NEBULA_API_KEY` или хранилище секретов. Nebula не отдаёт CORS-заголовки; помещать ключ в браузерный код нельзя.

## От поиска актива к настроениям

Пример ищет NVIDIA, сохраняет идентификатор вместе с классом актива и запрашивает настроения за неделю. Ошибки, пустая выдача и отсутствие заголовков с кредитами обработаны явно.

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
        fmt.Println("Активы не найдены")
        return nil
    }
    first := assets.Assets[0]
    asset := nebula.AssetRef{ID: first.AssetID, Class: first.AssetClass}
    sentiment, err := client.Social.Sentiment(ctx, &nebula.SentimentRequest{
        Asset: &asset, Hours: nebula.IntPtr(168),
    })
    if err != nil {
        return err
    }
    fmt.Printf("Получено наблюдений: %d\n", len(sentiment.Series))
    if remaining := sentiment.HTTPMetadata().CreditsRemaining; remaining != nil {
        fmt.Printf("Кредитов осталось: %d\n", *remaining)
    }
    return nil
}
```

В остальных примерах `client` и `ctx` уже созданы, а код находится внутри функции, возвращающей `error`. Есть и [компилируемые примеры](example_test.go).

## Что доступно

| Сервис | Возможности |
| --- | --- |
| `Discovery` | Скринер, импульс и популярные обсуждения |
| `Lookup` | Поиск, каталог активов, идентификаторы людей и тем |
| `Social` | Настроения, эмоции, доля внимания, связанные темы и социальный импульс |
| `Forecasts` | Каталог прогнозируемых активов и прогнозы на 72 часа |
| `MarketSentiment` | Торговые идеи и доходность, conviction, fear/greed, уровни цен, объём обсуждений |
| `MarketData` | История цен |
| `Intelligence` | Календарь событий, инсайты и внимание к ним |
| `Clusters` | Профили, сравнения, пересечения, ротация и позиции сообществ |
| `Predictions` | Рынки предсказаний, контракты площадок и история исходов |
| `Agent` | Обычный и потоковый OpenAI-совместимый чат |
| `Account` | План, ограничения запросов и остаток кредитов |

В [матрице покрытия](docs/API-COVERAGE.md) есть методы, типы, официальная документация и тестовые ответы для всех операций.

## Настройка

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

По умолчанию используются `X-API-Key`, **общий таймаут 30 секунд** и **три попытки суммарно** для подходящих GET-запросов. Bearer-режим отправляет только `Authorization`. Клиент задаёт `Accept: application/json` и `User-Agent: nebula-go/0.2.0`.

`WithHTTPClient` подключает прокси, TLS-настройки или инструменты наблюдения; `WithBaseURL` — шлюз или тестовый сервер; `WithUserAgent` — идентификатор приложения. Настройки HTTP-клиента копируются при создании SDK. Общие транспорты и логгеры должны поддерживать конкурентное использование. Перенаправления на другой origin блокируются, чтобы не передать ключ постороннему серверу.

`WithTimeout` охватывает все попытки и чтение ответа, включая весь поток Agent. Более короткий дедлайн контекста или таймаут собственного HTTP-клиента имеет приоритет. Ноль сохраняет значение по умолчанию; отрицательное время отклоняется. Для длинного ответа Agent задайте больше времени.

### Язык и параметры

Глобальная локаль отправляется в `X-Nebula-Locale`. Для одного вызова переопределите её через контекст:

```go
localized := nebula.WithRequestLocale(ctx, nebula.LocaleJapanese)
localized = nebula.WithRequestLocaleMode(localized, nebula.LocaleModeQuery)
// Передайте localized первым аргументом нужного метода.
```

`WithLocaleMode(LocaleModeQuery)` включает параметр `locale` глобально. Локаль — расширяемый строковый тип. Есть константы для испанского, итальянского, турецкого, русского, корейского, японского и китайского языков.

Для актива используйте `AssetRef{ID, Class}`, для человека или темы — `SubjectRef{ID}`. Совмещать их нельзя. `ListAssetsRequest.AssetClass` фильтрует каталог или поиск, `PredictionsRequest.AssetClass` — рынки предсказаний; `SentimentRequest.AssetClass` и `SocialMomentumRequest.AssetClass` выбирают рынок целиком. Отдельный фильтр класса и `AssetRef` используются по отдельности.

Необязательные поля — указатели: `nil` оставляет выбор за сервером, `IntPtr(0)` и `BoolPtr(false)` явно отправляют значение. Есть `StringPtr`, `Int64Ptr`, `Float64Ptr`, `TimePtr` и помощники для констант. Некорректные селекторы, обязательные поля, enum и числовые диапазоны возвращают `ErrValidation` до HTTP-вызова.

У списков есть `Limit`, `Offset` и `HasMore`. Сохраняйте фильтры, увеличивайте смещение на размер страницы и заканчивайте при `HasMore == false` или пустой выдаче. Фоновой загрузки нет.

`Clusters: []string{"kols", "media"}` превращается в повторяющиеся `clusters`. Фильтры календаря и инсайтов, а также списки классов активов в cluster-методах кодируются через запятую. SDK делает это сам.

## Квоты и метаданные

```go
account, err := client.Account.Get(ctx)
if err != nil {
    return err
}
fmt.Println("План:", account.Plan)
if account.Credits != nil {
    fmt.Println("Доступно кредитов:", account.Credits.TotalRemaining)
}
fmt.Println("Request ID:", account.HTTPMetadata().RequestID)
```

`HTTPMetadata()` есть у каждого ответа: статус, заголовки, request ID, ограничения запросов, кредиты, подсказка задержки и платёжная квитанция. Проверяйте необязательные значения на `nil`. `RateLimitReset` — секунды до сброса окна, а не Unix timestamp. `Retry-After` понимает секунды и HTTP-дату.

В некоторых телах ответа собственное поле `Metadata` описывает расчёт. `HTTPMetadata()` всегда возвращает HTTP-контекст; доступ через `BaseResponse.Metadata` сохранён. `Raw` хранит JSON-объект, включая ещё неизвестные SDK поля. Перед логированием raw-данных и заголовков учитывайте их содержимое.

## Ошибки и повторы

```go
account, err := client.Account.Get(ctx)
if err != nil {
    switch {
    case errors.Is(err, nebula.ErrUnauthorized):
        return fmt.Errorf("проверьте NEBULA_API_KEY: %w", err)
    case errors.Is(err, nebula.ErrOutOfCredits):
        return fmt.Errorf("закончились кредиты: %w", err)
    case errors.Is(err, context.DeadlineExceeded):
        return fmt.Errorf("истекло время запроса: %w", err)
    }
    var apiErr *nebula.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d; задержка %s; request %s\n",
            apiErr.StatusCode, apiErr.RetryAfter, apiErr.Metadata.RequestID)
    }
    return err
}
fmt.Println(account.Plan)
```

Нужен стандартный пакет `errors`. `APIError` сохраняет статус, сообщение, заголовки, метаданные и до 8 КиБ тела ответа. Для 402 доступны `PlanRenewsAt` и `PurchasedCredits`, если сервер их прислал. `TransportError` оборачивает сетевую ошибку, `DecodeError` сообщает о некорректном JSON при успешном HTTP-статусе. Работают `errors.Is`, `errors.As` и проверка отмены контекста.

GET в режиме API-ключа повторяется при временных сетевых сбоях и статусах **429, 500, 502, 503, 504**. Приоритет задержки: `Retry-After`, `RateLimit-Reset`, экспоненциальная задержка со случайным разбросом. `MaxDelay` ограничивает вычисленную задержку вместе с разбросом. Серверная подсказка может быть больше, но общий дедлайн продолжает действовать. `WithRetry(1, 0, 0)` отключает повторы.

POST Agent, x402-запросы, остальные 4xx, ошибки декодирования и чтения тела SDK автоматически не повторяет.

## Прогноз может быть недоступен — это нормальный ответ

```go
forecast, err := client.Forecasts.Mindshare(ctx, &nebula.MindshareForecastRequest{
    Asset: &nebula.AssetRef{ID: "NVDA", Class: nebula.AssetClassEquity},
})
if err != nil {
    return err
}
if !forecast.IsForecastable() {
    if forecast.Reason != nil {
        fmt.Println("Прогноз недоступен:", *forecast.Reason)
    }
    return nil
}
fmt.Println("Точек прогноза:", len(forecast.Series))
```

У обоих прогнозов есть `Status`, `Reason` и `IsForecastable()`. `status=not_forecastable` не считается ошибкой. `Refresh: BoolPtr(true)` запрашивает пересчёт и может стоить больше кредитов. Сохранённый контракт описывает около 5 секунд на обновление и возврат кешированного прогноза после 20 секунд ожидания. Учитывайте это в таймауте и проверяйте `Refreshed`. Стоимость и ограничения — в [документации прогнозов](https://docs.hiddensystems.ai/nebula/reference/get-forecasts-mindshare/).

## Рынки предсказаний

```go
markets, err := client.Predictions.ListMarkets(ctx, &nebula.PredictionsRequest{
    Query: nebula.StringPtr("NVIDIA"), Limit: nebula.IntPtr(10),
    AssetClass: nebula.AssetClassPtr(nebula.AssetClassEquity),
})
if err != nil {
    return err
}
if len(markets.Markets) > 0 {
    market, err := client.Predictions.Market(ctx, &nebula.PredictionMarketRequest{
        MarketID: markets.Markets[0].ID,
    })
    if err != nil {
        return err
    }
    fmt.Println(market.Question)
}
```

Используйте идентификаторы из API. `Predictions.History` принимает ID контракта площадки и ID исхода. У каждой площадки сохраняются собственные котировки и правила разрешения контракта.

## Nebula Agent

Agent не хранит диалог: каждый раз передавайте полную историю. Роли: `system`, `developer`, `user`, `assistant`; `developer` трактуется как серверная системная инструкция. `Model` принят для совместимости с OpenAI и игнорируется Nebula. `ThinkingMode` принимает `auto` или `deep`.

```go
request := &nebula.ChatCompletionRequest{
    Messages: []nebula.ChatMessage{
        {Role: nebula.ChatRoleUser, Content: "Что изменилось в обсуждении NVIDIA?"},
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
    fmt.Println("Потрачено кредитов:", answer.Usage.CreditsUsed)
}
```

Стоимость зависит от использования; в сохранённом контракте минимум — 10 кредитов. Для потока передайте тот же запрос в `StreamChatCompletion`: метод сам выставит флаг.

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

При досрочном выходе закройте поток или отмените контекст. `Err()` сообщает о некорректных событиях, ошибках внутри потока, отмене и обрыве до `[DONE]`. HTTP-метаданные — в `stream.Metadata()`. Событие может содержать только роль, причину завершения или usage: проверяйте элементы `Choices` и указатели на текст. Метод определяет режим независимо от совместимого поля `Stream`.

## x402: сначала предложение, затем разрешение на оплату

> **x402 может необратимо потратить настоящие USDC в Base mainnet.** По умолчанию SDK работает с API-ключом. Для платёжного клиента нужны лимит на один платёж и список доверенных получателей.

Просмотр предложения не требует кошелька или приватного ключа и ничего не подписывает:

```go
// Дополнительные импорты: "net/http" и "github.com/tigusigalpa/nebula-go/x402".
offer, err := x402.PreviewOffer(ctx, nil, http.MethodGet,
    nebula.X402BaseURL+"/assets?query=NVIDIA")
if err != nil {
    return err
}
fmt.Println("Ресурс:", offer.ResourceURL)
fmt.Println("Сеть:", offer.Network, "получатель:", offer.PayTo)
fmt.Println("Атомарных единиц USDC:", offer.Amount)
```

Передавайте `nil` или обычный HTTP-клиент без оплаты: клиент с платёжным транспортом сам может разрешить расходование. Preview показывает первый вариант оплаты, полное предложение остаётся в `Raw`. Проверяйте получателей независимо перед добавлением в `SpendPolicy.AllowedRecipients`: предложение само по себе не доказывает, что адресу можно доверять.

Для оплаты передайте свой `x402.Signer` и `SpendPolicy` в `x402.NewClient`. `MaxUSDCPerPayment` — десятичная строка с точностью до шести знаков; `AllowedRecipients` должен содержать доверенные адреса. Это лимит **одного платежа**: несколько вызовов могут суммарно потратить больше. Ограничения официального клиента остаются включёнными.

До подписи проверяются x402 v2, схема exact, сеть Base, USDC, получатель, сумма и соответствие ресурса адресу запроса. Для production нужен HTTPS. Доступны 37 GET-операций на отдельном адресе; `/account`, `/forecasts/mindshare`, `/forecasts/sentiment` и чат отклоняются локально. `/forecasts/assets` доступен.

Потерянный ответ после подписи, ошибка чтения оплаченного ответа, перенаправление или серверная ошибка могут вернуть `ErrPaymentUncertain`. Через `*nebula.PaymentUncertainError` получите исходный URL, request ID и доступную `Metadata.PaymentReceipt`; сверяйте платёж до новой попытки. Для подписанных GET отключён и скрытый повтор транспорта Go при обрыве соединения. Тесты не проводят настоящие платежи.

## Разработка

```sh
go generate ./...
go test ./...
go test -race ./...
go vet ./...
staticcheck ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
go mod tidy
go mod verify
```

Тесты используют локальные HTTP-серверы и искусственные ответы; ключи и кошельки не нужны. Проверка покрытия связывает 41 операцию спецификации с сервисными вызовами, документацией, fixtures и x402-доступностью. Отдельные контрактные тесты сверяют все query-параметры, тело чата и каждое поле типизированных ответов с сохранённой OpenAPI, включая nullable-поля и точность чисел и времени. `go generate` обновляет только query-валидацию. Генератор исходных DTO не включён; модели сопровождаются вместе со спецификацией и тестовыми ответами.

[Как внести изменения](CONTRIBUTING.md) · [Безопасность](SECURITY.md) · [История изменений](CHANGELOG.md)

Рыночные данные, настроения, прогнозы и ответы Agent **не являются инвестиционной рекомендацией**. SDK не округляет и не интерпретирует значения API. Поля `share` — доли от 0 до 1, `_pct` — процентные пункты. Время хранится в `time.Time` без перевода в локальную зону компьютера.

## Лицензия

[MIT](LICENSE) © [Igor Sazonov](https://github.com/tigusigalpa) · sovletig@gmail.com
