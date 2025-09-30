package main

import (
	"testing"
)

func TestAgentSelectAction(t *testing.T) {
	// Create a mock initial state
	initialState := &GameState{
		LanderX:    0,
		LanderY:    0,
		VelocityY:  0,
		IsDoneFlag: false,
	}

	agent := NewAgent(initialState)
	action := agent.SelectAction()

	if action < 0 || action > 3 {
		t.Errorf("Invalid action selected: %d", action)
	}
}
