# Price Momentum Analyzer

## Overview

Price Momentum Analyzer identifies trending markets through price action and volume analysis. Catch momentum trades early and ride trends with confidence.

## Key Features

- **Trend Detection**: Identify uptrends, downtrends, and sideways markets
- **Momentum Scoring**: STRONG BULLISH/BEARISH, BULLISH/BEARISH, or NEUTRAL
- **Volume Confirmation**: Validate trends with volume changes
- **Support/Resistance**: Identify key price levels
- **Volatility Analysis**: Measure price stability

## How to Run

```bash
cd examples/price_momentum_analyzer
go run main.go
```

## Momentum Signals

### 🟢 STRONG BUY
- Price up >3% with increasing volume
- Strong uptrend confirmed
- **Action**: Enter long positions

### 🟢 BUY
- Price up 1-3%
- Positive momentum
- **Action**: Consider long positions

### 🔴 STRONG SELL
- Price down >3% with increasing volume
- Strong downtrend confirmed
- **Action**: Enter short positions or exit longs

### 🔴 SELL
- Price down 1-3%
- Negative momentum
- **Action**: Consider short positions or reduce longs

### ⚪ HOLD
- Sideways price action
- No clear direction
- **Action**: Wait for breakout

## Trading Strategies

### Trend Following
1. Wait for STRONG momentum signal
2. Enter in direction of trend
3. Hold until momentum shifts
4. Use trailing stop loss

### Breakout Trading
1. Identify consolidation (NEUTRAL signal)
2. Wait for volume spike
3. Enter on breakout direction
4. Target resistance/support levels

### Mean Reversion
1. Find markets with extreme short-term momentum
2. Wait for initial reversal
3. Enter counter-trend
4. Quick exit (< 24 hours)

## Key Metrics Explained

### Price Change 24h
- **> +3%**: Strong uptrend
- **+1% to +3%**: Moderate uptrend
- **-1% to +1%**: Sideways
- **-3% to -1%**: Moderate downtrend
- **< -3%**: Strong downtrend

### Volume Change
- **> 1.5x**: Increasing interest, trend likely to continue
- **0.8x - 1.2x**: Normal volume
- **< 0.8x**: Decreasing interest, trend may be exhausting

### Trade Velocity
- **> 10 trades/hour**: High activity market
- **1-10 trades/hour**: Moderate activity
- **< 1 trade/hour**: Low activity, be cautious

## APIs Used

- `GetTrades()` - Price and volume data
- Trade velocity and order flow analysis
- Price range and volatility calculations

## Related Examples

- [Market Liquidity Analyzer](../market_liquidity_analyzer/) - Assess liquidity
- [Sentiment Reversal Detector](../sentiment_reversal_detector/) - Find reversals
- [Smart Money Tracker](../smart_money_tracker/) - Follow smart traders
