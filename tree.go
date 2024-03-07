package main

import "fmt"

type Node struct {
	ID       int
	Value    float64
	Children []*Node
}

func NewNode(id int, value float64) *Node {
	return &Node{
		ID:       id,
		Value:    value,
		Children: []*Node{},
	}
}

func (n *Node) AddChild(id int, value float64) *Node {
	child := NewNode(id, value)
	n.Children = append(n.Children, child)
	return child
}

func CreateTree(rootValue float64, depth int) *Node {
	root := NewNode(1, rootValue)
	createChildren(root, depth-1, 2, 0.5)
	return root
}

func createChildren(parent *Node, depth, numChildren int, value float64) {
	if depth == 0 {
		return
	}
	for i := 0; i < numChildren; i++ {
		child := parent.AddChild(i+1, value)
		if depth%2 != 0 {
			createChildren(child, depth-1, 2, value/2)
		} else {
			createChildren(child, depth-1, 5, value/5)
		}
	}
}

func printTree(node *Node, level int) {
	if node == nil {
		return
	}
	fmt.Printf("Level %d, Child %d: Value: %f\n", level, node.ID, node.Value)
	for _, child := range node.Children {
		printTree(child, level+1)
	}
}
func FindPath(node *Node, value float64, path []int, visited map[string]bool) []int {
	if node == nil {
		return nil
	}

	// Add the current node's ID to the path
	path = append(path, node.ID)

	// Convert the path to a string to use as a key in the visited map
	pathStr := fmt.Sprint(path)

	// If the current node's value matches the target value and the path has not been visited before, return the path
	if node.Value == value && !visited[pathStr] {
		visited[pathStr] = true
		return path
	}

	// If the current node's value does not match the target value or the path has been visited before, continue the search with the node's children
	for _, child := range node.Children {
		result := FindPath(child, value, path, visited)
		if result != nil {
			// If the value is found in the subtree rooted at the child, return the result
			return result
		}
	}

	// If the value was not found in the subtree rooted at the current node, remove the node from the path and return nil
	path = path[:len(path)-1]
	return nil
}

func main() {
	root := CreateTree(1.0, 7) // Creates a tree with root value 1, each node having 2 or 5 children depending on the depth, and a depth of 7

	values := []float64{0.001, 0.01, 0.01, 0.5}
	visited := make(map[string]bool)
	for _, value := range values {
		path := FindPath(root, value, []int{}, visited)
		fmt.Printf("Path to %f: %v\n", value, path)
	}
}
