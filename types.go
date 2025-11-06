package polymarketdata

import "github.com/shopspring/decimal"

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Data string `json:"data"`
}

// ErrorResponse represents a standard error response from the API
type ErrorResponse struct {
	Error string `json:"error"`
}

// SortBy represents the sort field for positions
type SortBy string

const (
	SortByCurrent    SortBy = "CURRENT"
	SortByInitial    SortBy = "INITIAL"
	SortByTokens     SortBy = "TOKENS"
	SortByCashPnl    SortBy = "CASHPNL"
	SortByPercentPnl SortBy = "PERCENTPNL"
	SortByTitle      SortBy = "TITLE"
	SortByResolving  SortBy = "RESOLVING"
	SortByPrice      SortBy = "PRICE"
	SortByAvgPrice   SortBy = "AVGPRICE"
)

// SortDirection represents the sort direction
type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc SortDirection = "DESC"
)

// Position represents a user's position in a market
type Position struct {
	ProxyWallet         string          `json:"proxyWallet"`
	Asset               string          `json:"asset"`
	ConditionId         string          `json:"conditionId"`
	Size                decimal.Decimal `json:"size"`
	AvgPrice            decimal.Decimal `json:"avgPrice"`
	InitialValue        decimal.Decimal `json:"initialValue"`
	CurrentValue        decimal.Decimal `json:"currentValue"`
	CashPnl             decimal.Decimal `json:"cashPnl"`
	PercentPnl          decimal.Decimal `json:"percentPnl"`
	TotalBought         decimal.Decimal `json:"totalBought"`
	RealizedPnl         decimal.Decimal `json:"realizedPnl"`
	PercentRealizedPnl  decimal.Decimal `json:"percentRealizedPnl"`
	CurPrice            decimal.Decimal `json:"curPrice"`
	Redeemable          bool            `json:"redeemable"`
	Mergeable           bool            `json:"mergeable"`
	Title               string          `json:"title"`
	Slug                string          `json:"slug"`
	Icon                string          `json:"icon"`
	EventSlug           string          `json:"eventSlug"`
	Outcome             string          `json:"outcome"`
	OutcomeIndex        int             `json:"outcomeIndex"`
	OppositeOutcome     string          `json:"oppositeOutcome"`
	OppositeAsset       string          `json:"oppositeAsset"`
	EndDate             string          `json:"endDate"`
	NegativeRisk        bool            `json:"negativeRisk"`
}

// GetPositionsParams represents parameters for getting user positions
type GetPositionsParams struct {
	User          string           // Required: User address. Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Market        []string         // Optional: Comma-separated list of condition IDs. 0x-prefixed 64-hex string. Mutually exclusive with EventId.
	EventId       []int            // Optional: Comma-separated list of event IDs. Mutually exclusive with Market.
	SizeThreshold *decimal.Decimal // Optional: Default 1, required range: x >= 0
	Redeemable    *bool            // Optional: Default false
	Mergeable     *bool            // Optional: Default false
	Limit         int              // Optional: Default 100, required range: 0 <= x <= 500. 0 means not set.
	Offset        int              // Optional: Default 0, required range: 0 <= x <= 10000. 0 means not set.
	SortBy        SortBy           // Optional: Default TOKENS. Available options: CURRENT, INITIAL, TOKENS, CASHPNL, PERCENTPNL, TITLE, RESOLVING, PRICE, AVGPRICE
	SortDirection SortDirection    // Optional: Default DESC. Available options: ASC, DESC
	Title         string           // Optional: Maximum length: 100
}

// TradeSide represents the side of a trade
type TradeSide string

const (
	TradeSideBuy  TradeSide = "BUY"
	TradeSideSell TradeSide = "SELL"
)

// FilterType represents the filter type for trades
type FilterType string

const (
	FilterTypeCash   FilterType = "CASH"
	FilterTypeTokens FilterType = "TOKENS"
)

// Trade represents a trade record
type Trade struct {
	ProxyWallet           string          `json:"proxyWallet"` // User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Side                  TradeSide       `json:"side"`        // Available options: BUY, SELL
	Asset                 string          `json:"asset"`
	ConditionId           string          `json:"conditionId"` // 0x-prefixed 64-hex string. Example: "0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"
	Size                  decimal.Decimal `json:"size"`
	Price                 decimal.Decimal `json:"price"`
	Timestamp             int64           `json:"timestamp"`
	Title                 string          `json:"title"`
	Slug                  string          `json:"slug"`
	Icon                  string          `json:"icon"`
	EventSlug             string          `json:"eventSlug"`
	Outcome               string          `json:"outcome"`
	OutcomeIndex          int             `json:"outcomeIndex"`
	Name                  string          `json:"name"`
	Pseudonym             string          `json:"pseudonym"`
	Bio                   string          `json:"bio"`
	ProfileImage          string          `json:"profileImage"`
	ProfileImageOptimized string          `json:"profileImageOptimized"`
	TransactionHash       string          `json:"transactionHash"`
}

// GetTradesParams represents parameters for getting trades
type GetTradesParams struct {
	Limit        int              // Optional: Default 100, required range: 0 <= x <= 10000. 0 means not set.
	Offset       int              // Optional: Default 0, required range: 0 <= x <= 10000. 0 means not set.
	TakerOnly    *bool            // Optional: Default true
	FilterType   FilterType       // Optional: Must be provided together with FilterAmount. Available options: CASH, TOKENS
	FilterAmount *decimal.Decimal // Optional: Must be provided together with FilterType. Required range: x >= 0
	Market       []string         // Optional: Comma-separated list of condition IDs (0x-prefixed 64-hex string). Mutually exclusive with EventId.
	EventId      []int            // Optional: Comma-separated list of event IDs. Mutually exclusive with Market.
	User         string           // Optional: User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Side         TradeSide        // Optional: Available options: BUY, SELL
}

// TradedMarketsCount represents the total number of markets a user has traded
type TradedMarketsCount struct {
	User   string `json:"user"`  // User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Traded int    `json:"traded"` // Total number of markets traded
}

// GetTradedMarketsCountParams represents parameters for getting traded markets count
type GetTradedMarketsCountParams struct {
	User string // Required: User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
}

// ActivityType represents the type of activity
type ActivityType string

const (
	ActivityTypeTrade      ActivityType = "TRADE"
	ActivityTypeSplit      ActivityType = "SPLIT"
	ActivityTypeMerge      ActivityType = "MERGE"
	ActivityTypeRedeem     ActivityType = "REDEEM"
	ActivityTypeReward     ActivityType = "REWARD"
	ActivityTypeConversion ActivityType = "CONVERSION"
)

// ActivitySortBy represents the sort field for activities
type ActivitySortBy string

const (
	ActivitySortByTimestamp ActivitySortBy = "TIMESTAMP"
	ActivitySortByTokens    ActivitySortBy = "TOKENS"
	ActivitySortByCash      ActivitySortBy = "CASH"
)

// Activity represents a user's on-chain activity
type Activity struct {
	ProxyWallet              string          `json:"proxyWallet"` // User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Timestamp                int64           `json:"timestamp"`
	ConditionId              string          `json:"conditionId"` // 0x-prefixed 64-hex string. Example: "0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"
	Type                     ActivityType    `json:"type"`        // Available options: TRADE, SPLIT, MERGE, REDEEM, REWARD, CONVERSION
	Size                     decimal.Decimal `json:"size"`
	UsdcSize                 decimal.Decimal `json:"usdcSize"`
	TransactionHash          string          `json:"transactionHash"`
	Price                    decimal.Decimal `json:"price"`
	Asset                    string          `json:"asset"`
	Side                     TradeSide       `json:"side"` // Available options: BUY, SELL
	OutcomeIndex             int             `json:"outcomeIndex"`
	Title                    string          `json:"title"`
	Slug                     string          `json:"slug"`
	Icon                     string          `json:"icon"`
	EventSlug                string          `json:"eventSlug"`
	Outcome                  string          `json:"outcome"`
	Name                     string          `json:"name"`
	Pseudonym                string          `json:"pseudonym"`
	Bio                      string          `json:"bio"`
	ProfileImage             string          `json:"profileImage"`
	ProfileImageOptimized    string          `json:"profileImageOptimized"`
}

// GetActivityParams represents parameters for getting user activity
type GetActivityParams struct {
	Limit         int            // Optional: Default 100, required range: 0 <= x <= 500. 0 means not set.
	Offset        int            // Optional: Default 0, required range: 0 <= x <= 10000. 0 means not set.
	User          string         // Required: User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Market        []string       // Optional: Comma-separated list of condition IDs (0x-prefixed 64-hex string). Mutually exclusive with EventId.
	EventId       []int          // Optional: Comma-separated list of event IDs. Mutually exclusive with Market.
	Type          []ActivityType // Optional: Activity type filters
	Start         int64          // Optional: Start timestamp, required range: x >= 0. 0 means not set.
	End           int64          // Optional: End timestamp, required range: x >= 0. 0 means not set.
	SortBy        ActivitySortBy // Optional: Default TIMESTAMP. Available options: TIMESTAMP, TOKENS, CASH
	SortDirection SortDirection  // Optional: Default DESC. Available options: ASC, DESC
	Side          TradeSide      // Optional: Available options: BUY, SELL
}

// Holder represents a holder of a market token
type Holder struct {
	ProxyWallet             string          `json:"proxyWallet"`             // User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Bio                     string          `json:"bio"`                     // User bio
	Asset                   string          `json:"asset"`                   // Asset address
	Pseudonym               string          `json:"pseudonym"`               // User pseudonym
	Amount                  decimal.Decimal `json:"amount"`                  // Amount held
	DisplayUsernamePublic   bool            `json:"displayUsernamePublic"`   // Whether username is public
	OutcomeIndex            int             `json:"outcomeIndex"`            // Outcome index
	Name                    string          `json:"name"`                    // User name
	ProfileImage            string          `json:"profileImage"`            // Profile image URL
	ProfileImageOptimized   string          `json:"profileImageOptimized"`   // Optimized profile image URL
}

// MarketHolders represents holders for a specific market token
type MarketHolders struct {
	Token   string   `json:"token"`   // Token address
	Holders []Holder `json:"holders"` // List of holders
}

// GetHoldersParams represents parameters for getting top holders
type GetHoldersParams struct {
	Limit      int      // Optional: Default 100, required range: 0 <= x <= 500. 0 means not set.
	Market     []string // Required: Comma-separated list of condition IDs (0x-prefixed 64-hex string)
	MinBalance int      // Optional: Default 1, required range: 0 <= x <= 999999. 0 means not set.
}

// UserValue represents the total value of a user's positions
type UserValue struct {
	User  string          `json:"user"`  // User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Value decimal.Decimal `json:"value"` // Total value
}

// GetValueParams represents parameters for getting user's total position value
type GetValueParams struct {
	User   string   // Required: User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Market []string // Optional: List of condition IDs (0x-prefixed 64-hex string)
}

// ClosedPositionSortBy represents the sort field for closed positions
type ClosedPositionSortBy string

const (
	ClosedPositionSortByRealizedPnl ClosedPositionSortBy = "REALIZEDPNL"
	ClosedPositionSortByTitle       ClosedPositionSortBy = "TITLE"
	ClosedPositionSortByPrice       ClosedPositionSortBy = "PRICE"
	ClosedPositionSortByAvgPrice    ClosedPositionSortBy = "AVGPRICE"
)

// ClosedPosition represents a closed position for a user
type ClosedPosition struct {
	ProxyWallet     string          `json:"proxyWallet"` // User Profile Address (0x-prefixed, 40 hex chars). Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Asset           string          `json:"asset"`
	ConditionId     string          `json:"conditionId"` // 0x-prefixed 64-hex string. Example: "0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"
	AvgPrice        decimal.Decimal `json:"avgPrice"`
	TotalBought     decimal.Decimal `json:"totalBought"`
	RealizedPnl     decimal.Decimal `json:"realizedPnl"`
	CurPrice        decimal.Decimal `json:"curPrice"`
	Title           string          `json:"title"`
	Slug            string          `json:"slug"`
	Icon            string          `json:"icon"`
	EventSlug       string          `json:"eventSlug"`
	Outcome         string          `json:"outcome"`
	OutcomeIndex    int             `json:"outcomeIndex"`
	OppositeOutcome string          `json:"oppositeOutcome"`
	OppositeAsset   string          `json:"oppositeAsset"`
	EndDate         string          `json:"endDate"`
}

// GetClosedPositionsParams represents parameters for getting closed positions
type GetClosedPositionsParams struct {
	User          string                   // Required: The address of the user. Example: "0x56687bf447db6ffa42ffe2204a05edaa20f55839"
	Market        []string                 // Optional: The conditionId of the market (0x-prefixed 64-hex string). Cannot be used with EventId.
	Title         string                   // Optional: Filter by market title. Maximum length: 100
	EventId       []int                    // Optional: The event id. Returns positions for all markets for those event ids. Cannot be used with Market.
	Limit         int                      // Optional: Default 50, required range: 0 <= x <= 500. 0 means not set.
	Offset        int                      // Optional: Default 0, required range: 0 <= x <= 10000. 0 means not set.
	SortBy        ClosedPositionSortBy     // Optional: Default REALIZEDPNL. Available options: REALIZEDPNL, TITLE, PRICE, AVGPRICE
	SortDirection SortDirection            // Optional: Default DESC. Available options: ASC, DESC
}

// OpenInterest represents the open interest for a market
type OpenInterest struct {
	Market string          `json:"market"` // 0x-prefixed 64-hex string. Example: "0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"
	Value  decimal.Decimal `json:"value"`  // Open interest value
}

// GetOpenInterestParams represents parameters for getting open interest
type GetOpenInterestParams struct {
	Market []string // Optional: List of market condition IDs (0x-prefixed 64-hex string)
}

// LiveVolumeMarket represents the volume for a specific market
type LiveVolumeMarket struct {
	Market string          `json:"market"` // 0x-prefixed 64-hex string. Example: "0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"
	Value  decimal.Decimal `json:"value"`  // Volume value
}

// LiveVolume represents the live volume for an event
type LiveVolume struct {
	Total   decimal.Decimal    `json:"total"`   // Total volume
	Markets []LiveVolumeMarket `json:"markets"` // List of market volumes
}

// GetLiveVolumeParams represents parameters for getting live volume
type GetLiveVolumeParams struct {
	Id int // Required: Event ID, required range: x >= 1
}
