package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var landerImage *ebiten.Image
var background *Background

type Game struct {
	Lander              *Lander
	TickLimit           int
	TickElapsed         int
	screenshotRequested bool
	crashed             bool
	won                 bool
	paused              bool
}

func (g *Game) Update() error {
	if g.paused {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			// Unpause and reset the game
			g.paused = false
			g.crashed = false
			g.won = false
			g.Lander = &Lander{X: 390, Y: 0}
			g.TickElapsed = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			return ebiten.Termination
		}
		return nil
	}

	// Escape key to quit
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.screenshotRequested = true
	}

	background.Update()
	g.Lander.Update()

	// Time tick logic
	g.TickElapsed++
	if g.TickLimit > 0 && g.TickElapsed >= g.TickLimit {
		return ebiten.Termination
	}

	// Create a GameState to check for landing status
	gs := &GameState{
		LanderX:   g.Lander.X,
		LanderY:   g.Lander.Y,
		VelocityX: g.Lander.VelocityX,
		VelocityY: g.Lander.VelocityY,
		Angle:     g.Lander.Angle,
	}

	landingStatus := gs.CheckLanding()
	if landingStatus == "Crash" {
		g.crashed = true
		g.paused = true
		g.Lander.Y = 485
		g.Lander.VelocityY = 0
		g.Lander.VelocityX = 0
	}
	if landingStatus == "Safe Landing" {
		g.won = true
		g.paused = true
	}

	// Off screen termination
	if g.Lander.X < -100 || g.Lander.X > 900 {
		return ebiten.Termination
	}
	if g.Lander.Y >= 700 || g.Lander.Y < -100 {
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	background.Draw(screen)
	g.Lander.Draw(screen)

	// draw thrust as bits, not booleans
	msg := fmt.Sprintf(
		"X: %4.2f, Y: %4.2f\nVelX: %4.2f, VelY: %4.2f\nAngle: %4.2f\nThrust: D:%d L:%d R:%d\nTick: %d/%d",
		g.Lander.X, g.Lander.Y, g.Lander.VelocityX, g.Lander.VelocityY, g.Lander.Angle,
		g.Lander.ThrustDown, g.Lander.ThrustLeft, g.Lander.ThrustRight, g.TickElapsed, g.TickLimit,
	)
	ebitenutil.DebugPrintAt(screen, msg, 0, 500)

	if g.crashed {
		ebitenutil.DebugPrintAt(screen, "You Crashed", 350, 300)
	}
	if g.won {
		ebitenutil.DebugPrintAt(screen, "You Won", 350, 300)
	}

	if g.screenshotRequested {
		g.saveScreenshot(screen)
	}
}

func (g *Game) saveScreenshot(screen *ebiten.Image) {
	g.screenshotRequested = false
	filename := time.Now().Format("2006.01.02_15.04.05") + ".png"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := png.Encode(file, screen); err != nil {
		log.Fatal(err)
	}
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

	background = NewBackground()
	game := &Game{
		Lander:    &Lander{X: 390, Y: 0},
		TickLimit: 1000,
	}
	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
