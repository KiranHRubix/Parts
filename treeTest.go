// package main

// import (
// 	"fmt"
// )

// // Node represents a node in the tree
// type Node struct {
// 	Value    float64
// 	Children []*Node
// }

// // NewNode creates a new node with the given value
// func NewNode(value float64) *Node {
// 	return &Node{
// 		Value:    value,
// 		Children: make([]*Node, 0, 5), // Initialize children slice with capacity for 5 children
// 	}
// }

// // AddChild adds a child to the node
// func (n *Node) AddChild(child *Node) {
// 	if len(n.Children) < 5 {
// 		n.Children = append(n.Children, child)
// 	} else {
// 		fmt.Println("Cannot add more children, already at maximum capacity")
// 	}
// }

// // TraverseDepthFirst traverses the tree in depth-first order
// func (n *Node) TraverseDepthFirst() {
// 	fmt.Println(n.Value)
// 	for _, child := range n.Children {
// 		child.TraverseDepthFirst()
// 	}
// }

// func BreadthFirstSearch(root *Node) {
// 	if root == nil {
// 		return
// 	}

// 	queue := []*Node{root}

// 	for len(queue) > 0 {
// 		currentNode := queue[0]
// 		queue = queue[1:]

// 		fmt.Printf("%.2f\n", currentNode.Value)

// 		for _, child := range currentNode.Children {
// 			queue = append(queue, child)
// 		}
// 	}
// }

// func main() {
// 	// Creating the root node
// 	root := NewNode(1.0)

// 	// Adding children to the root node
// 	for i := 0; i < 2; i++ {
// 		child := NewNode(0.5)
// 		for j := 0; j < 5; j++ {
// 			child.AddChild(NewNode(0.1))
// 		}
// 		root.AddChild(child)
// 	}

// 	// Traversing the tree in depth-first order
// 	fmt.Println("Depth-First Traversal:")
// 	root.TraverseDepthFirst()
// 	fmt.Println("Breadth-First Traversal:")
// 	BreadthFirstSearch(root)

// }
