package main

import (
	"time"
)

const SCALING = 20
const FPS = 60

var startTime time.Time
var now time.Time
var then time.Time
var elapsed time.Duration
var fpsInterval time.Duration

func main() {
	fpsInterval = time.Duration(time.Millisecond * 1000.0 / FPS)
	then = time.Now()
	startTime = then

	loadSprites()
	loadROM("ROMs/TEST.ch8")
	hexdump(memory[:])
	step()
	startDisplay()
}

func step() {
	elapsed = time.Since(then)

	if elapsed > fpsInterval {
		cycle()
	}
}
