# Smart Money Tracker

## Overview

The Smart Money Tracker helps you identify and follow the most profitable traders on Polymarket. By analyzing historical performance data, this tool reveals "smart money" movements and allows you to monitor their current positions and trading activity.

## Key Features

- **Identify Profitable Traders**: Analyze historical closed positions to find top performers
- **Monitor Current Positions**: See what markets smart money is trading right now
- **Track Trading Activity**: Real-time monitoring of their latest buy/sell operations
- **Multi-dimensional Analysis**: Evaluate traders by PnL, market count, portfolio value, and win rate

## Use Cases

### 1. Copy Trading Strategy
When you discover a trader with >$1,000,000 in historical profits opening a large position in a market, this could be a valuable trading signal.

### 2. Market Validation
If multiple profitable traders are heavily positioned in the same direction on a market, this increases confidence in that outcome.

### 3. Risk Avoidance
If your position is opposite to all profitable traders, this might be a warning signal to reconsider.

## How to Run

```bash
cd examples/smart_money_tracker
go run main.go
```

## Example Output

```
=== Smart Money Tracker ===
This example finds profitable traders and analyzes their current positions

=== Top Profitable Traders ===

#1 Trader: 0x56687bf447db6ffa42ffe2204a05edaa20f55839
  Total PnL: $22053933.75
  Markets Traded: 14
  Current Portfolio Value: $0.01
  Active Positions: 0

  Recent Trading Activity:
    1. [REDEEM] Which party wins 2024 US Presidential Election?
       Size: 120469.39 | USDC: $120469.39
    2. [TRADE] Henry Cavill announced as next James Bond?
       Side: BUY | Size: 199.95 | Price: $0.966 | USDC: $193.15
```

## Key Metrics Explained

### Total PnL
- Cumulative profit/loss across all closed positions
- **Higher is better**: Indicates strong historical performance
- **Recommended threshold**: Focus on traders with >$100,000

### Markets Traded
- Total number of different markets the trader has participated in
- **Sweet spot**: Too few might indicate inexperience, too many might indicate lack of focus
- **Ideal range**: 10-50 markets

### Current Portfolio Value
- Total value of all current positions
- **High value = High conviction**: Trader is actively engaged
- **Low value**: May be waiting for opportunities or has exited

### Active Positions
- Number of markets currently held
- **Concentrated vs Diversified**:
  - 1-3 positions: Highly concentrated, strong conviction
  - 5-10 positions: Diversified approach
  - >10 positions: Over-diversified, may lack focus

## Trading Strategy Recommendations

### Strategy 1: Direct Copy Trading
1. Find traders with >$500K in profits
2. Monitor their new positions
3. Enter within 24 hours of their entry
4. Set stop loss at 5-10%

### Strategy 2: Cluster Validation
1. Monitor 5-10 profitable traders
2. Enter when 3+ traders are heavily positioned in the same direction
3. This signals "smart money consensus"

### Strategy 3: Contrarian Indicator
1. Find the worst performing traders (modify code sorting)
2. Take opposite positions
3. Works well in emotional markets

## Code Customization

### Add More Traders
Modify the `testTraders` array:

```go
testTraders := []string{
    "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
    "0xYOUR_TRADER_ADDRESS_1",
    "0xYOUR_TRADER_ADDRESS_2",
}
```

### Change Sorting Criteria
Modify `GetClosedPositions` parameters:

```go
SortBy: polymarketdata.ClosedPositionSortByRealizedPnl,  // Sort by PnL
// or
SortBy: polymarketdata.ClosedPositionSortByAvgPrice,     // Sort by avg price
```

### Filter Small Trades
Add filtering in `showRecentActivity`:

```go
if activity.UsdcSize.LessThan(decimal.NewFromInt(100)) {
    continue  // Skip trades <$100
}
```

## Important Notes

⚠️ **Disclaimer**:
- Past performance does not guarantee future results
- Even profitable traders make mistakes
- Combine multiple signals with your own research

⚠️ **API Limitations**:
- Single query limited to 500 records
- For analyzing more traders, implement batch querying

⚠️ **Real-time Considerations**:
- Data may have a few minutes delay
- Not suitable for very short-term trading (seconds/minutes)

## Advanced Usage

### 1. Database Storage
Store analysis results to build a leaderboard:

```go
type TraderRank struct {
    Address     string
    TotalPnL    decimal.Decimal
    WinRate     float64
    UpdatedAt   time.Time
}
```

### 2. Alert System
Send notifications when profitable traders make new moves:

```go
func monitorTrader(address string) {
    lastActivity := getLastActivity(address)
    for {
        current := getLatestActivity(address)
        if current.Timestamp > lastActivity.Timestamp {
            sendAlert(current)
        }
        time.Sleep(5 * time.Minute)
    }
}
```

### 3. Historical Backtesting
Validate copy trading strategy effectiveness:

```go
func backtestStrategy(traderAddr string, startDate, endDate time.Time) {
    // Get historical trades
    // Simulate copying
    // Calculate returns
}
```

## APIs Used

This example uses the following APIs:

- `GetClosedPositions()` - Retrieve historical closed positions
- `GetTradedMarketsCount()` - Get number of markets traded
- `GetPositionsValue()` - Get total current position value
- `GetPositions()` - Get current position details
- `GetActivity()` - Get trading activity records

## FAQ

**Q: How do I find more trader addresses?**
A: Several ways:
1. Check leaderboards on Polymarket website
2. Use `GetHolders()` to find large holders in markets
3. Analyze trade records in popular markets

**Q: Why do some traders have negative Total PnL?**
A: This means the trader has overall losses. Not recommended to follow.

**Q: How to implement real-time monitoring?**
A: Put the code in a loop, periodically (e.g., every 5 minutes) re-query data and compare changes.

## Related Examples

- [Whale Watcher](../whale_watcher/) - Monitor large holders
- [Sentiment Reversal Detector](../sentiment_reversal_detector/) - Find contrarian opportunities
- [Price Momentum Analyzer](../price_momentum_analyzer/) - Trend analysis
