package nebula

import (
	"strings"
	"time"
)

// StringPtr returns a pointer to v for optional request fields.
func StringPtr(v string) *string { return &v }

// IntPtr returns a pointer to v for optional request fields.
func IntPtr(v int) *int { return &v }

// Int64Ptr returns a pointer to v for optional request fields.
func Int64Ptr(v int64) *int64 { return &v }

// BoolPtr returns a pointer to v for optional request fields.
func BoolPtr(v bool) *bool { return &v }

// Float64Ptr returns a pointer to v for optional request fields.
func Float64Ptr(v float64) *float64 { return &v }

// TimePtr returns a pointer to v for optional request fields.
func TimePtr(v time.Time) *time.Time { return &v }

// AssetClassPtr returns a pointer to v for optional request fields.
func AssetClassPtr(v AssetClass) *AssetClass { return &v }

// KindPtr returns a pointer to v for optional request fields.
func KindPtr(v Kind) *Kind { return &v }

// GranularityPtr returns a pointer to v for optional request fields.
func GranularityPtr(v Granularity) *Granularity { return &v }

// ScreenerSortPtr returns a pointer to v for optional request fields.
func ScreenerSortPtr(v ScreenerSort) *ScreenerSort { return &v }

// SortOrderPtr returns a pointer to v for optional request fields.
func SortOrderPtr(v SortOrder) *SortOrder { return &v }

// InsightOrderByPtr returns a pointer to v for optional request fields.
func InsightOrderByPtr(v InsightOrderBy) *InsightOrderBy { return &v }

// PredictionProviderPtr returns a pointer to v for optional request fields.
func PredictionProviderPtr(v PredictionProvider) *PredictionProvider { return &v }

// PredictionStatusPtr returns a pointer to v for optional request fields.
func PredictionStatusPtr(v PredictionStatus) *PredictionStatus { return &v }

// ThinkingModePtr returns a pointer to v for optional request fields.
func ThinkingModePtr(v ThinkingMode) *ThinkingMode { return &v }

// csvJoin serializes a form/explode=false array parameter as CSV.
func csvJoin[S ~string](items []S) string {
	parts := make([]string, len(items))
	for i, s := range items {
		parts[i] = string(s)
	}
	return strings.Join(parts, ",")
}

// minLTEPtr enforces "each minimum must not exceed its maximum".
func minLTEPtr[T ~int | ~float64](lo, hi *T, op, loName, hiName string) error {
	if lo != nil && hi != nil && *lo > *hi {
		return &ValidationError{Op: op, Field: loName, Msg: loName + " must not exceed " + hiName}
	}
	return nil
}

func intIn(v int, allowed []int) bool {
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}
