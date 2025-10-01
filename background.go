package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Triangle struct {
	X1, Y1 float64
	X2, Y2 float64
	X3, Y3 float64
}

type Background struct {
	Peaks []Triangle
}

func NewBackground() *Background {
	return &Background{
		Peaks: []Triangle{
			// Left peaks
			{X1: 0, Y1: 500, X2: 100, Y2: 400, X3: 200, Y3: 500},
			{X1: 200, Y1: 500, X2: 250, Y2: 450, X3: 300, Y3: 500},
			// Right peaks
			{X1: 500, Y1: 500, X2: 550, Y2: 470, X3: 600, Y3: 500},
			{X1: 600, Y1: 500, X2: 700, Y2: 450, X3: 800, Y3: 500},
			{X1: 800, Y1: 500, X2: 900, Y2: 400, X3: 1000, Y3: 500},
			{X1: 1000, Y1: 500, X2: 1050, Y2: 480, X3: 1100, Y3: 500},
		},
	}
}

func (b *Background) Update() {
	// No dynamic updates for the background yet
}

func (b *Background) Draw(screen *ebiten.Image) {
	// Draw the flat area
	ebitenutil.DrawLine(screen, 300, 500, 500, 500, color.White)

	// Draw flags
	ebitenutil.DrawLine(screen, 300, 500, 300, 480, color.White)
	ebitenutil.DrawLine(screen, 500, 500, 500, 480, color.White)

	// Draw peaks
	for _, peak := range b.Peaks {
		ebitenutil.DrawLine(screen, peak.X1, peak.Y1, peak.X2, peak.Y2, color.White)
		ebitenutil.DrawLine(screen, peak.X2, peak.Y2, peak.X3, peak.Y3, color.White)
		ebitenutil.DrawLine(screen, peak.X3, peak.Y3, peak.X1, peak.Y1, color.White)
	}
}

func (b *Background) CheckCollision(x, y float64) bool {
	// Check if a point (x, y) is inside any triangle
	for _, peak := range b.Peaks {
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
