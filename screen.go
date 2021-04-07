package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Allocate memory for graphics
var display = [64 * 32]bool{}

// Create a IMD drawing object
var imd = imdraw.New(nil)

// Function that acts as a keyboard event listener
var onNextPress func(uint8)

// Create window
var window pixelgl.Window

// Declare variable maps for controls
var CONTROLS = []pixelgl.Button{
	pixelgl.Key1,
	pixelgl.Key2,
	pixelgl.Key3,
	pixelgl.Key4,
	pixelgl.KeyQ,
	pixelgl.KeyW,
	pixelgl.KeyE,
	pixelgl.KeyR,
	pixelgl.KeyA,
	pixelgl.KeyS,
	pixelgl.KeyD,
	pixelgl.KeyF,
	pixelgl.KeyZ,
	pixelgl.KeyX,
	pixelgl.KeyC,
	pixelgl.KeyV,
}

var keysDown = make(map[uint8]bool)

var KEYMAP = map[uint8]pixelgl.Button{
	0x1: pixelgl.Key1,
	0x2: pixelgl.Key2,
	0x3: pixelgl.Key3,
	0xC: pixelgl.Key4,
	0x4: pixelgl.KeyQ,
	0x5: pixelgl.KeyW,
	0x6: pixelgl.KeyE,
	0xD: pixelgl.KeyR,
	0x7: pixelgl.KeyA,
	0x8: pixelgl.KeyS,
	0x9: pixelgl.KeyD,
	0xE: pixelgl.KeyF,
	0xA: pixelgl.KeyZ,
	0x0: pixelgl.KeyX,
	0xB: pixelgl.KeyC,
	0xF: pixelgl.KeyV,
}
var KEYMAP_REV = map[string]uint8{
	pixelgl.Key1.String(): 0x1,
	pixelgl.Key2.String(): 0x2,
	pixelgl.Key3.String(): 0x3,
	pixelgl.Key4.String(): 0xC,
	pixelgl.KeyQ.String(): 0x4,
	pixelgl.KeyW.String(): 0x5,
	pixelgl.KeyE.String(): 0x6,
	pixelgl.KeyR.String(): 0xD,
	pixelgl.KeyA.String(): 0x7,
	pixelgl.KeyS.String(): 0x8,
	pixelgl.KeyD.String(): 0x9,
	pixelgl.KeyF.String(): 0xE,
	pixelgl.KeyZ.String(): 0xA,
	pixelgl.KeyX.String(): 0x0,
	pixelgl.KeyC.String(): 0xB,
	pixelgl.KeyV.String(): 0xF,
}

// Run the graphics emulation via pixel
func run() {
	// Configure window
	cfg := pixelgl.WindowConfig{
		Title:  "Chip8",
		Bounds: pixel.R(0, 0, 64*SCALING, 32*SCALING),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Set ON pixel color to white
	imd.Color = pixel.RGB(1, 1, 1)

	// Update every frame
	for !window.Closed() {
		// Handle keyboard events
		for _, key := range CONTROLS {
			if window.Pressed(key) {
				keysDown[KEYMAP_REV[key.String()]] = true

				if onNextPress != nil && key != 0 {
					onNextPress(uint8(KEYMAP_REV[key.String()]))
					onNextPress = nil
				}
			} else {
				keysDown[KEYMAP_REV[key.String()]] = false
			}
		}
		// Take next step forward
		step()

		//	Clear window and draw the pixels
		window.Clear(colornames.Black)
		imd.Draw(window)
		window.Update()
	}
}

// Check if key is currently down
func isKeyPressed(keycode uint8) bool {
	return keysDown[keycode]
}

// Start display
func startDisplay() {
	pixelgl.Run(run)
}

// Render pixels on display according to graphics memory
func screenRender() {
	imd.Clear()
	for i := 0; i < 64*32; i++ {
		x := float64((i % 64) * SCALING)
		y := 31*SCALING - float64((i/64)*SCALING)
		if display[i] {
			imd.Push(pixel.V(x, y))
			imd.Push(pixel.V(x+SCALING, y+SCALING))
			imd.Rectangle(0)
		}
	}
}

// Dummy test render frunction
func testRender() {
	setPixel(0, 0)
	setPixel(63, 31)
}

// Set pixel on at x, y coords from top left corner
func setPixel(x int, y int) bool {
	if x >= 64 {
		x -= 64
	} else if x < 0 {
		x += 64
	}

	if y >= 32 {
		y -= 32
	} else if y < 0 {
		y += 32
	}
	pixelLoc := x + (y * 64)
	display[pixelLoc] = (true != display[pixelLoc])
	return !display[pixelLoc]
}

// Clear screen
func clearScreen() {
	display = [64 * 32]bool{}
}
