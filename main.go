package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var landerImage *ebiten.Image

type Game struct {
	landerX, landerY     float64
	velocityX, velocityY float64
	angle                float64
	thrustDown           int
	thrustLeft           int
	thrustRight          int
	tickLimit            int
	tickElapsed          int
}

func (g *Game) Update() error {
	// Space or Up arrow key for main thrust
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.thrustDown = 1
		g.velocityX += math.Sin(g.angle) * 0.1
		g.velocityY += math.Cos(g.angle) * -0.1
	} else {
		g.thrustDown = 0
	}

	// Left arrow key for left orientation engine
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.thrustLeft = 1
		g.angle -= 0.05
	} else {
		g.thrustLeft = 0
	}

	// Right arrow key for right orientation engine
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.thrustRight = 1
		g.angle += 0.05
	} else {
		g.thrustRight = 0
	}

	// Time tick logic
	g.tickElapsed++
	if g.tickLimit > 0 && g.tickElapsed >= g.tickLimit {
		return ebiten.Termination
	}

	// Update lander position and velocity
	g.velocityY += 0.05 // Gravity
	g.landerX += g.velocityX
	g.landerY += g.velocityY

	if g.landerY > 400 {
		g.landerY = 400
		g.velocityY = 0
		g.velocityX = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	drawLander(screen, g.landerX, g.landerY, g.angle)
	if g.thrustDown > 0 {
		// Draw flame when thrusting
		op := &ebiten.DrawImageOptions{}
		flameImage := ebiten.NewImage(6, 10)
		flameImage.Fill(color.RGBA{255, 165, 0, 255})
		op.GeoM.Translate(12, 20) // position relative to lander image
		op.GeoM.Translate(-15, -15) // move to origin
		op.GeoM.Rotate(g.angle)
		op.GeoM.Translate(g.landerX, g.landerY) // move to screen
		screen.DrawImage(flameImage, op)
	}
	if g.thrustLeft > 0 { // flame on the right
		op := &ebiten.DrawImageOptions{}
		flameImage := ebiten.NewImage(10, 4)
		flameImage.Fill(color.RGBA{255, 165, 0, 255})
		op.GeoM.Translate(25, 8) // position relative to lander image
		op.GeoM.Translate(-15, -15) // move to origin
		op.GeoM.Rotate(g.angle)
		op.GeoM.Translate(g.landerX, g.landerY) // move to screen
		screen.DrawImage(flameImage, op)
	}
	if g.thrustRight > 0 { // flame on the left
		op := &ebiten.DrawImageOptions{}
		flameImage := ebiten.NewImage(10, 4)
		flameImage.Fill(color.RGBA{255, 165, 0, 255})
		op.GeoM.Translate(-5, 8) // position relative to lander image
		op.GeoM.Translate(-15, -15) // move to origin
		op.GeoM.Rotate(g.angle)
		op.GeoM.Translate(g.landerX, g.landerY) // move to screen
		screen.DrawImage(flameImage, op)
	}

	// draw thrust as bits, not booleans
	msg := fmt.Sprintf(
		"X: %4.2f, Y: %4.2f\nVelX: %4.2f, VelY: %4.2f\nAngle: %4.2f\nThrust: D:%d L:%d R:%d\nTick: %d/%d",
		g.landerX, g.landerY, g.velocityX, g.velocityY, g.angle,
		g.thrustDown, g.thrustLeft, g.thrustRight, g.tickElapsed, g.tickLimit,
	)
	ebitenutil.DebugPrintAt(screen, msg, 0, 500)

}

func drawLander(screen *ebiten.Image, landerX, landerY, angle float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-15, -15)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(landerX, landerY)
	screen.DrawImage(landerImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	landerImage = ebiten.NewImage(30, 30)
	// box body
	ebitenutil.DrawRect(landerImage, 5, 0, 20, 20, color.RGBA{255, 0, 255, 255})
	// legs
	ebitenutil.DrawRect(landerImage, 0, 20, 5, 10, color.RGBA{255, 0, 255, 255})
	ebitenutil.DrawRect(landerImage, 25, 20, 5, 10, color.RGBA{255, 0, 255, 255})

	game := &Game{landerX: 390, landerY: 0, tickLimit: 1000}
	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
