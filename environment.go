package main

import (
	"math"
)

// GameState represents the state of the Lunar Lander environment.
type GameState struct {
	LanderX    float64
	LanderY    float64
	VelocityX  float64
	VelocityY  float64
	Angle      float64
	IsDoneFlag bool
}

// Safe landing thresholds
const (
	SafeVerticalSpeed   = 2.0  // Maximum safe vertical speed
	SafeHorizontalSpeed = 1.0  // Maximum safe horizontal speed
	SafeLandingAngle    = 0.26 // Maximum safe angle in radians (~15 degrees)
)

// Step simulates the environment for a given action and returns the new state.
func (g *GameState) Step(action int) *GameState {
	newState := &GameState{
		LanderX:    g.LanderX,
		LanderY:    g.LanderY,
		VelocityX:  g.VelocityX,
		VelocityY:  g.VelocityY,
		Angle:      g.Angle,
		IsDoneFlag: g.IsDoneFlag,
	}

	switch action {
	case 0: // Do nothing
		// Gravity affects the lander
		newState.VelocityY += 0.05
	case 1: // Fire left orientation engine
		newState.Angle -= 0.1
	case 2: // Fire main engine
		newState.VelocityY -= 0.1
	case 3: // Fire right orientation engine
		newState.Angle += 0.1
	}

	newState.LanderY += newState.VelocityY
	newState.LanderX += newState.VelocityX

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
		VelocityX:  g.VelocityX,
		VelocityY:  g.VelocityY,
		Angle:      g.Angle,
		IsDoneFlag: g.IsDoneFlag,
	}
}

// IsDone checks if the game is over.
func (g *GameState) IsDone() bool {
	return g.IsDoneFlag
}

// IsSafeLanding checks if the lander is within safe landing thresholds.
func (g *GameState) IsSafeLanding() bool {
	// Check if the lander is within safe landing thresholds
	return math.Abs(g.VelocityY) <= SafeVerticalSpeed &&
		math.Abs(g.VelocityX) <= SafeHorizontalSpeed &&
		math.Abs(g.Angle) <= SafeLandingAngle
}

// CheckLanding determines if the landing is safe or a crash.
func (g *GameState) CheckLanding() string {
	// Determine if the landing is safe or a crash
	if g.LanderY >= 485 { // Assuming 485 is the ground level
		if g.IsSafeLanding() {
			return "Safe Landing"
		}
		return "Crash"
	}
	return "In Air"
}
