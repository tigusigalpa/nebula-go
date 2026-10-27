package nebula

import "context"

type ctxKey int

const (
	ctxKeyLocale ctxKey = iota
	ctxKeyLocaleMode
)

// WithRequestLocale returns a context that overrides the client's locale for
// one call, sent as X-Nebula-Locale (or ?locale= under LocaleModeQuery).
//
//	ctx := nebula.WithRequestLocale(ctx, nebula.LocaleRussian)
//	res, err := client.Social.Sentiment(ctx, req)
func WithRequestLocale(ctx context.Context, l Locale) context.Context {
	return context.WithValue(ctx, ctxKeyLocale, l)
}

// WithRequestLocaleMode returns a context that overrides how the locale is
// transported for one call (header or query).
func WithRequestLocaleMode(ctx context.Context, m LocaleMode) context.Context {
	return context.WithValue(ctx, ctxKeyLocaleMode, m)
}

func localeFor(ctx context.Context, cfg *config) (Locale, LocaleMode) {
	locale := cfg.locale
	mode := cfg.localeMode
	if v, ok := ctx.Value(ctxKeyLocale).(Locale); ok {
		locale = v
	}
	if v, ok := ctx.Value(ctxKeyLocaleMode).(LocaleMode); ok {
		mode = v
	}
	return locale, mode
}
