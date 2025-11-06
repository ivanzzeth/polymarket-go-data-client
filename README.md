# Polymarket Go Data Client

A comprehensive Go SDK for the Polymarket Data API. This client library provides easy access to market data, user positions, trading activity, and more.

## Features

- ✅ Complete API coverage for all Polymarket Data endpoints
- ✅ Type-safe with proper decimal handling using `shopspring/decimal`
- ✅ Comprehensive error handling
- ✅ Full test coverage
- ✅ Production-ready examples for trading strategies

## Installation

```bash
go get github.com/ivanzzeth/polymarket-go-data-client
```

## Quick Start

```go
package main

import (
    "fmt"
    "net/http"

    polymarketdata "github.com/ivanzzeth/polymarket-go-data-client"
)

func main() {
    // Create client
    client, err := polymarketdata.NewDataClient(&http.Client{})
    if err != nil {
        panic(err)
    }

    // Check API health
    health, err := client.HealthCheck()
    if err != nil {
        panic(err)
    }
    fmt.Printf("API Status: %s\n", health.Data)

    // Get user positions
    positions, err := client.GetPositions(&polymarketdata.GetPositionsParams{
        User:  "0x56687bf447db6ffa42ffe2204a05edaa20f55839",
        Limit: 10,
    })
    if err != nil {
        panic(err)
    }

    for _, pos := range positions {
        fmt.Printf("Market: %s | PnL: %s\n", pos.Title, pos.CashPnl.String())
    }
}
```

## API Documentation

### Available Methods

#### Health & Status
- `HealthCheck()` - Check API availability

#### Positions
- `GetPositions(params)` - Get current user positions
- `GetClosedPositions(params)` - Get historical closed positions
- `GetPositionsValue(params)` - Get total position value

#### Trading
- `GetTrades(params)` - Get trade history
- `GetTradedMarketsCount(params)` - Get count of markets traded

#### Activity
- `GetActivity(params)` - Get user activity (trades, splits, merges, etc.)

#### Market Data
- `GetHolders(params)` - Get top token holders
- `GetOpenInterest(params)` - Get market open interest
- `GetLiveVolume(params)` - Get live trading volume

### Example Usage

```go
// Get recent trades for a market
trades, err := client.GetTrades(&polymarketdata.GetTradesParams{
    Market: []string{"0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"},
    Limit:  50,
})

// Get top holders
holders, err := client.GetHolders(&polymarketdata.GetHoldersParams{
    Market:     []string{marketId},
    Limit:      10,
    MinBalance: 1000,
})

// Get open interest
oi, err := client.GetOpenInterest(&polymarketdata.GetOpenInterestParams{
    Market: []string{marketId},
})
```

## Trading Strategy Examples

This repository includes 5 complete, production-ready examples demonstrating different trading strategies:

### 1. [Smart Money Tracker](examples/smart_money_tracker/)
**Track and follow profitable traders**

Identifies high-performing traders and monitors their positions and activities. Perfect for copy-trading strategies.

```bash
cd examples/smart_money_tracker && go run main.go
```

**Key Use Cases:**
- Copy trading profitable traders
- Validate market direction with smart money consensus
- Avoid positions opposite to successful traders

---

### 2. [Whale Watcher](examples/whale_watcher/)
**Monitor large holders and market concentration**

Analyzes token holder distribution to identify concentration risks and whale movements.

```bash
cd examples/whale_watcher && go run main.go
```

**Key Use Cases:**
- Assess concentration risk before entering markets
- Track whale accumulation/distribution
- Predict potential price movements from large holders

**Warning Levels:**
- ⚠️ Top 3 holders > 50% control
- ⚠️ Top 10 holders > 70% control

---

### 3. [Market Liquidity Analyzer](examples/market_liquidity_analyzer/)
**Find liquid markets and identify mispricing**

Evaluates market depth, volume, and liquidity to find optimal trading opportunities.

```bash
cd examples/market_liquidity_analyzer && go run main.go
```

**Key Use Cases:**
- Find liquid markets for large trades
- Identify illiquid markets prone to volatility
- Assess OI/Volume ratio for market health

**Liquidity Scoring:**
- 70-100: Excellent liquidity
- 30-70: Moderate liquidity
- <30: Poor liquidity, use caution

---

### 4. [Sentiment Reversal Detector](examples/sentiment_reversal_detector/)
**Identify overcrowded trades for contrarian opportunities**

Detects extreme sentiment imbalances (>85% one-sided) that often precede reversals.

```bash
cd examples/sentiment_reversal_detector && go run main.go
```

**Key Use Cases:**
- Find contrarian trading opportunities
- Identify overbought/oversold markets
- Fade overcrowded positions

**Signal Strength:**
- **STRONG**: >85% one-sided pressure
- **MODERATE**: 75-85% one-sided pressure
- **WEAK**: Balanced market

---

### 5. [Price Momentum Analyzer](examples/price_momentum_analyzer/)
**Analyze trends and catch momentum trades**

Identifies price trends, volume changes, and momentum to find trending markets.

```bash
cd examples/price_momentum_analyzer && go run main.go
```

**Key Use Cases:**
- Trend following strategies
- Breakout trading
- Support/resistance level identification

**Trading Signals:**
- 🟢 STRONG BUY: Price +3%+ with volume
- 🟢 BUY: Price +1-3%
- 🔴 STRONG SELL: Price -3%+ with volume
- 🔴 SELL: Price -1-3%
- ⚪ HOLD: Sideways action

---

## Project Structure

```
.
├── client.go           # Core client implementation
├── types.go            # All type definitions
├── health.go           # Health check endpoint
├── positions.go        # Position-related endpoints
├── trades.go           # Trading endpoints
├── activity.go         # Activity endpoints
├── holders.go          # Holder endpoints
├── misc.go             # Miscellaneous endpoints
├── *_test.go           # Comprehensive tests
└── examples/           # Trading strategy examples
    ├── smart_money_tracker/
    ├── whale_watcher/
    ├── market_liquidity_analyzer/
    ├── sentiment_reversal_detector/
    └── price_momentum_analyzer/
```

## Testing

Run all tests:

```bash
go test -v
```

Run specific test:

```bash
go test -v -run TestGetPositions
```

## Design Principles

### 1. Decimal Precision
All numeric values use `github.com/shopspring/decimal` for precise handling of financial data.

```go
type Position struct {
    Size     decimal.Decimal `json:"size"`
    AvgPrice decimal.Decimal `json:"avgPrice"`
    CashPnl  decimal.Decimal `json:"cashPnl"`
}
```

### 2. Pointer Parameters
Function parameters use pointers for mutability, while struct fields use values when zero is acceptable:

```go
// Function parameter is pointer
func (c *DataClient) GetPositions(params *GetPositionsParams) ([]Position, error)

// Field types - pointers only when needed
type GetPositionsParams struct {
    User          string           // Zero value OK
    Limit         int              // 0 means not set
    SizeThreshold *decimal.Decimal // Needs pointer
    Redeemable    *bool            // Needs pointer
}
```

### 3. Error Handling
Comprehensive error handling with context:

```go
if resp.StatusCode != http.StatusOK {
    var errResp ErrorResponse
    if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
        return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Error)
    }
    return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
}
```

## API Rate Limits

- The Polymarket Data API has rate limits
- Implement exponential backoff for production use
- Consider caching responses when appropriate

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Disclaimer

This SDK is for educational and research purposes. Trading prediction markets involves risk. Always:

- Do your own research
- Never risk more than you can afford to lose
- Past performance does not guarantee future results
- The examples are for demonstration only

## Resources

- [Polymarket Website](https://polymarket.com/)
- [Polymarket Data API Docs](https://docs.polymarket.com/)
- [Go Documentation](https://pkg.go.dev/)
- [Decimal Library](https://github.com/shopspring/decimal)

## Support

For issues, questions, or contributions:
- Open an issue on GitHub
- Check existing examples for guidance
- Review test files for usage patterns
