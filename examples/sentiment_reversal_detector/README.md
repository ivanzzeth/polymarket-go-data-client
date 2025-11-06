# Sentiment Reversal Detector

## Overview

Sentiment Reversal Detector identifies overcrowded trades that may be due for a reversal. When >85% of traders are on one side, the market may be overextended.

## Key Features

- **Buy/Sell Pressure Analysis**: Measure trading flow imbalance
- **Concentration Risk**: Factor in holder concentration
- **Signal Strength**: STRONG, MODERATE, or WEAK reversal signals
- **Trade Size Distribution**: Detect institutional vs retail activity

## How to Run

```bash
cd examples/sentiment_reversal_detector
go run main.go
```

## Signal Interpretation

### STRONG Signal (>85% one-sided)
- Extreme crowding detected
- High probability of reversal
- **Action**: Consider contrarian position

### MODERATE Signal (75-85% one-sided)
- Significant pressure building
- Watch for exhaustion signals
- **Action**: Prepare for potential reversal

### WEAK Signal
- Balanced market conditions
- No clear reversal setup
- **Action**: Wait for clearer signals

## Trading Strategies

### Contrarian Entry
1. Wait for STRONG signal
2. Enter opposite to crowd sentiment
3. Set tight stop loss (3-5%)
4. Take profit on reversal (5-10%)

### Fading the Move
1. Identify extreme buy/sell pressure
2. Wait for first sign of reversal
3. Enter with momentum shift
4. Exit when balance restores

## APIs Used

- `GetTrades()` - Analyze buy/sell flow
- `GetHolders()` - Check concentration
- `GetActivity()` - Recent market activity

## Related Examples

- [Smart Money Tracker](../smart_money_tracker/) - Follow profitable traders
- [Price Momentum Analyzer](../price_momentum_analyzer/) - Trend analysis
