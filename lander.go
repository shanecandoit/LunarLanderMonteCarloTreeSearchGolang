package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Lander struct {
	X, Y        float64
	VelocityX   float64
	VelocityY   float64
	Angle       float64
	ThrustDown  int
	ThrustLeft  int
	ThrustRight int
}

func (l *Lander) Update() {
	// Space or Up arrow key for main thrust
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		l.ThrustDown = 1
		l.VelocityX += math.Sin(l.Angle) * 0.1
		l.VelocityY += math.Cos(l.Angle) * -0.1
	} else {
		l.ThrustDown = 0
	}

	// Left arrow key for left orientation engine
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		l.ThrustLeft = 1
		l.Angle -= 0.05
	} else {
		l.ThrustLeft = 0
	}

	// Right arrow key for right orientation engine
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		l.ThrustRight = 1
		l.Angle += 0.05
	} else {
		l.ThrustRight = 0
	}

	// Update lander position and velocity
	l.VelocityY += 0.05 // Gravity
	l.X += l.VelocityX
	l.Y += l.VelocityY
}

func (l *Lander) Draw(screen *ebiten.Image) {
	drawLander(screen, l.X, l.Y, l.Angle)
	if l.ThrustDown > 0 {
		// Draw flame when thrusting
		op := &ebiten.DrawImageOptions{}
		flameImage := ebiten.NewImage(6, 10)
		flameImage.Fill(color.RGBA{255, 165, 0, 255})
		op.GeoM.Translate(12, 20)   // position relative to lander image
		op.GeoM.Translate(-15, -15) // move to origin
		op.GeoM.Rotate(l.Angle)
		op.GeoM.Translate(l.X, l.Y) // move to screen
		screen.DrawImage(flameImage, op)
	}
	if l.ThrustLeft > 0 { // flame on the right
		op := &ebiten.DrawImageOptions{}
		flameImage := ebiten.NewImage(10, 4)
		flameImage.Fill(color.RGBA{255, 165, 0, 255})
		op.GeoM.Translate(25, 8)    // position relative to lander image
		op.GeoM.Translate(-15, -15) // move to origin
		op.GeoM.Rotate(l.Angle)
		op.GeoM.Translate(l.X, l.Y) // move to screen
		screen.DrawImage(flameImage, op)
	}
	if l.ThrustRight > 0 { // flame on the left
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
