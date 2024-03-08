package main

import (
	"fmt"
)

type Node struct {
	Value    float64
	Children []*Node
}

func NewNode(value float64) *Node {
	return &Node{Value: value}
}

func BreadthFirstSearch(root *Node, targetValue float64) *Node {
	if root == nil {
		return nil
	}

	queue := []*Node{root}

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if currentNode.Value == targetValue {
			return currentNode
		}

		for _, child := range currentNode.Children {
			queue = append(queue, child)
		}
	}

	return nil
}

func PrintTree(root *Node, level int) {
	if root == nil {
		return
	}

	// Print the current node
	for i := 0; i < level; i++ {
		fmt.Print("\t")
	}
	fmt.Printf("%.3f\n", root.Value)

	// Print children recursively
	for _, child := range root.Children {
		PrintTree(child, level+1)
	}
}

func main() {
	// Create the tree structure
	root := NewNode(1.0)

	// Level 2
	for i := 0; i < 2; i++ {
		child := NewNode(0.5)
		root.Children = append(root.Children, child)
	}

	// Level 3
	for i := 0; i < 2; i++ {
		child := NewNode(0.1)
		for j := 0; j < 5; j++ {
			child.Children = append(child.Children, NewNode(0.1))
		}
		root.Children[i].Children = append(root.Children[i].Children, child)
	}

	// Level 4
	for i := 0; i < 2; i++ {
		child := NewNode(0.05)
		for j := 0; j < 2; j++ {
			child.Children = append(child.Children, NewNode(0.01))
		}
		root.Children[i].Children[0].Children = append(root.Children[i].Children[0].Children, child)
	}

	// Level 5
	for i := 0; i < 2; i++ {
		child := NewNode(0.005)
		for j := 0; j < 2; j++ {
			child.Children = append(child.Children, NewNode(0.001))
		}
		root.Children[i].Children[0].Children[0].Children = append(root.Children[i].Children[0].Children[0].Children, child)
	}

	// Perform Breadth First Search to find the node with value 0.001
	fmt.Println("Breadth First Search:")
	fmt.Println(root)
	fmt.Println("Tree Structure:")
	PrintTree(root, 0)
	targetValue := 0.001
	result := BreadthFirstSearch(root, targetValue)
	if result != nil {
		fmt.Printf("Node with value %.3f found in the tree.\n", targetValue)
	} else {
		fmt.Printf("Node with value %.3f not found in the tree.\n", targetValue)
	}
}
