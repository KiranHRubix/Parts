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
	//Level int
	Path  string
	Value float64
}

var Denominations = []float64{1.0, 0.500, 0.100, 0.050, 0.010, 0.005, 0.001}

type Result struct {
	RequiredDenomination map[float64]int
}

func SelectTokenPerLevel(amount float64) (Result, error) {
	result := Result{}
	result.RequiredDenomination = make(map[float64]int, 0)

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

		//fmt.Println("After considering denomination", denomination, "Remaining Amount:", remainingAmount)
		// Step 3: Check how many times the denomination can be subtracted from the remaining value without exceeding it

		for remainingAmount.Cmp(denomination) == 1 || remainingAmount.Cmp(denomination) == 0 {
			remainingAmount = remainingAmount.Sub(denomination)

			// Step 4: Update the count for the selected denomination
			/* if result.RequiredDenomination[i] == nil {
				result.RequiredDenomination[i] = make(map[float64]int)
			} */
			result.RequiredDenomination[denominationFloat]++
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

/*
*
generateToken generates token  based on the required denomination, parent token hash, and token value.
@param RequiredDenomination []map[float64]int - The required denomination for generating token hashes.
@param ParentTokenHash string - The parent token hash used for generating token hashes.
@param TokenValue float64 - The value of the token used for generating token hashes.
@return *TokenGenerationResult - The result of the token generation process.
@return error - An error if any occurred during the token generation process.
*/
func GenerateToken(RequiredDenomination map[float64]int, ParentTokenHash string, TokenValue float64, visitedOld map[string]bool) (*TokenGenerationResult, map[string]bool, error) {
	result := &TokenGenerationResult{
		Status: false,
	}

	var denominationsList []float64
	// Iterate over the keys of the RequiredDenomination map
	for denomination, count := range RequiredDenomination {
		// Append the denomination to the denominationsList 'count' number of times
		for i := 0; i < count; i++ {
			denominationsList = append(denominationsList, denomination)
		}
	}

	sort.Float64s(denominationsList)
	//fmt.Println(denominationsList)

	visitedNew := make(map[string]bool, 0)

	//fmt.Println("visitedOld", visitedOld)
	if len(visitedOld) != 0 {
		for key, value := range visitedOld {
			visitedNew[key] = value
		}
	}

	var root *Node
	if TokenValue == 1 && len(visitedOld) >= 0 {
		root = CreateTree(1.0, 7)
		//printTreeDFS(root, 1)
	}

	result.PartToken = make([]PartToken, len(denominationsList))
	for i, value := range denominationsList {
		//var path []int
		pathForValue := FindPathDFS(root, value, []int{}, visitedOld)
		pathStr := fmt.Sprint(pathForValue)
		if len(pathForValue) == 0 /* || visitedOld[pathStr] */ {
			result.RemainingAmount += value
		} else {
			result.PartToken[i].Path = fmt.Sprint(pathForValue)
			result.PartToken[i].Value = value
		}
		visitedNew[pathStr] = true

	}

	result.Status = true
	return result, visitedNew, nil
}

// To Do: storing visited path in DB ,

/* func main() {
	amount := 0.011
	requiredDenomination, err := SelectTokenPerLevel(amount)
	if err != nil {
		fmt.Println("Error:", err)
	}
	//fmt.Println(requiredDenomination.RequiredDenomination)

	visitedOld := make(map[string]bool, 0)
	TokenGenerationResult, _, err := GenerateToken(requiredDenomination.RequiredDenomination, "ParentToken", 1.0, visitedOld)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("VIsited NEw", visitedNew)
	visitedNew := make(map[string]bool, 0)
	for _, PartToken := range TokenGenerationResult.PartToken {
		fmt.Println("PArtToken", PartToken)
		visitedNew[fmt.Sprint(PartToken.Path)] = true
	}

	fmt.Println("VIsited NEw", visitedNew)
	//passing the visited new in previosu as the visited old
	fmt.Println("new set old token ")

	amount = 0.01
	requiredDenomination, err = SelectTokenPerLevel(amount)
	if err != nil {
		fmt.Println("Error:", err)
	}
	//fmt.Println(requiredDenomination.RequiredDenomination)

	//visitedOld := make(map[string]bool, 0)
	TokenGenerationResult, visitedNew2, err := GenerateToken(requiredDenomination.RequiredDenomination, "ParentToken", 1.0, visitedNew)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("VIsited NEw", visitedNew2)
	for _, PartToken := range TokenGenerationResult.PartToken {
		fmt.Println("PArtToken", PartToken)
	}
	fmt.Println("remaining amount", TokenGenerationResult.RemainingAmount)
} */
