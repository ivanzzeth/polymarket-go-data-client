# Market Liquidity Analyzer

## Overview

Market Liquidity Analyzer evaluates market depth and trading activity to identify opportunities. Find liquid markets for large trades or spot illiquid markets prone to price manipulation.

## Key Features

- **Liquidity Scoring**: Custom score (0-100) based on volume, trades, and OI/Volume ratio
- **Volume Analysis**: Compare open interest to trading volume
- **Price Action**: Recent price trends and volatility
- **Buy/Sell Pressure**: Analyze order flow imbalance

## How to Run

```bash
cd examples/market_liquidity_analyzer
go run main.go
```

## Key Metrics

### OI/Volume Ratio
- **< 2**: Very liquid, high activity
- **2-10**: Normal liquidity
- **> 10**: Low liquidity, watch for volatility

### Liquidity Score
- **70-100**: Excellent for trading
- **50-70**: Good liquidity
- **30-50**: Moderate, use caution
- **< 30**: Poor liquidity, avoid large orders

## Trading Strategies

### For Liquid Markets (Score > 70)
- Safe for larger position sizes
- Lower slippage on market orders
- Easier entry/exit

### For Illiquid Markets (Score < 30)
- Use limit orders only
- Split large orders
- Watch for sudden price spikes on news

## APIs Used

- `GetOpenInterest()` - Market open interest
- `GetTrades()` - Recent trading activity
- Volume and price analysis

## Related Examples

- [Whale Watcher](../whale_watcher/) - Monitor large holders
- [Price Momentum Analyzer](../price_momentum_analyzer/) - Trend analysis
