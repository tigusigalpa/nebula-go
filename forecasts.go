package nebula

// IsForecastable reports whether the result contains a usable forecast.
// A not_forecastable result is a successful business response, not an API error.
func (r *MindshareForecastResponse) IsForecastable() bool {
	return r != nil && r.Status == ForecastStatusOK
}

// IsForecastable reports whether the result contains a usable forecast.
func (r *SentimentForecastResponse) IsForecastable() bool {
	return r != nil && r.Status == ForecastStatusOK
}
