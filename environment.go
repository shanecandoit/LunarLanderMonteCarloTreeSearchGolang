package main

// GameState represents the state of the Lunar Lander environment.
type GameState struct {
	LanderX    float64
	LanderY    float64
	VelocityY  float64
	IsDoneFlag bool
}

// Step simulates the environment for a given action and returns the new state.
func (g *GameState) Step(action int) *GameState {
	newState := &GameState{
		LanderX:    g.LanderX,
		LanderY:    g.LanderY,
		VelocityY:  g.VelocityY,
		IsDoneFlag: g.IsDoneFlag,
	}

	switch action {
	case 0: // Do nothing
		// Gravity affects the lander
		newState.VelocityY += 0.05
	case 1: // Fire left orientation engine
		// No-op for now (can be extended later)
	case 2: // Fire main engine
		newState.VelocityY -= 0.1
	case 3: // Fire right orientation engine
		// No-op for now (can be extended later)
	}

	newState.LanderY += newState.VelocityY

	// Check if the lander has hit the ground
	if newState.LanderY > 400 {
		newState.LanderY = 400
		newState.VelocityY = 0
		newState.IsDoneFlag = true
	}

	return newState
}

// Copy creates a deep copy of the current game state.
func (g *GameState) Copy() *GameState {
	return &GameState{
		LanderX:    g.LanderX,
		LanderY:    g.LanderY,
		VelocityY:  g.VelocityY,
		IsDoneFlag: g.IsDoneFlag,
	}
}

// IsDone checks if the game is over.
func (g *GameState) IsDone() bool {
	return g.IsDoneFlag
}
