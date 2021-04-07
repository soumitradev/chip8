package main

import (
	"time"
)

// Adjust these according to your game/preferences
const SCALING = 15
const FPS = 60
const speed = 10

// Timing vars
var startTime time.Time
var now time.Time
var then time.Time
var elapsed time.Duration
var fpsInterval time.Duration

func main() {
	// Decide how often we render a frame
	fpsInterval = time.Duration(time.Millisecond * 1000.0 / FPS)
	then = time.Now()
	startTime = then

	// Load the chip8 font and ROM into memory
	loadSprites()
	loadROM("ROMs/SPACEINVADERS.ch8")

	// Start display and take one step forward (render one frame)
	startDisplay()
	step()
}

func step() {
	elapsed = time.Since(then)

	// Wait till it's time to render a frame
	if elapsed > fpsInterval {
		cycle()
	}
}
