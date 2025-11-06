package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	polymarketdata "github.com/ivanzzeth/polymarket-go-data-client"
	"github.com/shopspring/decimal"
)

// SentimentSignal represents a potential reversal opportunity
type SentimentSignal struct {
	MarketId            string
	MarketTitle         string
	BuyPressure         decimal.Decimal
	SellPressure        decimal.Decimal
	HolderConcentration decimal.Decimal
	SignalStrength      string
	Reasoning           string
}

func main() {
	fmt.Println("=== Sentiment Reversal Detector ===")
	fmt.Println("This example identifies overcrowded trades for potential reversals")
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

	var signals []SentimentSignal

	for _, marketId := range testMarkets {
		signal, err := detectSentimentReversal(client, marketId)
		if err != nil {
			log.Printf("Failed to analyze market %s: %v", marketId, err)
			continue
		}
		signals = append(signals, signal)
	}

	// Display signals
	fmt.Println("=== Reversal Signals ===")
	fmt.Println()

	for i, signal := range signals {
		fmt.Printf("#%d Market: %s\n", i+1, signal.MarketTitle)
		fmt.Printf("   Market ID: %s\n", signal.MarketId[:20]+"...")
		fmt.Printf("   Buy Pressure: %.1f%%\n", signal.BuyPressure.InexactFloat64())
		fmt.Printf("   Sell Pressure: %.1f%%\n", signal.SellPressure.InexactFloat64())
		fmt.Printf("   Holder Concentration (Top 10): %.1f%%\n", signal.HolderConcentration.InexactFloat64())
		fmt.Printf("   Signal Strength: %s\n", signal.SignalStrength)
		fmt.Printf("   📊 Analysis: %s\n", signal.Reasoning)

		// Show detailed market analysis
		fmt.Println("\n   Detailed Analysis:")
		showDetailedAnalysis(client, signal.MarketId)
		fmt.Println()
	}

	// Summary
	fmt.Println("\n=== Summary ===")
	strongSignals := 0
	moderateSignals := 0

	for _, signal := range signals {
		if signal.SignalStrength == "STRONG" {
			strongSignals++
		} else if signal.SignalStrength == "MODERATE" {
			moderateSignals++
		}
	}

	fmt.Printf("Strong Reversal Signals: %d\n", strongSignals)
	fmt.Printf("Moderate Reversal Signals: %d\n", moderateSignals)
	fmt.Printf("Weak/No Signals: %d\n", len(signals)-strongSignals-moderateSignals)

	if strongSignals > 0 {
		fmt.Println("\n⚠️  STRONG signals detected - consider contrarian positions!")
	}
}

func detectSentimentReversal(client *polymarketdata.DataClient, marketId string) (SentimentSignal, error) {
	signal := SentimentSignal{
		MarketId: marketId,
	}

	// Get recent trades to analyze sentiment
	trades, err := client.GetTrades(context.Background(), &polymarketdata.GetTradesParams{
		Market: []string{marketId},
		Limit:  200, // Last 200 trades
	})
	if err != nil {
		return signal, fmt.Errorf("failed to get trades: %w", err)
	}

	if len(trades) == 0 {
		return signal, fmt.Errorf("no trades found")
	}

	signal.MarketTitle = trades[0].Title

	// Calculate buy/sell pressure
	buyCount := 0
	sellCount := 0
	buyVolume := decimal.Zero
	sellVolume := decimal.Zero

	for _, trade := range trades {
		tradeVolume := trade.Size.Mul(trade.Price)
		if trade.Side == polymarketdata.TradeSideBuy {
			buyCount++
			buyVolume = buyVolume.Add(tradeVolume)
		} else {
			sellCount++
			sellVolume = sellVolume.Add(tradeVolume)
		}
	}

	totalVolume := buyVolume.Add(sellVolume)
	if !totalVolume.IsZero() {
		signal.BuyPressure = buyVolume.Div(totalVolume).Mul(decimal.NewFromInt(100))
		signal.SellPressure = sellVolume.Div(totalVolume).Mul(decimal.NewFromInt(100))
	}

	// Get holder concentration
	holders, err := client.GetHolders(context.Background(), &polymarketdata.GetHoldersParams{
		Market: []string{marketId},
		Limit:  20,
	})
	if err == nil && len(holders) > 0 {
		// Analyze first token (usually most relevant)
		if len(holders[0].Holders) > 0 {
			totalAmount := decimal.Zero
			top10Amount := decimal.Zero

			for i, holder := range holders[0].Holders {
				totalAmount = totalAmount.Add(holder.Amount)
				if i < 10 {
					top10Amount = top10Amount.Add(holder.Amount)
				}
			}

			if !totalAmount.IsZero() {
				signal.HolderConcentration = top10Amount.Div(totalAmount).Mul(decimal.NewFromInt(100))
			}
		}
	}

	// Determine signal strength and reasoning
	signal.SignalStrength = "WEAK"
	signal.Reasoning = "Market conditions are relatively balanced"

	// Strong buy-side crowding (>85% buy pressure)
	if signal.BuyPressure.GreaterThan(decimal.NewFromInt(85)) {
		signal.SignalStrength = "STRONG"
		signal.Reasoning = "Extreme buy-side crowding detected! >85% of recent volume is buys. " +
			"Market may be overbought and vulnerable to reversal. Consider selling or waiting for pullback."
	} else if signal.BuyPressure.GreaterThan(decimal.NewFromInt(75)) {
		signal.SignalStrength = "MODERATE"
		signal.Reasoning = "Significant buy-side pressure (>75%). Market is getting crowded on one side. " +
			"Watch for signs of exhaustion or negative catalysts."
	}

	// Strong sell-side crowding (>85% sell pressure)
	if signal.SellPressure.GreaterThan(decimal.NewFromInt(85)) {
		signal.SignalStrength = "STRONG"
		signal.Reasoning = "Extreme sell-side crowding detected! >85% of recent volume is sells. " +
			"Market may be oversold and vulnerable to reversal. Consider buying the dip."
	} else if signal.SellPressure.GreaterThan(decimal.NewFromInt(75)) {
		signal.SignalStrength = "MODERATE"
		signal.Reasoning = "Significant sell-side pressure (>75%). Market is getting crowded on one side. " +
			"Watch for signs of capitulation or positive catalysts."
	}

	// Add holder concentration risk
	if signal.HolderConcentration.GreaterThan(decimal.NewFromInt(70)) {
		if signal.SignalStrength == "WEAK" {
			signal.SignalStrength = "MODERATE"
		}
		signal.Reasoning += " Additionally, top 10 holders control >70% of positions - " +
			"high concentration risk. Large holder movements could trigger sharp price changes."
	}

	return signal, nil
}

func showDetailedAnalysis(client *polymarketdata.DataClient, marketId string) {
	// Get recent trading patterns
	trades, err := client.GetTrades(context.Background(), &polymarketdata.GetTradesParams{
		Market: []string{marketId},
		Limit:  50,
	})
	if err != nil {
		fmt.Printf("      Failed to get trades: %v\n", err)
		return
	}

	if len(trades) == 0 {
		fmt.Println("      No trade data available")
		return
	}

	// Analyze trade size distribution
	largeTradeCount := 0
	smallTradeCount := 0
	medianSize := decimal.Zero

	for _, trade := range trades {
		tradeValue := trade.Size.Mul(trade.Price)
		if tradeValue.GreaterThan(decimal.NewFromInt(1000)) {
			largeTradeCount++
		} else if tradeValue.LessThan(decimal.NewFromInt(100)) {
			smallTradeCount++
		}
		medianSize = medianSize.Add(tradeValue)
	}

	if len(trades) > 0 {
		medianSize = medianSize.Div(decimal.NewFromInt(int64(len(trades))))
	}

	fmt.Printf("      Recent Trades: %d\n", len(trades))
	fmt.Printf("      Large Trades (>$1000): %d (%.1f%%)\n",
		largeTradeCount,
		float64(largeTradeCount)*100/float64(len(trades)))
	fmt.Printf("      Small Trades (<$100): %d (%.1f%%)\n",
		smallTradeCount,
		float64(smallTradeCount)*100/float64(len(trades)))
	fmt.Printf("      Average Trade Size: $%s\n", medianSize.StringFixed(2))

	// Interpretation
	if largeTradeCount > len(trades)/3 {
		fmt.Println("      💡 High proportion of large trades - institutional activity")
	}
	if smallTradeCount > len(trades)*2/3 {
		fmt.Println("      💡 Mostly small trades - retail dominated")
	}

	// Price trend
	if len(trades) >= 2 {
		latestPrice := trades[0].Price
		earlierPrice := trades[len(trades)-1].Price
		priceChange := latestPrice.Sub(earlierPrice).Div(earlierPrice).Mul(decimal.NewFromInt(100))

		fmt.Printf("      Price Movement: %.2f%% (over last %d trades)\n",
			priceChange.InexactFloat64(), len(trades))

		if priceChange.Abs().GreaterThan(decimal.NewFromInt(5)) {
			fmt.Println("      📊 Significant price movement detected")
		}
	}
}
