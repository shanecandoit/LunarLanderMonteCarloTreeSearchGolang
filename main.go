package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var background *Background

type Game struct {
	Lander              *Lander
	TickLimit           int
	TickElapsed         int
	screenshotRequested bool
	crashed             bool
	won                 bool
	paused              bool
	Score               float64
}

func (g *Game) Update() error {
	// Handle input
	if g.paused {
		return g.handlePausedInput()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.screenshotRequested = true
	}

	// Update game state
	background.Update()
	g.Lander.Update()
	g.TickElapsed++

	// Calculate score
	// Proximity to the landing pad
	distance := math.Sqrt(math.Pow(g.Lander.X-390, 2) + math.Pow(g.Lander.Y-485, 2))
	g.Score -= distance * 0.01

	// Speed
	speed := math.Sqrt(math.Pow(g.Lander.VelocityX, 2) + math.Pow(g.Lander.VelocityY, 2))
	g.Score -= speed * 0.01

	// Angle
	g.Score -= math.Abs(g.Lander.Angle) * 0.01

	// Engine Usage
	if g.Lander.ThrustDown > 0 {
		g.Score -= 0.3
	}
	if g.Lander.ThrustLeft > 0 || g.Lander.ThrustRight > 0 {
		g.Score -= 0.03
	}

	// Check game logic
	if g.TickLimit > 0 && g.TickElapsed >= g.TickLimit {
		return ebiten.Termination
	}
	if g.checkLandingStatus() {
		return nil
	}
	if g.checkOffScreen() {
		return ebiten.Termination
	}

	return nil
}

func (g *Game) handlePausedInput() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if g.won || g.crashed {
		// If game is over, any key (other than escape) resets the game
		if len(inpututil.AppendJustPressedKeys(nil)) > 0 {
			g.paused = false
			g.crashed = false
			g.won = false
			g.Lander = &Lander{X: 390, Y: 0}
			g.TickElapsed = 0
			g.Score = 0
		}
	} else {
		// If manually paused, only space unpauses
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.paused = false
		}
	}
	return nil
}

func (g *Game) checkLandingStatus() bool {
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
		g.Score -= 100
		return true
	}
	if landingStatus == "Safe Landing" {
		g.won = true
		g.paused = true
		g.Score += 100
		return true
	}
	return false
}

func (g *Game) checkOffScreen() bool {
	return g.Lander.X < -100 || g.Lander.X > 900 || g.Lander.Y >= 700 || g.Lander.Y < -100
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	background.Draw(screen)
	g.Lander.Draw(screen)

	// draw thrust as bits, not booleans
	msg := fmt.Sprintf(
		"X: %4.2f, Y: %4.2f\nVelX: %4.2f, VelY: %4.2f\nAngle: %4.2f\nThrust: D:%d L:%d R:%d\nTick: %d/%d\nScore: %4.2f",
		g.Lander.X, g.Lander.Y, g.Lander.VelocityX, g.Lander.VelocityY, g.Lander.Angle,
		g.Lander.ThrustDown, g.Lander.ThrustLeft, g.Lander.ThrustRight, g.TickElapsed, g.TickLimit, g.Score,
	)
	ebitenutil.DebugPrintAt(screen, msg, 0, 500)

	if g.crashed {
		ebitenutil.DebugPrintAt(screen, "You Crashed", 350, 300)
	} else if g.won {
		ebitenutil.DebugPrintAt(screen, "You Won", 350, 300)
	} else if g.paused {
		ebitenutil.DebugPrintAt(screen, "Paused", 350, 300)
	}

	if g.screenshotRequested {
		g.screenshotRequested = false
		img := ebiten.NewImageFromImage(screen)
		go g.saveScreenshot(img)
	}
}

func (g *Game) saveScreenshot(screen *ebiten.Image) {
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	ebiten.SetWindowResizable(true)

	background = NewBackground()
	game := &Game{
		Lander:    &Lander{X: 390, Y: 0},
		TickLimit: 1000,
		Score:     0,
	}
	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
