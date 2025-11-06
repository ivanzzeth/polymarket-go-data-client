# Whale Watcher

## Overview

Whale Watcher monitors large holders ("whales") in Polymarket markets. By tracking whale positions and activities, you can identify concentration risks, predict potential price movements, and make more informed trading decisions.

## Key Features

- **Identify Top Holders**: See who controls the most tokens in each market
- **Concentration Analysis**: Calculate how much top holders control (e.g., Top 3, Top 10)
- **Activity Monitoring**: Track recent trades and position changes by whales
- **Risk Assessment**: Automatic warnings for highly concentrated markets

## Use Cases

### 1. Risk Management
Avoid entering markets where Top 3 holders control >70% of supply - these markets are vulnerable to manipulation.

### 2. Whale Tracking
When a whale makes a large exit, it often signals a shift in market sentiment. Follow their moves to stay ahead.

### 3. Liquidity Assessment
High concentration means low liquidity. Your large orders could significantly move the price.

## How to Run

```bash
cd examples/whale_watcher
go run main.go
```

## Example Output

```
=== Whale Watcher ===
This example monitors large holders (whales) in specific markets

Analyzing market: 0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917

=== Token 1: 48331043336612883890938759509493159234755048973500640148014422747788308965732 ===
Total Holders: 20

Top 10 Whales:
  #1: 0xa5ef39c3...
      Address: 0xa5ef39c3d3e10d0b270233af41cac69796b12966
      Amount: 196965002.56 tokens (86.11% of total)
      Recent Activity:
        No recent activity in this market

  #2: I95153360
      Address: 0x5557f74a8b21634d2daa199b6364cef137276f2c
      Amount: 7107980.00 tokens (3.11% of total)
      Recent Activity:
        - [BUY TRADE] Size: 192434.49 @ $0.056 (Total: $10759.39)

=== Concentration Metrics ===
Top 3 holders control: 90.44% of tokens
Top 10 holders control: 96.44% of tokens
⚠️  WARNING: High concentration! Top 3 holders control >50% of tokens
⚠️  WARNING: Very high concentration! Top 10 holders control >70% of tokens
```

## Key Metrics Explained

### Token Concentration
- **Top 3 Control**: If >50%, market is highly concentrated
- **Top 10 Control**: If >70%, extremely concentrated
- **Impact**: High concentration = high manipulation risk

### Holder Distribution
- **Balanced** (<40% in Top 10): Healthy, decentralized market
- **Concentrated** (40-70% in Top 10): Moderate risk
- **Highly Concentrated** (>70% in Top 10): High risk, proceed with caution

### Activity Patterns
- **Recent Buys**: Whale is accumulating, bullish signal
- **Recent Sells**: Whale is distributing, bearish signal
- **No Activity**: Whale is holding, neutral or waiting

## Trading Strategy Recommendations

### Strategy 1: Avoid Concentrated Markets
1. Check concentration before entering
2. If Top 3 > 70%, skip this market
3. Find more decentralized alternatives

### Strategy 2: Follow Whale Accumulation
1. Monitor whale activity
2. When whales are buying heavily, consider following
3. Exit when whales start selling

### Strategy 3: Frontrun Whale Exits
1. Set up alerts for whale sell activity
2. Exit before the whale dumps their full position
3. Re-enter after price stabilizes

## Code Customization

### Monitor Multiple Markets
Modify the market list:

```go
markets := []string{
    "0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917",
    "0xYOUR_MARKET_ID_1",
    "0xYOUR_MARKET_ID_2",
}
```

### Adjust Minimum Balance Filter
Change the minimum tokens to display:

```go
GetHoldersParams{
    Market:     []string{marketId},
    Limit:      20,
    MinBalance: 5000,  // Only show holders with 5000+ tokens
}
```

### Customize Alert Thresholds
Modify concentration warning levels:

```go
if top3Pct.GreaterThan(decimal.NewFromInt(60)) {  // Change from 50 to 60
    fmt.Println("⚠️  WARNING: High concentration!")
}
```

## Important Notes

⚠️ **Risk Warnings**:
- Concentrated markets are vulnerable to manipulation
- Whale movements can cause extreme volatility
- Always consider concentration in your risk assessment

⚠️ **Data Limitations**:
- Holder data may be slightly delayed
- Some whales use multiple addresses
- Smart contract addresses may appear as holders

⚠️ **Market Impact**:
- In highly concentrated markets, a single large order can move prices significantly
- Use limit orders and smaller position sizes

## Advanced Usage

### 1. Whale Alert System
Monitor whale activity and send real-time alerts:

```go
func monitorWhales(marketId string) {
    previousHolders := getHolders(marketId)

    for {
        currentHolders := getHolders(marketId)

        // Detect significant changes
        for i, holder := range currentHolders {
            if holder.Amount.Sub(previousHolders[i].Amount).Abs().GreaterThan(threshold) {
                sendWhaleAlert(holder, "Position changed significantly")
            }
        }

        previousHolders = currentHolders
        time.Sleep(5 * time.Minute)
    }
}
```

### 2. Cross-Market Whale Analysis
Track the same whale across multiple markets:

```go
func trackWhaleAcrossMarkets(whaleAddr string) {
    // Get all markets
    markets := getAllMarkets()

    for _, market := range markets {
        position := getWhalePosition(whaleAddr, market)
        if position.Size.GreaterThan(threshold) {
            fmt.Printf("Whale has large position in: %s\n", market.Title)
        }
    }
}
```

### 3. Concentration History Tracking
Build a database of historical concentration:

```go
type ConcentrationHistory struct {
    MarketId    string
    Timestamp   time.Time
    Top3Pct     decimal.Decimal
    Top10Pct    decimal.Decimal
}

// Track how concentration changes over time
func trackConcentrationTrend(marketId string) {
    // Store daily snapshots
    // Analyze if concentration is increasing or decreasing
    // Alert if suddenly increasing (accumulation phase)
}
```

## APIs Used

This example uses the following APIs:

- `GetHolders()` - Retrieve top token holders for a market
- `GetActivity()` - Get recent trading activity for specific addresses

## FAQ

**Q: How do I find market IDs to monitor?**
A: You can:
1. Browse Polymarket website and extract from URLs
2. Use `GetOpenInterest()` to get all active markets
3. Check popular markets from the Polymarket API documentation

**Q: What if a whale uses multiple addresses?**
A: Unfortunately, it's difficult to link addresses without additional data. Focus on monitoring obvious large holders.

**Q: Is high concentration always bad?**
A: Not always. Some markets naturally have knowledgeable concentrated holders. However, it does increase risk.

**Q: How often should I check whale activity?**
A: Depends on your strategy:
- Day trading: Every 15-30 minutes
- Swing trading: Every few hours
- Position trading: Once or twice daily

## Concentration Risk Levels

| Top 10 Control | Risk Level | Recommendation |
|----------------|------------|----------------|
| < 40% | Low | Safe to trade normally |
| 40-60% | Moderate | Use caution, smaller positions |
| 60-80% | High | Reduce position size significantly |
| > 80% | Extreme | Avoid or trade very carefully |

## Related Examples

- [Smart Money Tracker](../smart_money_tracker/) - Follow profitable traders
- [Sentiment Reversal Detector](../sentiment_reversal_detector/) - Find overcrowded trades
- [Market Liquidity Analyzer](../market_liquidity_analyzer/) - Assess market liquidity
