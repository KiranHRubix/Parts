package main

import (
	"database/sql"
	"fmt"
	"sort"

	_ "github.com/mattn/go-sqlite3"
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

type TokenList struct {
	Name  string
	Value float64
}

type TokenInfo struct {
	TokenID       string
	ParentTokenID string
	TokenValue    float64
	DID           string
	TokenStatus   int
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

func checkForPartTokens() ([]TokenList, error) {
	db, err := sql.Open("sqlite3", "rubixtest.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare SQL statement
	stmt, err := db.Prepare("SELECT token_id, token_value FROM TokensTable WHERE token_value BETWEEN ? AND ? AND token_value NOT IN (0, 1) ORDER BY token_value ASC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(0, 1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process rows
	var tokenList []TokenList
	for rows.Next() {
		var tokenID string
		var tokenValue float64
		if err := rows.Scan(&tokenID, &tokenValue); err != nil {
			return nil, err
		}
		tokenList = append(tokenList, TokenList{Name: tokenID, Value: tokenValue})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Sort the slice based on token values
	sort.Slice(tokenList, func(i, j int) bool {
		return tokenList[i].Value < tokenList[j].Value
	})

	return tokenList, nil
}

func getTokensForTransfer(tokenList []TokenList, amount float64) {
	var selectedTokens []TokenList
	// Step 2: Iterate through the sorted token list
	for i := 0; i < len(tokenList); i++ {
		token := tokenList[i]
		// Step 3: Check if the token value is less than or equal to the remaining amount
		if token.Value <= amount {
			// Step 4: Subtract the token value from the remaining amount
			amount -= token.Value
			// remove the selected tokens from the list and store the token details in a map
			selectedTokens = append(selectedTokens, token)
			tokenList = append(tokenList[:i], tokenList[i+1:]...)
			fmt.Printf("Selected Token: %s, Token Value: %.2f\n, amount : %.3f \n", token.Name, token.Value, amount)
			fmt.Println("Token List:", tokenList)
			fmt.Println("Selected Tokens:", selectedTokens)
			if amount == 0 {
				fmt.Println("Amount is 0, we got the required tokens")
			} else {
				fmt.Println("Amount is not 0, we need to select more tokens")
				//check whether the remaining amount is greater than the smallest token value
				if amount < tokenList[0].Value {
					fmt.Println("Amount is less than the smallest token value, we need to select more tokens")
					fmt.Printf("Remaining Amount: %.3f\n ", amount, "Split the tokens")
				}
			}
			// store the token details in a map
		}
	}

}

func main() {
	// Example usage
	tokenList, err := checkForPartTokens()
	getTokensForTransfer(tokenList, 0.597)
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
