package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GameState represents the state of the Lunar Lander environment.
// LanderX and LanderY represent the CENTER point of the lander
type GameState struct {
	LanderX    float64
	LanderY    float64
	VelocityX  float64
	VelocityY  float64
	Angle      float64
	IsDoneFlag bool
}

// Environment represents the landing environment with peaks.
type Environment struct {
	Peaks       []Triangle
	TargetX     float64 // X-coordinate of the center of the landing pad
	TargetY     float64 // Y-coordinate of the landing pad
	TargetWidth float64 // Width of the landing pad
}

type Triangle struct {
	X1, Y1 float64
	X2, Y2 float64
	X3, Y3 float64
}

// NewEnvironment creates a new environment with predefined peaks.
func NewEnvironment() *Environment {
	return &Environment{
		Peaks: []Triangle{
			// Left peaks
			{X1: 0, Y1: GroundLevel, X2: 100, Y2: GroundLevel - 100, X3: 200, Y3: GroundLevel},
			{X1: 200, Y1: GroundLevel, X2: 250, Y2: GroundLevel - 50, X3: LandingPadLeft, Y3: GroundLevel},
			// Right peaks
			{X1: LandingPadRight, Y1: GroundLevel, X2: 550, Y2: GroundLevel - 30, X3: 600, Y3: GroundLevel},
			{X1: 600, Y1: GroundLevel, X2: 700, Y2: GroundLevel - 50, X3: 800, Y3: GroundLevel},
		},
		TargetX:     (LandingPadLeft + LandingPadRight) / 2.0,
		TargetY:     GroundLevel,
		TargetWidth: LandingPadRight - LandingPadLeft,
	}
}

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
		newState.VelocityY += Gravity
	case 1: // Fire left orientation engine
		newState.Angle -= SideThrust
	case 2: // Fire main engine
		newState.VelocityX += math.Sin(newState.Angle) * MainThrust
		newState.VelocityY -= math.Cos(newState.Angle) * MainThrust
	case 3: // Fire right orientation engine
		newState.Angle += SideThrust
	}

	// Gravity always applies (even when thrusting)
	newState.VelocityY += Gravity
	newState.LanderY += newState.VelocityY
	newState.LanderX += newState.VelocityX

	// Check if the lander has hit the ground
	if IsLanderOnGround(newState.LanderY) {
		// Snap to ground level (center point)
		newState.LanderY = GroundLevel - LanderBottomOffset
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

func (e *Environment) Update() {
	// Placeholder for dynamic updates to the environment
	// For now, no dynamic updates are needed
}

func (e *Environment) Draw(screen *ebiten.Image) {
	// Draw the flat landing area
	ebitenutil.DrawLine(screen, LandingPadLeft, GroundLevel, LandingPadRight, GroundLevel, color.White)

	// Draw flags at the landing pad boundaries
	ebitenutil.DrawLine(screen, LandingPadLeft, GroundLevel, LandingPadLeft, GroundLevel-20, color.White)
	ebitenutil.DrawLine(screen, LandingPadRight, GroundLevel, LandingPadRight, GroundLevel-20, color.White)

	// Draw peaks from the environment
	for _, peak := range e.Peaks {
		ebitenutil.DrawLine(screen, peak.X1, peak.Y1, peak.X2, peak.Y2, color.White)
		ebitenutil.DrawLine(screen, peak.X2, peak.Y2, peak.X3, peak.Y3, color.White)
		ebitenutil.DrawLine(screen, peak.X3, peak.Y3, peak.X1, peak.Y1, color.White)
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
	if IsLanderOnGround(g.LanderY) {
		if g.IsSafeLanding() && IsOnLandingPad(g.LanderX) {
			return "Safe Landing"
		}
		return "Crash"
	}
	return "In Air"
}

// CheckCollision checks if a point (x, y) is inside any triangle in the environment.
func (e *Environment) CheckCollision(x, y float64) bool {
	// Check if a point (x, y) is inside any triangle
	for _, peak := range e.Peaks {
		if pointInTriangle(x, y, peak) {
			return true
		}
	}
	return false
}

func pointInTriangle(px, py float64, t Triangle) bool {
	// Barycentric technique to check if a point is inside a triangle
	area := 0.5 * (-t.Y2*t.X3 + t.Y1*(-t.X2+t.X3) + t.X1*(t.Y2-t.Y3) + t.X2*t.Y3)
	s := 1 / (2 * area) * (t.Y1*t.X3 - t.X1*t.Y3 + (t.Y3-t.Y1)*px + (t.X1-t.X3)*py)
	tCoord := 1 / (2 * area) * (t.X1*t.Y2 - t.Y1*t.X2 + (t.Y1-t.Y2)*px + (t.X2-t.X1)*py)

	return s > 0 && tCoord > 0 && (s+tCoord) < 1
}

// Distance calculates the Euclidean distance from the lander to the center of the landing pad
func (e *Environment) Distance(lander *Lander) float64 {
	dx := lander.X - e.TargetX
	dy := lander.Y - e.TargetY
	return math.Sqrt(dx*dx + dy*dy)
}
