package main

import (
	"time"

	"github.com/alexflint/go-arg"
)

// Adjust these according to your game preferences
var SCALING = 15
var FPS = 60
var speed = 10

type args struct {
	ROMPATH  string `arg:"positional" help:"Path to CHIP8 ROM file."`
	SCALING  int    `default:"15" arg:"-s" help:"Pixel scaling. Adjusts size of display."`
	FPS      int    `default:"60" arg:"-f" help:"FPS to run display at."`
	CPUspeed int    `default:"10" arg:"-c" help:"Speed of CPU relative to FPS."`
}

// Return version info
func (args) Version() string {
	return "chip8 v0.1.0"
}

// Timing vars
var startTime time.Time
var now time.Time
var then time.Time
var elapsed time.Duration
var fpsInterval time.Duration

func main() {
	// Use prefs in args
	var args args

	arg.MustParse(&args)
	var ROMPATH = args.ROMPATH
	SCALING = args.SCALING
	FPS = args.FPS
	speed = args.CPUspeed

	// Decide how often we render a frame
	fpsInterval = time.Duration(time.Millisecond * 1000.0 / time.Duration(FPS))
	then = time.Now()
	startTime = then

	// Load the chip8 font and ROM into memory
	loadSprites()
	loadROM(ROMPATH)

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
