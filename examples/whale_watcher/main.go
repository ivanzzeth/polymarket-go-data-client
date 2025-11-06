package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	polymarketdata "github.com/ivanzzeth/polymarket-go-data-client"
	"github.com/shopspring/decimal"
)

// WhaleInfo represents information about a large holder
type WhaleInfo struct {
	Address               string
	Amount                decimal.Decimal
	Outcome               string
	OutcomeIndex          int
	Name                  string
	Pseudonym             string
	DisplayUsernamePublic bool
}

func main() {
	fmt.Println("=== Whale Watcher ===")
	fmt.Println("This example monitors large holders (whales) in specific markets")
	fmt.Println()

	// Create client
	client, err := polymarketdata.NewDataClient(&http.Client{})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example market: A popular market condition ID
	// In production, you'd monitor multiple markets or fetch from a list
	marketId := "0xdd22472e552920b8438158ea7238bfadfa4f736aa4cee91a6b86c39ead110917"

	fmt.Printf("Analyzing market: %s\n\n", marketId)

	// Get holders for the market
	holders, err := client.GetHolders(&polymarketdata.GetHoldersParams{
		Market:     []string{marketId},
		Limit:      20,
		MinBalance: 1000, // Only show holders with at least 1000 tokens
	})
	if err != nil {
		log.Fatalf("Failed to get holders: %v", err)
	}

	if len(holders) == 0 {
		fmt.Println("No holder data available for this market")
		return
	}

	// Analyze each token (Yes/No outcomes)
	for tokenIdx, marketHolder := range holders {
		fmt.Printf("=== Token %d: %s ===\n", tokenIdx+1, marketHolder.Token)
		fmt.Printf("Total Holders: %d\n\n", len(marketHolder.Holders))

		if len(marketHolder.Holders) == 0 {
			fmt.Println("No holders found")
			fmt.Println()
			continue
		}

		// Convert to WhaleInfo for easier sorting
		var whales []WhaleInfo
		for _, holder := range marketHolder.Holders {
			whales = append(whales, WhaleInfo{
				Address:               holder.ProxyWallet,
				Amount:                holder.Amount,
				Outcome:               "", // Will be determined by outcome index
				OutcomeIndex:          holder.OutcomeIndex,
				Name:                  holder.Name,
				Pseudonym:             holder.Pseudonym,
				DisplayUsernamePublic: holder.DisplayUsernamePublic,
			})
		}

		// Sort by amount (largest first)
		sort.Slice(whales, func(i, j int) bool {
			return whales[i].Amount.GreaterThan(whales[j].Amount)
		})

		// Calculate concentration
		totalAmount := decimal.Zero
		for _, whale := range whales {
			totalAmount = totalAmount.Add(whale.Amount)
		}

		// Show top 10 whales
		fmt.Println("Top 10 Whales:")
		for i, whale := range whales {
			if i >= 10 {
				break
			}

			displayName := whale.Address[:10] + "..."
			if whale.DisplayUsernamePublic {
				if whale.Name != "" {
					displayName = whale.Name
				} else if whale.Pseudonym != "" {
					displayName = whale.Pseudonym
				}
			}

			percentage := whale.Amount.Div(totalAmount).Mul(decimal.NewFromInt(100))

			fmt.Printf("  #%d: %s\n", i+1, displayName)
			fmt.Printf("      Address: %s\n", whale.Address)
			fmt.Printf("      Amount: %s tokens (%.2f%% of total)\n",
				whale.Amount.StringFixed(2),
				percentage.InexactFloat64(),
			)

			// Get recent activity for this whale
			fmt.Println("      Recent Activity:")
			showWhaleActivity(client, whale.Address, marketId, 3)
			fmt.Println()
		}

		// Calculate concentration metrics
		top3Amount := decimal.Zero
		top10Amount := decimal.Zero
		for i, whale := range whales {
			if i < 3 {
				top3Amount = top3Amount.Add(whale.Amount)
			}
			if i < 10 {
				top10Amount = top10Amount.Add(whale.Amount)
			}
		}

		top3Pct := top3Amount.Div(totalAmount).Mul(decimal.NewFromInt(100))
		top10Pct := top10Amount.Div(totalAmount).Mul(decimal.NewFromInt(100))

		fmt.Printf("\n=== Concentration Metrics ===\n")
		fmt.Printf("Top 3 holders control: %.2f%% of tokens\n", top3Pct.InexactFloat64())
		fmt.Printf("Top 10 holders control: %.2f%% of tokens\n", top10Pct.InexactFloat64())

		if top3Pct.GreaterThan(decimal.NewFromInt(50)) {
			fmt.Println("⚠️  WARNING: High concentration! Top 3 holders control >50% of tokens")
		}
		if top10Pct.GreaterThan(decimal.NewFromInt(70)) {
			fmt.Println("⚠️  WARNING: Very high concentration! Top 10 holders control >70% of tokens")
		}

		fmt.Println()
	}
}

func showWhaleActivity(client *polymarketdata.DataClient, address string, marketId string, limit int) {
	activities, err := client.GetActivity(&polymarketdata.GetActivityParams{
		User:          address,
		Market:        []string{marketId},
		Limit:         limit,
		SortBy:        polymarketdata.ActivitySortByTimestamp,
		SortDirection: polymarketdata.SortDirectionDesc,
	})
	if err != nil {
		fmt.Printf("        Failed to get activity: %v\n", err)
		return
	}

	if len(activities) == 0 {
		fmt.Println("        No recent activity in this market")
		return
	}

	for _, activity := range activities {
		actionDesc := string(activity.Type)
		if activity.Type == polymarketdata.ActivityTypeTrade {
			actionDesc = fmt.Sprintf("%s %s", activity.Side, activity.Type)
		}

		fmt.Printf("        - [%s] Size: %s",
			actionDesc,
			activity.Size.StringFixed(2),
		)

		if activity.Type == polymarketdata.ActivityTypeTrade {
			fmt.Printf(" @ $%s (Total: $%s)",
				activity.Price.StringFixed(3),
				activity.UsdcSize.StringFixed(2),
			)
		}
		fmt.Println()
	}
}
