package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	polymarketdata "github.com/ivanzzeth/polymarket-go-data-client"
	"github.com/shopspring/decimal"
)

// MomentumSignal represents price momentum analysis
type MomentumSignal struct {
	MarketId        string
	MarketTitle     string
	CurrentPrice    decimal.Decimal
	PriceChange24h  decimal.Decimal
	VolumeChange    decimal.Decimal
	Trend           string
	Momentum        string
	TradingSignal   string
	SupportLevel    decimal.Decimal
	ResistanceLevel decimal.Decimal
}

func main() {
	fmt.Println("=== Price Momentum Analyzer ===")
	fmt.Println("This example analyzes price trends and volume to identify momentum trading opportunities")
	fmt.Println()

	// Create client
	client, err := polymarketdata.NewDataClient(&http.Client{})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Analyze specific markets
	testMarkets := []string{
		"0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917",
	}

	var signals []MomentumSignal

	for _, marketId := range testMarkets {
		signal, err := analyzeMomentum(client, marketId)
		if err != nil {
			log.Printf("Failed to analyze market %s: %v", marketId, err)
			continue
		}
		signals = append(signals, signal)
	}

	// Display signals
	fmt.Println("=== Momentum Analysis ===")
	fmt.Println()

	for i, signal := range signals {
		fmt.Printf("#%d Market: %s\n", i+1, signal.MarketTitle)
		fmt.Printf("   Market ID: %s\n", signal.MarketId[:20]+"...")
		fmt.Printf("   Current Price: $%s\n", signal.CurrentPrice.StringFixed(3))
		fmt.Printf("   24h Price Change: %.2f%%\n", signal.PriceChange24h.InexactFloat64())
		fmt.Printf("   Volume Trend: %s\n", signal.VolumeChange.StringFixed(1)+"x")
		fmt.Printf("   Price Trend: %s\n", signal.Trend)
		fmt.Printf("   Momentum: %s\n", signal.Momentum)

		// Show support/resistance
		if !signal.SupportLevel.IsZero() {
			fmt.Printf("   Support Level: $%s\n", signal.SupportLevel.StringFixed(3))
		}
		if !signal.ResistanceLevel.IsZero() {
			fmt.Printf("   Resistance Level: $%s\n", signal.ResistanceLevel.StringFixed(3))
		}

		// Trading signal
		fmt.Printf("\n   🎯 Trading Signal: %s\n", signal.TradingSignal)

		// Show detailed momentum analysis
		fmt.Println("\n   Detailed Momentum Analysis:")
		showMomentumDetails(client, signal.MarketId)
		fmt.Println()
	}

	// Summary
	fmt.Println("\n=== Trading Recommendations ===")
	bullishCount := 0
	bearishCount := 0
	neutralCount := 0

	for _, signal := range signals {
		if signal.Momentum == "STRONG BULLISH" || signal.Momentum == "BULLISH" {
			bullishCount++
		} else if signal.Momentum == "STRONG BEARISH" || signal.Momentum == "BEARISH" {
			bearishCount++
		} else {
			neutralCount++
		}
	}

	fmt.Printf("\nBullish Markets: %d\n", bullishCount)
	fmt.Printf("Bearish Markets: %d\n", bearishCount)
	fmt.Printf("Neutral Markets: %d\n", neutralCount)

	if bullishCount > 0 {
		fmt.Println("\n📈 Bullish momentum detected in some markets - consider long positions")
	}
	if bearishCount > 0 {
		fmt.Println("\n📉 Bearish momentum detected in some markets - consider short positions")
	}
}

func analyzeMomentum(client *polymarketdata.DataClient, marketId string) (MomentumSignal, error) {
	signal := MomentumSignal{
		MarketId: marketId,
	}

	// Get recent trades
	trades, err := client.GetTrades(context.Background(), &polymarketdata.GetTradesParams{
		Market: []string{marketId},
		Limit:  500,
	})
	if err != nil {
		return signal, fmt.Errorf("failed to get trades: %w", err)
	}

	if len(trades) == 0 {
		return signal, fmt.Errorf("no trades found")
	}

	signal.MarketTitle = trades[0].Title
	signal.CurrentPrice = trades[0].Price

	// Calculate price change over last 24 hours (approximated by trade sequence)
	cutoffTime := time.Now().Add(-24 * time.Hour).Unix()

	var recentTrades []polymarketdata.Trade
	var olderTrades []polymarketdata.Trade

	for _, trade := range trades {
		if trade.Timestamp > cutoffTime {
			recentTrades = append(recentTrades, trade)
		} else {
			olderTrades = append(olderTrades, trade)
		}
	}

	// Calculate price change
	if len(olderTrades) > 0 {
		oldPrice := olderTrades[0].Price
		priceChange := signal.CurrentPrice.Sub(oldPrice)
		signal.PriceChange24h = priceChange.Div(oldPrice).Mul(decimal.NewFromInt(100))
	}

	// Calculate volume change
	recentVolume := decimal.Zero
	olderVolume := decimal.Zero

	for _, trade := range recentTrades {
		recentVolume = recentVolume.Add(trade.Size.Mul(trade.Price))
	}
	for _, trade := range olderTrades {
		olderVolume = olderVolume.Add(trade.Size.Mul(trade.Price))
	}

	if !olderVolume.IsZero() && len(recentTrades) > 0 && len(olderTrades) > 0 {
		// Normalize by trade count
		avgRecentVol := recentVolume.Div(decimal.NewFromInt(int64(len(recentTrades))))
		avgOlderVol := olderVolume.Div(decimal.NewFromInt(int64(len(olderTrades))))
		signal.VolumeChange = avgRecentVol.Div(avgOlderVol)
	} else {
		signal.VolumeChange = decimal.NewFromInt(1)
	}

	// Determine trend
	if signal.PriceChange24h.GreaterThan(decimal.NewFromFloat(2.0)) {
		signal.Trend = "STRONG UPTREND 📈"
	} else if signal.PriceChange24h.GreaterThan(decimal.Zero) {
		signal.Trend = "UPTREND ↗️"
	} else if signal.PriceChange24h.LessThan(decimal.NewFromFloat(-2.0)) {
		signal.Trend = "STRONG DOWNTREND 📉"
	} else if signal.PriceChange24h.LessThan(decimal.Zero) {
		signal.Trend = "DOWNTREND ↘️"
	} else {
		signal.Trend = "SIDEWAYS ➡️"
	}

	// Determine momentum
	volumeIncreasing := signal.VolumeChange.GreaterThan(decimal.NewFromFloat(1.2))

	if signal.PriceChange24h.GreaterThan(decimal.NewFromFloat(3.0)) && volumeIncreasing {
		signal.Momentum = "STRONG BULLISH"
		signal.TradingSignal = "🟢 STRONG BUY - Uptrend with increasing volume"
	} else if signal.PriceChange24h.GreaterThan(decimal.NewFromFloat(1.0)) {
		signal.Momentum = "BULLISH"
		signal.TradingSignal = "🟢 BUY - Positive momentum"
	} else if signal.PriceChange24h.LessThan(decimal.NewFromFloat(-3.0)) && volumeIncreasing {
		signal.Momentum = "STRONG BEARISH"
		signal.TradingSignal = "🔴 STRONG SELL - Downtrend with increasing volume"
	} else if signal.PriceChange24h.LessThan(decimal.NewFromFloat(-1.0)) {
		signal.Momentum = "BEARISH"
		signal.TradingSignal = "🔴 SELL - Negative momentum"
	} else {
		signal.Momentum = "NEUTRAL"
		signal.TradingSignal = "⚪ HOLD - Wait for clear direction"
	}

	// Calculate support and resistance (using recent price range)
	minPrice := signal.CurrentPrice
	maxPrice := signal.CurrentPrice

	for _, trade := range recentTrades[:min(len(recentTrades), 100)] {
		if trade.Price.LessThan(minPrice) {
			minPrice = trade.Price
		}
		if trade.Price.GreaterThan(maxPrice) {
			maxPrice = trade.Price
		}
	}

	signal.SupportLevel = minPrice
	signal.ResistanceLevel = maxPrice

	return signal, nil
}

func showMomentumDetails(client *polymarketdata.DataClient, marketId string) {
	// Get activity data for more insights
	trades, err := client.GetTrades(context.Background(), &polymarketdata.GetTradesParams{
		Market: []string{marketId},
		Limit:  100,
	})
	if err != nil {
		fmt.Printf("      Failed to get trades: %v\n", err)
		return
	}

	if len(trades) == 0 {
		fmt.Println("      No trade data available")
		return
	}

	// Analyze trade velocity (trades per hour)
	if len(trades) >= 2 {
		timeSpan := trades[0].Timestamp - trades[len(trades)-1].Timestamp
		if timeSpan > 0 {
			hoursSpan := float64(timeSpan) / 3600.0
			tradesPerHour := float64(len(trades)) / hoursSpan
			fmt.Printf("      Trade Velocity: %.1f trades/hour\n", tradesPerHour)

			if tradesPerHour > 10 {
				fmt.Println("      🔥 High trading activity - strong interest")
			} else if tradesPerHour < 1 {
				fmt.Println("      ❄️  Low trading activity - limited interest")
			}
		}
	}

	// Analyze buy/sell imbalance in recent trades
	recentBuys := 0
	recentSells := 0
	for i, trade := range trades {
		if i >= 20 { // Last 20 trades
			break
		}
		if trade.Side == polymarketdata.TradeSideBuy {
			recentBuys++
		} else {
			recentSells++
		}
	}

	buyRatio := float64(recentBuys) / float64(recentBuys+recentSells) * 100
	fmt.Printf("      Recent Order Flow (last 20 trades): %.0f%% buy / %.0f%% sell\n",
		buyRatio, 100-buyRatio)

	if buyRatio > 70 {
		fmt.Println("      ⬆️  Strong buying pressure in recent trades")
	} else if buyRatio < 30 {
		fmt.Println("      ⬇️  Strong selling pressure in recent trades")
	}

	// Calculate price volatility
	prices := make([]decimal.Decimal, 0)
	for i, trade := range trades {
		if i >= 50 {
			break
		}
		prices = append(prices, trade.Price)
	}

	if len(prices) > 1 {
		avgPrice := decimal.Zero
		for _, p := range prices {
			avgPrice = avgPrice.Add(p)
		}
		avgPrice = avgPrice.Div(decimal.NewFromInt(int64(len(prices))))

		variance := decimal.Zero
		for _, p := range prices {
			diff := p.Sub(avgPrice)
			variance = variance.Add(diff.Mul(diff))
		}
		variance = variance.Div(decimal.NewFromInt(int64(len(prices))))

		volatility := variance.Div(avgPrice).Mul(decimal.NewFromInt(100))
		fmt.Printf("      Price Volatility: %.2f%%\n", volatility.InexactFloat64())

		if volatility.GreaterThan(decimal.NewFromFloat(5.0)) {
			fmt.Println("      ⚡ High volatility - larger price swings expected")
		} else if volatility.LessThan(decimal.NewFromFloat(1.0)) {
			fmt.Println("      📊 Low volatility - stable price action")
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
