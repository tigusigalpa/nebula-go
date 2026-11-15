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
