package main

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

const (
	LevelOne int = iota + 1
	LevelTwo
	LevelThree
	LevelFour
	LevelFive
	LevelSix
	LevelSeven
)

// Token represents a node in the tree
type Token struct {
	Level int
	Path  string
	Value float64
}

var Denominations = []float64{1.0, 0.500, 0.100, 0.050, 0.010, 0.005, 0.001}

type Result struct {
	RequiredDenomination []map[float64]int
}

func selectTokenPerLevel(amount float64) (Result, error) {
	result := Result{}
	result.RequiredDenomination = make([]map[float64]int, len(Denominations))

	// Step 1: Sort denominations in descending order
	sort.Sort(sort.Reverse(sort.Float64Slice(Denominations)))

	remainingAmount := decimal.NewFromFloat(amount)

	// Step 2: Iterate through the sorted denominations
	for i := 0; i < len(Denominations); i++ {
		denominationFloat := Denominations[i]
		denomination := decimal.NewFromFloat(Denominations[i])

		// Skip denominations larger than the remaining amount

		compare1 := denomination.Cmp(remainingAmount)

		if compare1 == 1 /* if denomination > remainingAmount */ {
			//fmt.Println("Skipping Level", i+1)
			continue
		}

		fmt.Println("After considering denomination", denomination, "Remaining Amount:", remainingAmount)
		// Step 3: Check how many times the denomination can be subtracted from the remaining value without exceeding it

		for remainingAmount.Cmp(denomination) == 1 || remainingAmount.Cmp(denomination) == 0 {
			remainingAmount = remainingAmount.Sub(denomination)

			// Step 4: Update the count for the selected denomination
			if result.RequiredDenomination[i] == nil {
				result.RequiredDenomination[i] = make(map[float64]int)
			}
			result.RequiredDenomination[i][denominationFloat]++
		}
	}

	// Step 5: Check if the remaining amount is zero
	if remainingAmount.Cmp(decimal.New(0, 0)) == 1 {
		return Result{}, fmt.Errorf("unable to represent the amount with the available denominations")
	}
	return result, nil
}

func main() {
	// Example usage
	amountToTransfer := 0.597
	fmt.Printf("Amount to Transfer: %.3f\n", amountToTransfer)

	result, err := selectTokenPerLevel(amountToTransfer)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Selected Denominations and Counts:")
		for i, m := range result.RequiredDenomination {
			if m != nil {
				for denom, count := range m {
					fmt.Printf("Level %d - Denomination: %.3f, Count: %d\n", i+1, denom, count)
				}
			}
		}
	}
}
