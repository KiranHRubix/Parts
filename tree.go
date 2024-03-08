package main

import (
	"fmt"
	"strings"
)

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

/* func FindPath(node *Node, value float64, path []int, visited map[string]bool) []int {
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
} */

func isSubpath(path1, path2 []int) bool {
	// Check if path1 is a subpath of path2
	//fmt.Println("path1", path1)
	//fmt.Println("path 2", path2)
	if len(path1) > len(path2) {
		return false
	}
	for i := range path1 {
		if path1[i] != path2[i] {
			return false
		}
	}
	return true
}

func FindPathDFS(node *Node, value float64, path []int, visited map[string]bool) []int {
	if node == nil {
		return nil
	}

	// Add the current node's ID to the path
	path = append(path, node.ID)

	// Convert the path to a string to use as a key in the visited map
	pathStr := fmt.Sprint(path)

	isSub := false
	// Check if the current path is a subpath of any previously visited paths
	for prevPath := range visited {
		if strings.Contains(prevPath, pathStr) {
			// If the current path is a subpath of a previously visited path, set the flag and break the loop
			isSub = true
			break
		}
	}

	//fmt.Println(isSub)
	// If the current path is a subpath of a previously visited path, add it to the visited paths
	if isSub {
		visited[pathStr] = true
	}

	// If the current node's value matches the target value and the path has not been visited before, return the path
	if node.Value == value && !visited[pathStr] {
		visited[pathStr] = true
		return path
	}

	// Continue the search with the node's children
	for _, child := range node.Children {
		result := FindPathDFS(child, value, path, visited)
		if result != nil {
			// If the value is found in the subtree rooted at the child, return the result
			return result
		}
	}

	// If the value was not found in the subtree rooted at the current node, remove the node from the path and return nil
	path = path[:len(path)-1]
	return nil
}

func parsePath(pathStr string) []int {
	// Parse the path string into a slice of integers
	var path []int

	if len(pathStr) > 0 {
		var nodeID int
		_, err := fmt.Sscanf(pathStr, "[%d", &nodeID)
		if err != nil {
			return nil
		}
		path = append(path, nodeID)

		for {
			_, err := fmt.Sscanf(pathStr, " %d", &nodeID)
			if err != nil {
				break
			}
			path = append(path, nodeID)
		}
	}

	return path
}

func FindPathBFS(root *Node, value float64) []int {
	fmt.Println("checking for value", value)
	if root == nil {
		return nil
	}

	// Create a queue for BFS and enqueue the root
	queue := []*Node{root}

	// Create a map to store paths
	paths := make(map[*Node][]int)
	paths[root] = []int{root.ID}

	// Create a map to store visited paths
	visited := make(map[string]bool)

	for len(queue) > 0 {
		// Dequeue a node from the front of the queue
		node := queue[0]
		queue = queue[1:]

		// Convert the current path to a string to use as a key in the visited map
		pathStr := fmt.Sprint(paths[node])

		// If the current path is a subpath of any previously visited paths, continue with the next node in the queue
		isSub := false
		for prevPath := range visited {
			fmt.Println("prevPath", prevPath)
			if isSubpath(paths[node], parsePath(prevPath)) {
				isSub = true
				break
			}
		}
		if isSub {
			continue
		}

		fmt.Println("isSub", isSub)

		// If the current node's value matches the target value, return the path
		if node.Value == value {
			path := paths[node]
			// Remove the last node from the visited paths
			delete(visited, pathStr)
			return path
		}

		// Add the current path to the visited paths
		visited[pathStr] = true

		// Enqueue all children of the current node
		for _, child := range node.Children {
			// Add the child to the queue
			queue = append(queue, child)

			// Add the child's path to the paths map
			paths[child] = append(paths[node], child.ID)
		}
	}

	// If the value was not found, return nil
	return nil
}

func main() {
	root := CreateTree(1.0, 7) // Creates a tree with root value 1, each node having 2 or 5 children depending on the depth, and a depth of 7

	//0.521
	values := []float64{0.001, 0.001, 0.001, 0.005, 0.005}
	visited := make(map[string]bool)
	fmt.Println("using DFS")
	for _, value := range values {
		path := FindPathDFS(root, value, []int{}, visited)
		fmt.Printf("Path to %f: %v\n", value, path)
	}

	/* fmt.Println("usinf BFS")
	for _, value := range values {
		path := FindPathBFS(root, value)
		fmt.Printf("Path to %f: %v\n", value, path)
	} */
}
