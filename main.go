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
		if len(pathForValue) == 0 {
			result.RemainingAmount += value
		} else {
			result.PartToken[i].Path = fmt.Sprint(pathForValue)
			result.PartToken[i].Value = value
			visitedNew[result.PartToken[i].Path] = true
		}
	}

	result.Status = true
	return result, visitedNew, nil
}

/* func main() {
	root := CreateTree(1.0, 7) // Creates a tree with root value 1, each node having 2 or 5 children depending on the depth, and a depth of 7

	//0.521
	//values := []float64{0.001, 0.001, 0.001, 0.001, 0.001, 0.001, 0.005, 0.005, 0.005, 0.005, 0.005, 0.05, 0.05, 0.05}
	//values := []float64{0.05, 0.05, 0.05, 0.005, 0.005, 0.005, 0.005, 0.005, 0.001, 0.001, 0.001, 0.001, 0.001, 0.001}
	values := []float64{0.001, 0.005, 0.01, 0.05, 0.5}
	visited := make(map[string]bool)
	visited2 := make(map[string]bool)

	fmt.Println("using DFS")
	for _, value := range values {
		path := FindPathDFS(root, value, []int{}, visited)
		visited2[formatPath(path)] = true
		fmt.Printf("Path to %f: %v\n", value, path)
	}

	//visited2 := make(map[string]bool)
	values2 := []float64{0.001, 0.05, 0.5}
	fmt.Println("using DFS and new values but old visited")
	for _, value := range values2 {
		path := FindPathDFS(root, value, []int{}, visited2)
		if len(path) == 0 {
			fmt.Printf("no path to value %v in this tree \n", value)
		}
		fmt.Printf("Path to %f: %v\n", value, path)
	}
} */

func main() {
	amount := 0.011
	requiredDenomination, err := SelectTokenPerLevel(amount)
	if err != nil {
		fmt.Println("Error:", err)
	}
	//fmt.Println(requiredDenomination.RequiredDenomination)

	visitedOld := make(map[string]bool, 0)
	TokenGenerationResult, visitedNew, err := GenerateToken(requiredDenomination.RequiredDenomination, "ParentToken", 1.0, visitedOld)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("VIsited NEw", visitedNew)
	for _, PartToken := range TokenGenerationResult.PartToken {
		fmt.Println("PArtToken", PartToken)
	}

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
}
