package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"

	polymarketdata "github.com/ivanzzeth/polymarket-go-data-client"
	"github.com/shopspring/decimal"
)

// TraderProfile represents a trader with their performance metrics
type TraderProfile struct {
	Address         string
	TotalPnL        decimal.Decimal
	WinRate         decimal.Decimal
	MarketsTraded   int
	CurrentValue    decimal.Decimal
	ActivePositions int
}

func main() {
	fmt.Println("=== Smart Money Tracker ===")
	fmt.Println("This example finds profitable traders and analyzes their current positions")
	fmt.Println()

	// Create client
	client, err := polymarketdata.NewClient(&http.Client{})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// For demonstration, we'll analyze a known active trader
	// In production, you'd scan multiple traders or use a database of traders
	testTraders := []string{
		"0x56687bf447db6ffa42ffe2204a05edaa20f55839",
		// Add more trader addresses here
	}

	var profiles []TraderProfile

	for _, traderAddr := range testTraders {
		profile, err := analyzeTrader(client, traderAddr)
		if err != nil {
			log.Printf("Failed to analyze trader %s: %v", traderAddr, err)
			continue
		}
		profiles = append(profiles, profile)
	}

	// Sort by total PnL (highest first)
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].TotalPnL.GreaterThan(profiles[j].TotalPnL)
	})

	// Display top traders
	fmt.Println("=== Top Profitable Traders ===")
	for i, profile := range profiles {
		fmt.Printf("\n#%d Trader: %s\n", i+1, profile.Address)
		fmt.Printf("  Total PnL: $%s\n", profile.TotalPnL.StringFixed(2))
		fmt.Printf("  Markets Traded: %d\n", profile.MarketsTraded)
		fmt.Printf("  Current Portfolio Value: $%s\n", profile.CurrentValue.StringFixed(2))
		fmt.Printf("  Active Positions: %d\n", profile.ActivePositions)

		// Show their current top positions
		if profile.ActivePositions > 0 {
			fmt.Println("\n  Current Top Positions:")
			showTopPositions(client, profile.Address, 3)
		}

		// Show recent activity
		fmt.Println("\n  Recent Trading Activity:")
		showRecentActivity(client, profile.Address, 5)
	}
}

func analyzeTrader(client *polymarketdata.Client, address string) (TraderProfile, error) {
	profile := TraderProfile{
		Address: address,
	}

	// Get closed positions to calculate historical PnL
	closedPositions, err := client.GetClosedPositions(context.Background(), &polymarketdata.GetClosedPositionsParams{
		User:          address,
		Limit:         500,
		SortBy:        polymarketdata.ClosedPositionSortByRealizedPnl,
		SortDirection: polymarketdata.SortDirectionDesc,
	})
	if err != nil {
		return profile, fmt.Errorf("failed to get closed positions: %w", err)
	}

	// Calculate total PnL
	totalPnL := decimal.Zero
	for _, pos := range closedPositions {
		totalPnL = totalPnL.Add(pos.RealizedPnl)
	}
	profile.TotalPnL = totalPnL

	// Get markets traded count
	tradedCount, err := client.GetTradedMarketsCount(context.Background(), &polymarketdata.GetTradedMarketsCountParams{
		User: address,
	})
	if err != nil {
		return profile, fmt.Errorf("failed to get traded markets count: %w", err)
	}
	profile.MarketsTraded = tradedCount.Traded

	// Get current positions value
	values, err := client.GetPositionsValue(context.Background(), &polymarketdata.GetValueParams{
		User: address,
	})
	if err != nil {
		return profile, fmt.Errorf("failed to get positions value: %w", err)
	}
	if len(values) > 0 {
		profile.CurrentValue = values[0].Value
	}

	// Get active positions count
	positions, err := client.GetPositions(context.Background(), &polymarketdata.GetPositionsParams{
		User:  address,
		Limit: 500,
	})
	if err != nil {
		return profile, fmt.Errorf("failed to get positions: %w", err)
	}
	profile.ActivePositions = len(positions)

	return profile, nil
}

func showTopPositions(client *polymarketdata.Client, address string, limit int) {
	positions, err := client.GetPositions(context.Background(), &polymarketdata.GetPositionsParams{
		User:          address,
		Limit:         limit,
		SortBy:        polymarketdata.SortByCurrent,
		SortDirection: polymarketdata.SortDirectionDesc,
	})
	if err != nil {
		log.Printf("Failed to get positions: %v", err)
		return
	}

	for i, pos := range positions {
		if i >= limit {
			break
		}
		fmt.Printf("    %d. %s\n", i+1, pos.Title)
		fmt.Printf("       Outcome: %s | Size: %s | Avg Price: $%s | Current: $%s\n",
			pos.Outcome,
			pos.Size.StringFixed(2),
			pos.AvgPrice.StringFixed(3),
			pos.CurPrice.StringFixed(3),
		)
		fmt.Printf("       PnL: $%s (%.2f%%)\n",
			pos.CashPnl.StringFixed(2),
			pos.PercentPnl.InexactFloat64(),
		)
	}
}

func showRecentActivity(client *polymarketdata.Client, address string, limit int) {
	activities, err := client.GetActivity(context.Background(), &polymarketdata.GetActivityParams{
		User:          address,
		Limit:         limit,
		SortBy:        polymarketdata.ActivitySortByTimestamp,
		SortDirection: polymarketdata.SortDirectionDesc,
	})
	if err != nil {
		log.Printf("Failed to get activity: %v", err)
		return
	}

	for i, activity := range activities {
		if i >= limit {
			break
		}
		fmt.Printf("    %d. [%s] %s\n", i+1, activity.Type, activity.Title)
		if activity.Type == polymarketdata.ActivityTypeTrade {
			fmt.Printf("       Side: %s | Size: %s | Price: $%s | USDC: $%s\n",
				activity.Side,
				activity.Size.StringFixed(2),
				activity.Price.StringFixed(3),
				activity.UsdcSize.StringFixed(2),
			)
		} else {
			fmt.Printf("       Size: %s | USDC: $%s\n",
				activity.Size.StringFixed(2),
				activity.UsdcSize.StringFixed(2),
			)
		}
	}
}
