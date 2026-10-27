// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// Emotions - Emotion distribution. GET /emotions
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-emotions/
func (s *SocialService) Emotions(ctx context.Context, req *EmotionsRequest) (*EmotionsResponse, error) {
	return get[EmotionsResponse](ctx, s.client, "/emotions", req, "get_emotions")
}

// Mindshare - Mindshare. GET /mindshare
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-mindshare/
func (s *SocialService) Mindshare(ctx context.Context, req *MindshareRequest) (*MindshareResponse, error) {
	return get[MindshareResponse](ctx, s.client, "/mindshare", req, "get_mindshare")
}

// RelatedSubjects - Related subjects. GET /related-subjects
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-related-subjects/
func (s *SocialService) RelatedSubjects(ctx context.Context, req *RelatedSubjectsRequest) (*RelatedSubjectsResponse, error) {
	return get[RelatedSubjectsResponse](ctx, s.client, "/related-subjects", req, "get_related_subjects")
}

// Sentiment - Sentiment. GET /sentiment
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-sentiment/
func (s *SocialService) Sentiment(ctx context.Context, req *SentimentRequest) (*SentimentResponse, error) {
	return get[SentimentResponse](ctx, s.client, "/sentiment", req, "get_sentiment")
}

// SentimentTimeframes - Sentiment by timeframe. GET /sentiment/timeframes
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-sentiment-timeframes/
func (s *SocialService) SentimentTimeframes(ctx context.Context, req *SentimentTimeframesRequest) (*SentimentTimeframesResponse, error) {
	return get[SentimentTimeframesResponse](ctx, s.client, "/sentiment/timeframes", req, "get_sentiment_timeframes")
}

// SocialMomentum - Social momentum. GET /social-momentum
// Docs: https://docs.hiddensystems.ai/nebula/reference/get-social-momentum/
func (s *SocialService) SocialMomentum(ctx context.Context, req *SocialMomentumRequest) (*SocialMomentumResponse, error) {
	return get[SocialMomentumResponse](ctx, s.client, "/social-momentum", req, "get_social_momentum")
}
