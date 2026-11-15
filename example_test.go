package nebula_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	nebula "github.com/tigusigalpa/nebula-go"
)

// Examples without an Output directive are compiled but never call the live API
// during go test. Set NEBULA_API_KEY when adapting them into your application.
func ExampleClient() {
	client, err := nebula.New(os.Getenv("NEBULA_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	assets, err := client.Lookup.ListAssets(ctx, &nebula.ListAssetsRequest{Query: nebula.StringPtr("NVIDIA")})
	if err != nil {
		log.Fatal(err)
	}
	if len(assets.Assets) == 0 {
		return
	}
	first := assets.Assets[0]
	asset := nebula.AssetRef{ID: first.AssetID, Class: first.AssetClass}
	sentiment, err := client.Social.Sentiment(ctx, &nebula.SentimentRequest{Asset: &asset, Hours: nebula.IntPtr(168)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(sentiment.Series))
	if remaining := sentiment.HTTPMetadata().CreditsRemaining; remaining != nil {
		fmt.Println(*remaining)
	}
}

func ExampleForecastsService_Mindshare() {
	client, err := nebula.New(os.Getenv("NEBULA_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	forecast, err := client.Forecasts.Mindshare(context.Background(), &nebula.MindshareForecastRequest{Asset: &nebula.AssetRef{ID: "NVDA", Class: nebula.AssetClassEquity}})
	if err != nil {
		log.Fatal(err)
	}
	if !forecast.IsForecastable() {
		if forecast.Reason != nil {
			fmt.Println(*forecast.Reason)
		}
		return
	}
	fmt.Println(len(forecast.Series))
}

func ExamplePredictionsService_Market() {
	client, err := nebula.New(os.Getenv("NEBULA_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	markets, err := client.Predictions.ListMarkets(ctx, &nebula.PredictionsRequest{Query: nebula.StringPtr("NVIDIA"), Limit: nebula.IntPtr(10)})
	if err != nil {
		log.Fatal(err)
	}
	if len(markets.Markets) == 0 {
		return
	}
	market, err := client.Predictions.Market(ctx, &nebula.PredictionMarketRequest{MarketID: markets.Markets[0].ID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(market.Question)
}

func ExampleAgentService_StreamChatCompletion() {
	client, err := nebula.New(os.Getenv("NEBULA_API_KEY"), nebula.WithTimeout(2*time.Minute))
	if err != nil {
		log.Fatal(err)
	}
	stream, err := client.Agent.StreamChatCompletion(context.Background(), &nebula.ChatCompletionRequest{Messages: []nebula.ChatMessage{{Role: nebula.ChatRoleUser, Content: "Summarize the recent NVIDIA conversation."}}})
	if err != nil {
		log.Fatal(err)
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
		log.Print(err)
	}
}
