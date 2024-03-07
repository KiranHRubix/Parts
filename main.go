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
type PartToken struct {
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

type TokenGenerationResult struct {
	PartToken       []PartToken
	RemainingAmount float64
	Message         string
	Status          bool
}

/*
Starting with level 1 as root token

Even levels split parent to 2 equal parts
Odd levels split parent to 5 equal parts

Level: 1 - Node Value: 1.000
Level: 2 - Node Value: 0.500
Level: 3 - Node Value: 0.100
Level: 4 - Node Value: 0.050
Level: 5 - Node Value: 0.010
Level: 6 - Node Value: 0.005
Level: 7 - Node Value: 0.001
*/

/**


generateToken generates token  based on the required denomination, parent token hash, and token value.


@param RequiredDenomination []map[float64]int - The required denomination for generating token hashes.

@param ParentTokenHash string - The parent token hash used for generating token hashes.

@param TokenValue float64 - The value of the token used for generating token hashes.


@return *TokenGenerationResult - The result of the token generation process.

@return error - An error if any occurred during the token generation process.
*/

// This method creates the required token based on the nonbinary tree structure where we have levels from 1-7 and denominatiosn for each level
// the PartToken struct has the value path where the path from the ParentToken to the crequired denomination token would be laid out

func generateToken(RequiredDenomination []map[float64]int, ParentTokenHash string, TokenValue float64) (*TokenGenerationResult, error) {
	result := &TokenGenerationResult{
		Status:    false,
		PartToken: []PartToken{},
	}

	// Iterate over each level
	for level, denominations := range RequiredDenomination {
		// Iterate over each denomination at this level
		for denomination, count := range denominations {
			// Generate tokens for this denomination
			for i := 0; i < count; i++ {
				// Generate token hash
				//tokenHash := sha3.Sum256([]byte(fmt.Sprintf("%s%d", ParentTokenHash, i)))
				//tokenHashString := hex.EncodeToString(tokenHash[:])

				// Create a new PartToken
				partToken := PartToken{
					Level: level + 1,
					Path:  fmt.Sprintf("%s%d", ParentTokenHash, i),
					Value: denomination,
				}

				// Add PartToken to result
				result.PartToken = append(result.PartToken, partToken)
			}
		}
	}

	result.Status = true
	return result, nil
}

func main() {
	// Example usage
	amountToTransfer := 0.611
	fmt.Printf("Amount to Transfer: %.3f\n", amountToTransfer)

	reqDenom, err := selectTokenPerLevel(amountToTransfer)
	if err != nil {
		fmt.Println("Error:", err)
	}

	result, err := generateToken(reqDenom.RequiredDenomination, "TokenHash1", 1)
	if err != nil {
		fmt.Println("Error:", err)
	}
	for _, partToken := range result.PartToken {
		fmt.Println("partToken.Level", partToken.Level)
		fmt.Println("partToken.Path", partToken.Path)
		fmt.Println("partToken.Value", partToken.Value)
	}

	/* if err != nil {
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
	} */
}
