// Package nebula is the Go SDK for the Nebula Public API
// (https://docs.hiddensystems.ai/nebula/). It covers all 41 operations of the
// account-authenticated API plus the separate opt-in x402 pay-per-request
// surface via the x402 adapter package.
//
//	client, err := nebula.New(os.Getenv("NEBULA_API_KEY"))
//	assets, err := client.Lookup.ListAssets(ctx, &nebula.ListAssetsRequest{
//	    Query: nebula.StringPtr("NVIDIA"),
//	})
//
// All methods take a context.Context first. All responses carry typed data
// plus ResponseMetadata (rate limit and credits) via the embedded
// BaseResponse. Market data is not investment advice.
package nebula

// AssetClass is the class of a market asset.
type AssetClass string

const (
	AssetClassEquity    AssetClass = "equity"
	AssetClassCommodity AssetClass = "commodity"
	AssetClassIndex     AssetClass = "index"
	AssetClassCrypto    AssetClass = "crypto"
)

// Kind distinguishes market assets from people/topics (subjects).
type Kind string

const (
	KindAsset   Kind = "asset"
	KindSubject Kind = "subject"
)

// Granularity is the bucket size of a time series.
type Granularity string

const (
	Granularity1h  Granularity = "1h"
	Granularity2h  Granularity = "2h"
	Granularity4h  Granularity = "4h"
	Granularity12h Granularity = "12h"
	Granularity1d  Granularity = "1d"
)

// ScreenerSort is the ranking measure of the screener.
type ScreenerSort string

const (
	ScreenerSortShare             ScreenerSort = "share"
	ScreenerSortShareChange       ScreenerSort = "share_change"
	ScreenerSortShareChangePct    ScreenerSort = "share_change_pct"
	ScreenerSortPosts             ScreenerSort = "posts"
	ScreenerSortNetBullishPosts   ScreenerSort = "net_bullish_posts"
	ScreenerSortPriceChange24hPct ScreenerSort = "price_change_24h_pct"
)

// SortOrder is a result ordering direction.
type SortOrder string

const (
	SortOrderDesc SortOrder = "desc"
	SortOrderAsc  SortOrder = "asc"
)

// InsightOrderBy is the ordering of the insights list.
type InsightOrderBy string

const (
	InsightOrderByStartTimeAsc        InsightOrderBy = "start_time_asc"
	InsightOrderByStartTimeDesc       InsightOrderBy = "start_time_desc"
	InsightOrderByUpdatedAtAsc        InsightOrderBy = "updated_at_asc"
	InsightOrderByUpdatedAtDesc       InsightOrderBy = "updated_at_desc"
	InsightOrderByCreatedAtAsc        InsightOrderBy = "created_at_asc"
	InsightOrderByCreatedAtDesc       InsightOrderBy = "created_at_desc"
	InsightOrderByImportanceScoreAsc  InsightOrderBy = "importance_score_asc"
	InsightOrderByImportanceScoreDesc InsightOrderBy = "importance_score_desc"
)

// PredictionProvider is a prediction-market venue.
type PredictionProvider string

const (
	PredictionProviderPolymarket PredictionProvider = "polymarket"
	PredictionProviderKalshi     PredictionProvider = "kalshi"
)

// PredictionStatus is the lifecycle state of a prediction market listing.
type PredictionStatus string

const (
	PredictionStatusOpen     PredictionStatus = "open"
	PredictionStatusClosed   PredictionStatus = "closed"
	PredictionStatusResolved PredictionStatus = "resolved"
	PredictionStatusDisputed PredictionStatus = "disputed"
	PredictionStatusAll      PredictionStatus = "all"
)

// CalendarDomain is an event domain on the calendar endpoint.
type CalendarDomain string

const (
	CalendarDomainCrypto            CalendarDomain = "crypto"
	CalendarDomainEquities          CalendarDomain = "equities"
	CalendarDomainCommodities       CalendarDomain = "commodities"
	CalendarDomainMacro             CalendarDomain = "macro"
	CalendarDomainPolitics          CalendarDomain = "politics"
	CalendarDomainTechnology        CalendarDomain = "technology"
	CalendarDomainPredictionMarkets CalendarDomain = "prediction_markets"
)

// InsightDomain is an event domain on the insights endpoint.
type InsightDomain string

const (
	InsightDomainCrypto            InsightDomain = "crypto"
	InsightDomainEquities          InsightDomain = "equities"
	InsightDomainCommodities       InsightDomain = "commodities"
	InsightDomainMacro             InsightDomain = "macro"
	InsightDomainPolitics          InsightDomain = "politics"
	InsightDomainTechnology        InsightDomain = "technology"
	InsightDomainSports            InsightDomain = "sports"
	InsightDomainPredictionMarkets InsightDomain = "prediction_markets"
)

// EventType is a calendar event classification. See the OpenAPI spec for the
// full set; the constants below enumerate every value documented in v1.0.2.
type EventType string

const (
	EventTypeEarnings                   EventType = "earnings"
	EventTypeGuidance                   EventType = "guidance"
	EventTypeInvestorDay                EventType = "investor_day"
	EventTypeShareholderMeeting         EventType = "shareholder_meeting"
	EventTypeDividend                   EventType = "dividend"
	EventTypeStockSplit                 EventType = "stock_split"
	EventTypeBuyback                    EventType = "buyback"
	EventTypeSpinOff                    EventType = "spin_off"
	EventTypeAcquisition                EventType = "acquisition"
	EventTypeFundraising                EventType = "fundraising"
	EventTypeBankruptcy                 EventType = "bankruptcy"
	EventTypeLayoffs                    EventType = "layoffs"
	EventTypePartnership                EventType = "partnership"
	EventTypeAppointment                EventType = "appointment"
	EventTypePricingChange              EventType = "pricing_change"
	EventTypeClinicalTrial              EventType = "clinical_trial"
	EventTypeCorporateAction            EventType = "corporate_action"
	EventTypeIPO                        EventType = "ipo"
	EventTypeListing                    EventType = "listing"
	EventTypeDelisting                  EventType = "delisting"
	EventTypeIndexRebalance             EventType = "index_rebalance"
	EventTypeContractExpiry             EventType = "contract_expiry"
	EventTypeMarketHoliday              EventType = "market_holiday"
	EventTypeETF                        EventType = "etf"
	EventTypeEconomicRelease            EventType = "economic_release"
	EventTypeMonetaryPolicy             EventType = "monetary_policy"
	EventTypeCentralBankSpeech          EventType = "central_bank_speech"
	EventTypeBudget                     EventType = "budget"
	EventTypeDebtAuction                EventType = "debt_auction"
	EventTypeTradePolicy                EventType = "trade_policy"
	EventTypeSanctions                  EventType = "sanctions"
	EventTypeElection                   EventType = "election"
	EventTypeReferendum                 EventType = "referendum"
	EventTypeVote                       EventType = "vote"
	EventTypeHearing                    EventType = "hearing"
	EventTypeSummit                     EventType = "summit"
	EventTypeTreaty                     EventType = "treaty"
	EventTypePolicyDeadline             EventType = "policy_deadline"
	EventTypeRegulatoryDecision         EventType = "regulatory_decision"
	EventTypeRegulatoryDeadline         EventType = "regulatory_deadline"
	EventTypeFilingDeadline             EventType = "filing_deadline"
	EventTypeCourtRuling                EventType = "court_ruling"
	EventTypeLegal                      EventType = "legal"
	EventTypeUnlock                     EventType = "unlock"
	EventTypeBurn                       EventType = "burn"
	EventTypeMigration                  EventType = "migration"
	EventTypeSnapshot                   EventType = "snapshot"
	EventTypeAirdrop                    EventType = "airdrop"
	EventTypeTokenSale                  EventType = "token_sale"
	EventTypeStaking                    EventType = "staking"
	EventTypeHalving                    EventType = "halving"
	EventTypeMainnet                    EventType = "mainnet"
	EventTypeTestnet                    EventType = "testnet"
	EventTypeFork                       EventType = "fork"
	EventTypeUpgrade                    EventType = "upgrade"
	EventTypeGovernance                 EventType = "governance"
	EventTypeLaunch                     EventType = "launch"
	EventTypeConference                 EventType = "conference"
	EventTypeHackathon                  EventType = "hackathon"
	EventTypeAMA                        EventType = "ama"
	EventTypeCompetition                EventType = "competition"
	EventTypePredictionMarketResolution EventType = "prediction_market_resolution"
	EventTypeDeadline                   EventType = "deadline"
	EventTypeSports                     EventType = "sports"
	EventTypeOther                      EventType = "other"
)

// EventStatus is a calendar event lifecycle state.
type EventStatus string

const (
	EventStatusScheduled EventStatus = "scheduled"
	EventStatusTentative EventStatus = "tentative"
	EventStatusCancelled EventStatus = "cancelled"
	EventStatusPostponed EventStatus = "postponed"
)

// ChatRole is the role of an agent conversation message. "developer" is
// OpenAI's newer name for system instructions and is treated as system.
type ChatRole string

const (
	ChatRoleSystem    ChatRole = "system"
	ChatRoleDeveloper ChatRole = "developer"
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
)

// ThinkingMode selects the Agent reasoning mode.
type ThinkingMode string

const (
	ThinkingModeAuto ThinkingMode = "auto"
	ThinkingModeDeep ThinkingMode = "deep"
)

// ForecastStatus discriminates a forecast response.
type ForecastStatus string

const (
	ForecastStatusOK              ForecastStatus = "ok"
	ForecastStatusNotForecastable ForecastStatus = "not_forecastable"
)

// ForecastReason explains why an asset cannot be forecast.
type ForecastReason string

const (
	ForecastReasonTooLittleHistory ForecastReason = "too_little_history"
	ForecastReasonStaleData        ForecastReason = "stale_data"
	ForecastReasonMissingPrice     ForecastReason = "missing_price"
)

// Locale selects the language of generated text. Sent as X-Nebula-Locale
// (or the `locale` query parameter). It is an open string type: a future
// locale works without a new SDK release.
type Locale string

const (
	LocaleSpanish  Locale = "es"
	LocaleItalian  Locale = "it"
	LocaleTurkish  Locale = "tr"
	LocaleRussian  Locale = "ru"
	LocaleKorean   Locale = "ko"
	LocaleJapanese Locale = "ja"
	LocaleChinese  Locale = "zh"
)

// Cluster IDs documented in the OpenAPI `clusters`/`cluster_id` descriptions.
// Cluster IDs are an extensible string type; these constants cover the
// catalog published in spec v1.0.2.
const (
	ClusterSmartMoney      = "smart_money"
	ClusterDumbMoney       = "dumb_money"
	ClusterTopDegens       = "top_degens"
	ClusterEarlyMovers     = "early_movers"
	ClusterLateMovers      = "late_movers"
	ClusterDevelopers      = "developers"
	ClusterFounders        = "founders"
	ClusterTeamLeaders     = "team_leaders"
	ClusterProjectAccounts = "project_accounts"
	ClusterVCs             = "vcs"
	ClusterAssetManagers   = "asset_managers"
	ClusterDAOGuildFunds   = "dao_guild_funds"
	ClusterIncubators      = "incubators"
	ClusterExchanges       = "exchanges"
	ClusterKOLs            = "kols"
	ClusterCelebrities     = "celebrities"
	ClusterMedia           = "media"
	ClusterMarketer        = "marketer"
	ClusterShiller         = "shiller"
	ClusterTAAnalyst       = "ta_analyst"
	ClusterFAAnalyst       = "fa_analyst"
	ClusterMacroAnalyst    = "macro_analyst"
	ClusterOnchainAnalyst  = "onchain_analyst"
	ClusterInvestigators   = "investigators"
	ClusterPositionTrader  = "position_trader"
	ClusterSwingTrader     = "swing_trader"
	ClusterDayTrader       = "day_trader"
	ClusterPermaBull       = "perma_bull"
	ClusterDoomer          = "doomer"
	ClusterTribalist       = "tribalist"
	ClusterPolitical       = "political"
	ClusterNFT             = "nft"
	ClusterBots            = "bots"
	ClusterAutomatedAI     = "automated_ai"
	ClusterWallStreetBets  = "wallstreetbets"
	ClusterReddit          = "reddit"
)

var communityClusters = map[string]bool{
	ClusterWallStreetBets: true,
	ClusterReddit:         true,
}

// IsCommunityCluster reports whether id is a posting community
// (wallstreetbets, reddit) rather than an author cluster. Community data only
// exists for traditional asset classes, which the SDK validates locally.
func IsCommunityCluster(id string) bool { return communityClusters[id] }

func hasCommunity(clusters []string) bool {
	for _, c := range clusters {
		if communityClusters[c] {
			return true
		}
	}
	return false
}

func isCommunity(id string) bool { return communityClusters[id] }
