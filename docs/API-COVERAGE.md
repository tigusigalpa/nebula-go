# API Coverage — nebula-go

All 41 operations of the Nebula Public API (OpenAPI 1.0.2). Methods live on service fields of `*nebula.Client`. Every operation has a request type, response type, fixture and test.

| # | Endpoint | Method | Request | Response | Docs | x402 |
|---|---|---|---|---|---|---|
| 1 | `GET /screener` | ``client.Discovery.Screener(ctx, req)`` | `ScreenerRequest` | `ScreenerResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-screener/) | yes |
| 2 | `GET /top-momentum` | ``client.Discovery.TopMomentum(ctx, req)`` | `TopMomentumRequest` | `TopMomentumResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-top-momentum/) | yes |
| 3 | `GET /top-trending` | ``client.Discovery.TopTrending(ctx, req)`` | `TopTrendingRequest` | `TopTrendingResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-top-trending/) | yes |
| 4 | `GET /assets` | ``client.Lookup.ListAssets(ctx, req)`` | `ListAssetsRequest` | `ListAssetsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-assets/) | yes |
| 5 | `GET /search` | ``client.Lookup.Search(ctx, req)`` | `SearchRequest` | `SearchResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-search/) | yes |
| 6 | `GET /subjects` | ``client.Lookup.ListSubjects(ctx, req)`` | `ListSubjectsRequest` | `ListSubjectsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-subjects/) | yes |
| 7 | `GET /emotions` | ``client.Social.Emotions(ctx, req)`` | `EmotionsRequest` | `EmotionsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-emotions/) | yes |
| 8 | `GET /mindshare` | ``client.Social.Mindshare(ctx, req)`` | `MindshareRequest` | `MindshareResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-mindshare/) | yes |
| 9 | `GET /related-subjects` | ``client.Social.RelatedSubjects(ctx, req)`` | `RelatedSubjectsRequest` | `RelatedSubjectsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-related-subjects/) | yes |
| 10 | `GET /sentiment` | ``client.Social.Sentiment(ctx, req)`` | `SentimentRequest` | `SentimentResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-sentiment/) | yes |
| 11 | `GET /sentiment/timeframes` | ``client.Social.SentimentTimeframes(ctx, req)`` | `SentimentTimeframesRequest` | `SentimentTimeframesResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-sentiment-timeframes/) | yes |
| 12 | `GET /social-momentum` | ``client.Social.SocialMomentum(ctx, req)`` | `SocialMomentumRequest` | `SocialMomentumResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-social-momentum/) | yes |
| 13 | `GET /forecasts/assets` | ``client.Forecasts.ListAssets(ctx, req)`` | `ForecastAssetsRequest` | `ForecastAssetsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-forecasts-assets/) | yes |
| 14 | `GET /forecasts/mindshare` | ``client.Forecasts.Mindshare(ctx, req)`` | `MindshareForecastRequest` | `MindshareForecastResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-forecasts-mindshare/) | no |
| 15 | `GET /forecasts/sentiment` | ``client.Forecasts.Sentiment(ctx, req)`` | `SentimentForecastRequest` | `SentimentForecastResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-forecasts-sentiment/) | no |
| 16 | `GET /calls` | ``client.MarketSentiment.Calls(ctx, req)`` | `CallsRequest` | `CallsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-calls/) | yes |
| 17 | `GET /calls/returns` | ``client.MarketSentiment.CallReturns(ctx, req)`` | `CallReturnsRequest` | `CallReturnsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-calls-returns/) | yes |
| 18 | `GET /conviction-index` | ``client.MarketSentiment.ConvictionIndex(ctx, req)`` | `ConvictionIndexRequest` | `ConvictionIndexResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-conviction-index/) | yes |
| 19 | `GET /fear-greed-index` | ``client.MarketSentiment.FearGreedIndex(ctx, req)`` | `FearGreedIndexRequest` | `FearGreedIndexResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-fear-greed-index/) | yes |
| 20 | `GET /price-levels` | ``client.MarketSentiment.PriceLevels(ctx, req)`` | `PriceLevelsRequest` | `PriceLevelsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-price-levels/) | yes |
| 21 | `GET /social-volume-index` | ``client.MarketSentiment.SocialVolumeIndex(ctx, req)`` | `SocialVolumeIndexRequest` | `SocialVolumeIndexResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-social-volume-index/) | yes |
| 22 | `GET /social-volume-profile` | ``client.MarketSentiment.SocialVolumeProfile(ctx, req)`` | `SocialVolumeProfileRequest` | `SocialVolumeProfileResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-social-volume-profile/) | yes |
| 23 | `GET /price-history` | ``client.MarketData.PriceHistory(ctx, req)`` | `PriceHistoryRequest` | `PriceHistoryResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-price-history/) | yes |
| 24 | `GET /calendar-events` | ``client.Intelligence.CalendarEvents(ctx, req)`` | `CalendarEventsRequest` | `CalendarEventsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-calendar-events/) | yes |
| 25 | `GET /insights` | ``client.Intelligence.Insights(ctx, req)`` | `InsightsRequest` | `InsightsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-insights/) | yes |
| 26 | `GET /insights/attention` | ``client.Intelligence.InsightAttention(ctx, req)`` | `InsightAttentionRequest` | `InsightAttentionResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-insight-attention/) | yes |
| 27 | `GET /clusters` | ``client.Clusters.List(ctx, req)`` | `ClustersRequest` | `ClustersResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters/) | yes |
| 28 | `GET /clusters/comparison` | ``client.Clusters.Comparison(ctx, req)`` | `ClustersComparisonRequest` | `ClustersComparisonResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison/) | yes |
| 29 | `GET /clusters/comparison/divisive-assets` | ``client.Clusters.DivisiveAssets(ctx, req)`` | `DivisiveAssetsRequest` | `DivisiveAssetsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-divisive-assets/) | yes |
| 30 | `GET /clusters/comparison/hype-cycles` | ``client.Clusters.ComparisonHypeCycles(ctx, req)`` | `ComparisonHypeCyclesRequest` | `ComparisonHypeCyclesResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-hype-cycles/) | yes |
| 31 | `GET /clusters/comparison/traits` | ``client.Clusters.ComparisonTraits(ctx, req)`` | `ComparisonTraitsRequest` | `ComparisonTraitsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-comparison-traits/) | yes |
| 32 | `GET /clusters/hype-cycles` | ``client.Clusters.HypeCycles(ctx, req)`` | `ClusterHypeCyclesRequest` | `ClusterHypeCyclesResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-hype-cycles/) | yes |
| 33 | `GET /clusters/overlap` | ``client.Clusters.Overlap(ctx, req)`` | `ClusterOverlapRequest` | `ClusterOverlapResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-overlap/) | yes |
| 34 | `GET /clusters/profile` | ``client.Clusters.Profile(ctx, req)`` | `ClusterProfileRequest` | `ClusterProfileResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-profile/) | yes |
| 35 | `GET /clusters/rotation` | ``client.Clusters.Rotation(ctx, req)`` | `ClusterRotationRequest` | `ClusterRotationResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-rotation/) | yes |
| 36 | `GET /clusters/stances` | ``client.Clusters.Stances(ctx, req)`` | `ClusterStancesRequest` | `ClusterStancesResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-clusters-stances/) | yes |
| 37 | `POST /chat/completions` | ``client.Agent.CreateChatCompletion` + `StreamChatCompletion`` | `ChatCompletionRequest` | `ChatCompletionResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/post-chat-completions/) | no |
| 38 | `GET /account` | ``client.Account.Get(ctx, req)`` | - | `AccountResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-account/) | no |
| 39 | `GET /predictions` | ``client.Predictions.ListMarkets(ctx, req)`` | `PredictionsRequest` | `PredictionsResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-predictions/) | yes |
| 40 | `GET /predictions/history` | ``client.Predictions.History(ctx, req)`` | `PredictionHistoryRequest` | `PredictionHistoryResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-prediction-history/) | yes |
| 41 | `GET /predictions/market` | ``client.Predictions.Market(ctx, req)`` | `PredictionMarketRequest` | `PredictionMarketResponse` | [docs](https://docs.hiddensystems.ai/nebula/reference/get-prediction-market/) | yes |

Fixture per op: `testdata/<operationId>.json`. Coverage guard: `coverage_test.go` calls all 41 methods against httptest and asserts path, auth and decode; `X402Eligible` lists the 37 eligible ops.