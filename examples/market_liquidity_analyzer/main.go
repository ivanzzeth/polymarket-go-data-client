package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	polymarketdata "github.com/ivanzzeth/polymarket-go-data-client"
	"github.com/shopspring/decimal"
)

// MarketMetrics represents liquidity and activity metrics for a market
type MarketMetrics struct {
	MarketId        string
	OpenInterest    decimal.Decimal
	Volume          decimal.Decimal
	OIToVolumeRatio decimal.Decimal
	LiquidityScore  decimal.Decimal
}

func main() {
	fmt.Println("=== Market Liquidity Analyzer ===")
	fmt.Println("This example analyzes market liquidity to find trading opportunities")
	fmt.Println()

	// Create client
	client, err := polymarketdata.NewDataClient(&http.Client{})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get global open interest to see all active markets
	fmt.Println("Fetching global open interest data...")
	allOI, err := client.GetOpenInterest(&polymarketdata.GetOpenInterestParams{})
	if err != nil {
		log.Fatalf("Failed to get open interest: %v", err)
	}

	fmt.Printf("Found %d markets\n\n", len(allOI))

	// Analyze specific markets (in production, you'd analyze all or filter by criteria)
	testMarkets := []string{
		"0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917",
	}

	var metrics []MarketMetrics

	for _, marketId := range testMarkets {
		if marketId == "GLOBAL" {
			continue
		}

		metric, err := analyzeMarketLiquidity(client, marketId)
		if err != nil {
			log.Printf("Failed to analyze market %s: %v", marketId, err)
			continue
		}
		metrics = append(metrics, metric)
	}

	// Sort by liquidity score (custom metric)
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].LiquidityScore.GreaterThan(metrics[j].LiquidityScore)
	})

	// Display results
	fmt.Println("=== Market Liquidity Analysis ===")
	fmt.Println()

	for i, metric := range metrics {
		fmt.Printf("#%d Market: %s\n", i+1, metric.MarketId[:20]+"...")
		fmt.Printf("   Open Interest: $%s\n", metric.OpenInterest.StringFixed(2))
		fmt.Printf("   24h Volume: $%s\n", metric.Volume.StringFixed(2))
		fmt.Printf("   OI/Volume Ratio: %.2f\n", metric.OIToVolumeRatio.InexactFloat64())
		fmt.Printf("   Liquidity Score: %.2f\n", metric.LiquidityScore.InexactFloat64())

		// Interpretation
		fmt.Println("\n   Analysis:")
		if metric.OIToVolumeRatio.GreaterThan(decimal.NewFromInt(10)) {
			fmt.Println("   ⚠️  Low trading activity relative to open interest")
			fmt.Println("   💡 Potential for volatility on news or large orders")
		} else if metric.OIToVolumeRatio.LessThan(decimal.NewFromInt(2)) {
			fmt.Println("   ✅ High trading activity - very liquid market")
			fmt.Println("   💡 Good for larger position sizes")
		} else {
			fmt.Println("   📊 Moderate liquidity - normal market conditions")
		}

		// Get recent trades to analyze price action
		fmt.Println("\n   Recent Price Action:")
		analyzePriceAction(client, metric.MarketId, 10)

		fmt.Println()
	}

	// Summary recommendations
	fmt.Println("=== Trading Recommendations ===")
	fmt.Println("\n🎯 Best Markets for Trading:")
	for i, metric := range metrics {
		if i >= 3 { // Top 3
			break
		}
		if metric.LiquidityScore.GreaterThan(decimal.NewFromInt(50)) {
			fmt.Printf("   • Market %s (Score: %.1f)\n",
				metric.MarketId[:20]+"...",
				metric.LiquidityScore.InexactFloat64())
		}
	}

	fmt.Println("\n⚠️  Markets to Watch (Low Liquidity):")
	for i := len(metrics) - 1; i >= 0 && i >= len(metrics)-3; i-- {
		if metrics[i].LiquidityScore.LessThan(decimal.NewFromInt(30)) {
			fmt.Printf("   • Market %s (Score: %.1f)\n",
				metrics[i].MarketId[:20]+"...",
				metrics[i].LiquidityScore.InexactFloat64())
			fmt.Println("     → Use limit orders and smaller sizes")
		}
	}
}

func analyzeMarketLiquidity(client *polymarketdata.DataClient, marketId string) (MarketMetrics, error) {
	metric := MarketMetrics{
		MarketId: marketId,
	}

	// Get open interest
	oiData, err := client.GetOpenInterest(&polymarketdata.GetOpenInterestParams{
		Market: []string{marketId},
	})
	if err != nil {
		return metric, fmt.Errorf("failed to get open interest: %w", err)
	}

	if len(oiData) > 0 {
		metric.OpenInterest = oiData[0].Value
	}

	// Get recent trades to estimate volume
	trades, err := client.GetTrades(&polymarketdata.GetTradesParams{
		Market: []string{marketId},
		Limit:  1000, // Last 1000 trades as proxy for recent activity
	})
	if err != nil {
		return metric, fmt.Errorf("failed to get trades: %w", err)
	}

	// Calculate total volume
	totalVolume := decimal.Zero
	for _, trade := range trades {
		tradeValue := trade.Size.Mul(trade.Price)
		totalVolume = totalVolume.Add(tradeValue)
	}
	metric.Volume = totalVolume

	// Calculate OI to Volume ratio
	if !metric.Volume.IsZero() {
		metric.OIToVolumeRatio = metric.OpenInterest.Abs().Div(metric.Volume)
	} else {
		metric.OIToVolumeRatio = decimal.NewFromInt(999) // Very high = no liquidity
	}

	// Calculate liquidity score (0-100)
	// Score based on: volume (50%), number of trades (30%), low OI/Vol ratio (20%)
	volumeScore := decimal.Min(metric.Volume.Div(decimal.NewFromInt(10000)), decimal.NewFromInt(50))
	tradesScore := decimal.Min(decimal.NewFromInt(int64(len(trades))).Div(decimal.NewFromInt(10)), decimal.NewFromInt(30))

	ratioScore := decimal.NewFromInt(20)
	if metric.OIToVolumeRatio.GreaterThan(decimal.NewFromInt(10)) {
		ratioScore = decimal.NewFromInt(5)
	} else if metric.OIToVolumeRatio.LessThan(decimal.NewFromInt(3)) {
		ratioScore = decimal.NewFromInt(20)
	} else {
		ratioScore = decimal.NewFromInt(10)
	}

	metric.LiquidityScore = volumeScore.Add(tradesScore).Add(ratioScore)

	return metric, nil
}

func analyzePriceAction(client *polymarketdata.DataClient, marketId string, limit int) {
	trades, err := client.GetTrades(&polymarketdata.GetTradesParams{
		Market: []string{marketId},
		Limit:  limit,
	})
	if err != nil {
		fmt.Printf("      Failed to get trades: %v\n", err)
		return
	}

	if len(trades) == 0 {
		fmt.Println("      No recent trades")
		return
	}

	// Analyze last N trades
	buyVolume := decimal.Zero
	sellVolume := decimal.Zero
	minPrice := trades[0].Price
	maxPrice := trades[0].Price

	for _, trade := range trades {
		tradeValue := trade.Size.Mul(trade.Price)

		if trade.Side == polymarketdata.TradeSideBuy {
			buyVolume = buyVolume.Add(tradeValue)
		} else {
			sellVolume = sellVolume.Add(tradeValue)
		}

		if trade.Price.LessThan(minPrice) {
			minPrice = trade.Price
		}
		if trade.Price.GreaterThan(maxPrice) {
			maxPrice = trade.Price
		}
	}

	totalVolume := buyVolume.Add(sellVolume)
	buyPct := decimal.Zero
	if !totalVolume.IsZero() {
		buyPct = buyVolume.Div(totalVolume).Mul(decimal.NewFromInt(100))
	}

	priceRange := maxPrice.Sub(minPrice)
	volatility := decimal.Zero
	if !maxPrice.IsZero() {
		volatility = priceRange.Div(maxPrice).Mul(decimal.NewFromInt(100))
	}

	fmt.Printf("      Latest Price: $%s\n", trades[0].Price.StringFixed(3))
	fmt.Printf("      Price Range (last %d trades): $%s - $%s\n",
		limit, minPrice.StringFixed(3), maxPrice.StringFixed(3))
	fmt.Printf("      Volatility: %.2f%%\n", volatility.InexactFloat64())
	fmt.Printf("      Buy/Sell Ratio: %.1f%% buy / %.1f%% sell\n",
		buyPct.InexactFloat64(),
		decimal.NewFromInt(100).Sub(buyPct).InexactFloat64())

	// Sentiment interpretation
	if buyPct.GreaterThan(decimal.NewFromInt(65)) {
		fmt.Println("      📈 Strong buying pressure")
	} else if buyPct.LessThan(decimal.NewFromInt(35)) {
		fmt.Println("      📉 Strong selling pressure")
	} else {
		fmt.Println("      ⚖️  Balanced market")
	}
}
