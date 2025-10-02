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

var Env *Environment

type Game struct {
	Lander              *Lander
	TickLimit           int
	TickElapsed         int
	screenshotRequested bool
	crashed             bool
	won                 bool
	paused              bool
	Score               float64
	prevDistance        float64 // Track previous distance for reward calculation
	prevSpeed           float64 // Track previous speed for reward calculation
	hasLanded           bool    // Track if legs have touched ground this episode
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
	Env.Update()
	g.Lander.Update(Env)
	g.TickElapsed++

	// Check for landing/crash
	if g.Lander.Crashed && !g.paused {
		g.crashed = true
		g.paused = true
		g.Score -= 100
		g.Lander.VelocityX = 0
		g.Lander.VelocityY = 0
		return nil
	}

	// Check if it's on the landing pad
	onPad := IsOnLandingPad(g.Lander.X)

	// Check if landing conditions are safe
	safe := g.Lander.SafeToLand()

	if IsLanderOnGround(g.Lander.Y) {
		if safe && onPad {
			// Safe landing!
			g.won = true
			g.paused = true
			g.Score += 100
			g.Lander.VelocityX = 0
			g.Lander.VelocityY = 0
		} else {
			// Crash - either too fast, wrong angle, or off pad
			g.crashed = true
			g.paused = true
			g.Score -= 100
			g.Lander.VelocityX = 0
			g.Lander.VelocityY = 0
		}
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
			g.hasLanded = false
			// Initialize distance and speed tracking
			targetX := (LandingPadLeft + LandingPadRight) / 2.0
			targetY := GroundLevel - LanderBottomOffset
			g.prevDistance = math.Sqrt(math.Pow(g.Lander.X-targetX, 2) + math.Pow(g.Lander.Y-targetY, 2))
			g.prevSpeed = 0
		}
	} else {
		// If manually paused, only space unpauses
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.paused = false
		}
	}
	return nil
}

func (g *Game) checkOffScreen() bool {
	return g.Lander.X < -100 || g.Lander.X > 900 || g.Lander.Y >= 700 || g.Lander.Y < -100
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	Env.Draw(screen)
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
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	Env = NewEnvironment()

	// Initialize game state
	initialLander := &Lander{X: 390, Y: 0}
	initialDistance := Env.Distance(initialLander)

	game := &Game{
		Lander:       initialLander,
		TickLimit:    1000,
		Score:        0,
		prevDistance: initialDistance,
		prevSpeed:    0,
		hasLanded:    false,
	}
	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
