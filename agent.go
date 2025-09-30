package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

type Node struct {
	state       *GameState
	action      int
	visitCount  int
	totalReward float64
	children    []*Node
	parent      *Node
}

type Tree struct {
	Root        *Node
	Simulations int
}

type Agent struct {
	Tree *Tree
}

func NewAgent(initialState *GameState) *Agent {
	return &Agent{
		Tree: &Tree{
			Root: &Node{
				state:       initialState,
				visitCount:  0,
				totalReward: 0,
				children:    nil,
				parent:      nil,
			},
			Simulations: 0,
		},
	}
}

func (a *Agent) SelectAction() int {
	// Perform MCTS to select the best action
	for i := 0; i < 1000; i++ { // Number of simulations
		node := a.treePolicy(a.Tree.Root)
		reward := a.simulate(node.state)
		a.backpropagate(node, reward)
	}

	// Choose the best action based on visit count
	bestAction := -1
	maxVisits := -1
	for _, child := range a.Tree.Root.children {
		if child.visitCount > maxVisits {
			maxVisits = child.visitCount
			bestAction = child.action
		}
	}
	return bestAction
}

func (a *Agent) treePolicy(node *Node) *Node {
	// Expand the tree or select the best child
	if len(node.children) == 0 {
		return a.expand(node)
	}
	return a.bestChild(node, true)
}

func (a *Agent) expand(node *Node) *Node {
	// Generate all possible actions
	for action := 0; action < 4; action++ {
		newState := node.state.Step(action)
		child := &Node{
			state:       newState,
			action:      action,
			visitCount:  0,
			totalReward: 0,
			children:    nil,
			parent:      node,
		}
		node.children = append(node.children, child)
	}
	// Return a random child for now
	return node.children[rand.Intn(len(node.children))]
}

func (a *Agent) bestChild(node *Node, useExploration bool) *Node {
	// Use UCB1 to select the best child
	bestChild := node.children[0]
	bestValue := -math.MaxFloat64
	for _, child := range node.children {
		value := child.totalReward / float64(child.visitCount+1)
		if useExploration {
			exploration := math.Sqrt(2 * math.Log(float64(node.visitCount+1)) / float64(child.visitCount+1))
			value += exploration
		}
		if value > bestValue {
			bestValue = value
			bestChild = child
		}
	}
	return bestChild
}

func (a *Agent) simulate(state *GameState) float64 {
	// Simulate a random rollout
	simulatedState := state.Copy()
	totalReward := 0.0
	for i := 0; i < 100; i++ { // Limit the simulation depth
		if simulatedState.IsDone() {
			break
		}
		action := rand.Intn(4)
		simulatedState = simulatedState.Step(action) // Handle single return value
		totalReward += 1.0                           // Example reward increment (adjust as needed)
	}
	return totalReward
}

func (a *Agent) backpropagate(node *Node, reward float64) {
	// Update the node and its ancestors
	for node != nil {
		node.visitCount++
		node.totalReward += reward
		node = node.parent
	}
}

func (a *Agent) SaveTreeToFile(filename string) error {
	// Serialize the tree to a file (e.g., JSON format)

	// if filename == "" then filename = date_time_rand4
	if filename == "" {
		date_str := time.Now().Format("20060102_150405")
		rand_suffix := rand.Intn(1000)
		filename = date_str + "_" + fmt.Sprintf("%03d", rand_suffix)
	}

	fileName := filename + ".json"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(a.Tree)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) LoadTreeFromFile(filename string) error {
	// Deserialize the tree from a file
	fileName := filename + ".json"
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(a.Tree)
	if err != nil {
		return err
	}
	return nil
}
