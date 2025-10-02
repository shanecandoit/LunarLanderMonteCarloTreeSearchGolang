package main

// Game world constants
const (
	// Screen dimensions
	ScreenWidth  = 800
	ScreenHeight = 600

	// Ground and landing pad
	GroundLevel     = 500.0 // Y coordinate of the ground surface
	LandingPadLeft  = 300.0
	LandingPadRight = 500.0

	// Lander dimensions
	LanderWidth         = 30.0
	LanderHeight        = 30.0
	LanderCenterOffsetX = 15.0 // Distance from left edge to center
	LanderCenterOffsetY = 15.0 // Distance from top edge to center
	LanderBottomOffset  = 15.0 // Distance from center to bottom (legs)

	// Physics constants
	Gravity    = 0.05
	MainThrust = 0.1
	SideThrust = 0.05 // Rotational thrust

	// Safe landing thresholds
	SafeVerticalSpeed   = 2.0  // Maximum safe vertical speed
	SafeHorizontalSpeed = 1.0  // Maximum safe horizontal speed
	SafeLandingAngle    = 0.26 // Maximum safe angle in radians (~15 degrees)
)

// GetLanderBottomY returns the Y coordinate of the lander's bottom (legs)
func GetLanderBottomY(landerCenterY float64) float64 {
	return landerCenterY + LanderBottomOffset
}

// IsLanderOnGround checks if the lander has touched the ground
func IsLanderOnGround(landerCenterY float64) bool {
	return GetLanderBottomY(landerCenterY) >= GroundLevel
}

// IsOnLandingPad checks if the lander is positioned over the landing pad
func IsOnLandingPad(landerCenterX float64) bool {
	return landerCenterX >= LandingPadLeft && landerCenterX <= LandingPadRight
}
