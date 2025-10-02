package main

import (
	"testing"
)

func TestStep(t *testing.T) {
	// Test case 1: No action
	gs := &GameState{LanderY: 300}
	newState := gs.Step(0)
	if newState.VelocityY == 0 {
		t.Errorf("Expected VelocityY to increase with gravity, but it was 0")
	}

	// Test case 2: Fire main engine
	gs = &GameState{LanderY: 300}
	newState = gs.Step(2)
	if newState.VelocityY == 0 {
		t.Errorf("Expected VelocityY to decrease with main engine thrust, but it was 0")
	}
}

func TestIsSafeLanding(t *testing.T) {
	// Test case 1: Safe landing
	gs := &GameState{
		VelocityY: 1.0,
		VelocityX: 0.5,
		Angle:     0.1,
	}
	if !gs.IsSafeLanding() {
		t.Errorf("Expected a safe landing, but it was not")
	}

	// Test case 2: Unsafe landing (high vertical speed)
	gs = &GameState{
		VelocityY: 3.0,
		VelocityX: 0.5,
		Angle:     0.1,
	}
	if gs.IsSafeLanding() {
		t.Errorf("Expected an unsafe landing, but it was safe")
	}
}

func TestCheckLanding(t *testing.T) {
	// Test case 1: In air
	gs := &GameState{LanderY: 300}
	status := gs.CheckLanding()
	if status != "In Air" {
		t.Errorf("Expected status 'In Air', but got '%s'", status)
	}

	// Test case 2: Safe landing (on landing pad with safe speeds)
	gs = &GameState{
		LanderX:   400, // On the landing pad (300-500)
		LanderY:   485,
		VelocityY: 1.0,
		VelocityX: 0.5,
		Angle:     0.1,
	}
	status = gs.CheckLanding()
	if status != "Safe Landing" {
		t.Errorf("Expected status 'Safe Landing', but got '%s'", status)
	}

	// Test case 3: Crash (too fast)
	gs = &GameState{
		LanderX:   400, // On the landing pad
		LanderY:   485,
		VelocityY: 3.0, // Too fast!
		VelocityX: 0.5,
		Angle:     0.1,
	}
	status = gs.CheckLanding()
	if status != "Crash" {
		t.Errorf("Expected status 'Crash', but got '%s'", status)
	}

	// Test case 4: Crash (off landing pad)
	gs = &GameState{
		LanderX:   250, // Outside the landing pad
		LanderY:   485,
		VelocityY: 1.0,
		VelocityX: 0.5,
		Angle:     0.1,
	}
	status = gs.CheckLanding()
	if status != "Crash" {
		t.Errorf("Expected status 'Crash' (off pad), but got '%s'", status)
	}
}
