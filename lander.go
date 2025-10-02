package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var landerImage *ebiten.Image

func init() {
	landerImage = ebiten.NewImage(30, 30)
	// box body
	ebitenutil.DrawRect(landerImage, 5, 0, 20, 20, color.RGBA{255, 0, 255, 255})
	// legs
	ebitenutil.DrawRect(landerImage, 0, 20, 5, 10, color.RGBA{255, 0, 255, 255})
	ebitenutil.DrawRect(landerImage, 25, 20, 5, 10, color.RGBA{255, 0, 255, 255})
}

type Lander struct {
	X, Y             float64
	VelocityX        float64
	VelocityY        float64
	Angle            float64
	ThrustDown       int
	ThrustLeft       int
	ThrustRight      int
	LeftLegOnGround  bool
	RightLegOnGround bool
	Crashed          bool
}

func (l *Lander) Update(env *Environment) {
	// Up arrow key for main thrust
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		l.ThrustDown = 1
		l.VelocityX += math.Sin(l.Angle) * MainThrust
		l.VelocityY += math.Cos(l.Angle) * -MainThrust
	} else {
		l.ThrustDown = 0
	}

	// Left arrow key for left orientation engine
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		l.ThrustLeft = 1
		l.Angle -= SideThrust
	} else {
		l.ThrustLeft = 0
	}

	// Right arrow key for right orientation engine
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		l.ThrustRight = 1
		l.Angle += SideThrust
	} else {
		l.ThrustRight = 0
	}

	// Update lander position and velocity
	l.VelocityY += Gravity // Gravity
	l.X += l.VelocityX
	l.Y += l.VelocityY

	// Check for collisions with the environment
	if env.CheckCollision(l.X-15, l.Y+20) || env.CheckCollision(l.X+15, l.Y+20) {
		l.Crashed = true
	}
}

// SafeToLand checks if the lander's speed and angle are within safe landing parameters
func (l *Lander) SafeToLand() bool {
	return math.Abs(l.VelocityY) <= SafeVerticalSpeed &&
		math.Abs(l.VelocityX) <= SafeHorizontalSpeed &&
		math.Abs(l.Angle) <= SafeLandingAngle
}

func (l *Lander) Draw(screen *ebiten.Image) {
	// Draw the lander body
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-15, -15)
	op.GeoM.Rotate(l.Angle)
	op.GeoM.Translate(l.X, l.Y)
	screen.DrawImage(landerImage, op)

	// Draw thrust flames if active
	l.drawThrustFlames(screen)
}

func (l *Lander) drawThrustFlames(screen *ebiten.Image) {
	if l.ThrustDown > 0 {
		// Draw flame for main thrust
		op := &ebiten.DrawImageOptions{}
		flameImage := ebiten.NewImage(6, 10)
		flameImage.Fill(color.RGBA{255, 165, 0, 255})
		op.GeoM.Translate(12, 20)   // position relative to lander image
		op.GeoM.Translate(-15, -15) // move to origin
		op.GeoM.Rotate(l.Angle)
		op.GeoM.Translate(l.X, l.Y) // move to screen
		screen.DrawImage(flameImage, op)
	}
	if l.ThrustLeft > 0 {
		// Draw flame for left thrust
		op := &ebiten.DrawImageOptions{}
		flameImage := ebiten.NewImage(10, 4)
		flameImage.Fill(color.RGBA{255, 165, 0, 255})
		op.GeoM.Translate(25, 8)    // position relative to lander image
		op.GeoM.Translate(-15, -15) // move to origin
		op.GeoM.Rotate(l.Angle)
		op.GeoM.Translate(l.X, l.Y) // move to screen
		screen.DrawImage(flameImage, op)
	}
	if l.ThrustRight > 0 {
		// Draw flame for right thrust
		op := &ebiten.DrawImageOptions{}
		flameImage := ebiten.NewImage(10, 4)
		flameImage.Fill(color.RGBA{255, 165, 0, 255})
		op.GeoM.Translate(-5, 8)    // position relative to lander image
		op.GeoM.Translate(-15, -15) // move to origin
		op.GeoM.Rotate(l.Angle)
		op.GeoM.Translate(l.X, l.Y) // move to screen
		screen.DrawImage(flameImage, op)
	}
}

func (l *Lander) checkCollisionWithPeaks(peaks []float64) {
	// Reset leg states
	l.LeftLegOnGround = false
	l.RightLegOnGround = false

	landerLeftX := l.X - 15  // Left edge of the lander
	landerRightX := l.X + 15 // Right edge of the lander

	for x, peakY := range peaks {
		if float64(x) >= landerLeftX && float64(x) <= landerRightX {
			if l.Y+20 >= peakY { // Check if legs are touching the peak
				if float64(x) < l.X {
					l.LeftLegOnGround = true
				} else {
					l.RightLegOnGround = true
				}

				// If either leg touches the peak, the lander crashes
				l.Crashed = true
			}
		}
	}
}
