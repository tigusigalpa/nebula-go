package nebula

//go:generate go run ./tools/nebula-gen

import (
	"math"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type queryRule struct {
	name                 string
	kind                 string
	required             bool
	csv                  bool
	enum                 []string
	min, max             *float64
	minLength, maxLength int
	minItems, maxItems   *int
}

// validateQuery closes gaps that typed nil requests and extensible string types
// would otherwise leave in the hand-maintained, cross-field validation.
func validateQuery(op string, q url.Values) error {
	invalid := func(field, msg string) error { return &ValidationError{Op: op, Field: field, Msg: msg} }
	for _, rule := range queryRules[op] {
		values, present := q[rule.name]
		if rule.required && (!present || len(values) == 0 || strings.TrimSpace(values[0]) == "") {
			return invalid(rule.name, "required")
		}
		if present && (rule.minItems != nil || rule.maxItems != nil) {
			count := len(values)
			if rule.csv {
				count = 0
				for _, value := range values {
					count += len(strings.Split(value, ","))
				}
			}
			if rule.minItems != nil && count < *rule.minItems || rule.maxItems != nil && count > *rule.maxItems {
				return invalid(rule.name, "invalid array length")
			}
		}
		for _, value := range values {
			items := []string{value}
			if rule.csv {
				items = strings.Split(value, ",")
			}
			for _, item := range items {
				if len(rule.enum) > 0 && !slices.Contains(rule.enum, item) {
					return invalid(rule.name, "must be one of "+strings.Join(rule.enum, ", "))
				}
				length := utf8.RuneCountInString(item)
				if length < rule.minLength || rule.maxLength > 0 && length > rule.maxLength {
					return invalid(rule.name, "invalid string length")
				}
				if rule.kind == "integer" || rule.kind == "number" {
					n, err := strconv.ParseFloat(item, 64)
					if err != nil || math.IsNaN(n) || math.IsInf(n, 0) || rule.min != nil && n < *rule.min || rule.max != nil && n > *rule.max {
						return invalid(rule.name, "number is outside the documented range")
					}
				}
			}
		}
	}
	return validateRelations(op, q)
}
