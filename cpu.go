package main

import (
	"math/rand"
	"os"
)

// Allocate memory for registers, memory, timers, stack, other vars
var memory [4096]uint8

var registers [16]uint8
var mem_register uint = 0

var delayTimer uint = 0
var soundTimer uint = 0

var ip uint16 = 0x200

var stack = []uint16{}

var paused = false

// Load sprites into lower memory
func loadSprites() {
	// Font sprites
	sprites := [...]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	// Copy to low memory
	for i := 0; i < len(sprites); i++ {
		memory[i] = sprites[i]
	}
}

// Load ROM into memory (starting at 0x200)
func loadProgram(program []uint8) {
	// Error on overflow
	if len(program) > (0xE00) {
		panic("Program too large!")
	}
	for i := 0; i < len(program); i++ {
		memory[0x200+i] = program[i]
	}
}

// Read and load ROM
func loadROM(romPath string) {
	romFile, err := os.ReadFile(romPath)
	if err != nil {
		panic(err)
	}
	loadProgram(romFile)
}

// Every display cycle, process a bunch of CPU instructions from ROM if not paused,
// update timers, play sound and render graphics.
func cycle() {
	for i := 0; i < speed; i++ {
		if !paused {
			opcode := ((uint(memory[ip]) << 8) | uint(memory[ip+1]))
			executeInstruction(opcode)
		}
	}

	if !paused {
		updateTimers()
	}

	playSound()
	screenRender()
}

// Execute instruction provided according to chip8 arch
func executeInstruction(opcode uint) {
	ip += 2

	// Gte parameters in opcode
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := (opcode & 0xF000)

	switch n {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			// Clear screen
			clearScreen()
			break
		case 0x00EE:
			// Return from coroutine
			ip, stack = pop(len(stack)-1, stack)
			break
		}
		break
	case 0x1000:
		// Jump to lowest 12 bits of instruction
		ip = uint16(opcode & 0xFFF)
		break
	case 0x2000:
		// Call coroutine to go to lowest 12 bits of instruction
		stack = append(stack, ip)
		ip = uint16(opcode & 0xFFF)
		break
	case 0x3000:
		// Skip next instruction if value in register x is equal to lowest 12 bits of instruction
		if registers[x] == uint8(opcode&0xFF) {
			ip += 2
		}
		break
	case 0x4000:
		// Skip next instruction if value in register x is not equal to lowest 12 bits of instruction
		if registers[x] != uint8(opcode&0xFF) {
			ip += 2
		}
		break
	case 0x5000:
		// Skip next instruction if value in register x is equal to value on register y
		if registers[x] == registers[y] {
			ip += 2
		}
		break
	case 0x6000:
		// Load value of lowest byte into register x
		registers[x] = uint8(opcode & 0xFF)
		break
	case 0x7000:
		// Add value of lowest byte to register x
		registers[x] += uint8(opcode & 0xFF)
		break
	case 0x8000:
		switch opcode & 0xF {
		case 0x0:
			// Load value of register y into register x
			registers[x] = registers[y]
			break
		case 0x1:
			// Store value of OR of register x and y into register x
			registers[x] |= registers[y]
			break
		case 0x2:
			// Store value of AND of register x and y into register x
			registers[x] &= registers[y]
			break
		case 0x3:
			// Store value of XOR of register x and y into register x
			registers[x] ^= registers[y]
			break
		case 0x4:
			// Add register x and y, and store it in register x.
			// If an overflow is caused, set register F (flag register) to 1
			sum := uint(registers[x]) + uint(registers[y])
			registers[0xF] = 0
			if sum > 0xFF {
				registers[0xF] = 1
			}
			registers[x] = uint8(sum)
			break
		case 0x5:
			// Subtract value in register y from value in register x, and store it in register x.
			// If result is positive (i.e. value in register x > that in y), set register F to 1
			registers[0xF] = 0
			if registers[x] > registers[y] {
				registers[0xF] = 1
			}
			registers[x] -= registers[y]
			break
		case 0x6:
			// Store least significant bit of register x in register F,
			// and then divide value in register x
			registers[0xF] = registers[x] & 0x1
			registers[x] >>= 1
			break
		case 0x7:
			// Subtract value in register x from value in register y, and store it in register x.
			// If result is positive (i.e. value in register y > that in x), set register F to 1
			registers[0xF] = 0
			if registers[y] > registers[x] {
				registers[0xF] = 1
			}
			registers[x] = registers[y] - registers[x]
			break
		case 0xE:
			// WARNING: DOCS FAULTY ON THIS ONE.
			// REFER http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xyE FOR CORRECT VERSION
			// Store most significant bit in register F, and then multiply value in register x
			registers[0xF] = ((registers[x] >> 7) & 1)
			registers[x] <<= 1
			break
		}
		break
	case 0x9000:
		// Skip next instruction if values in registers x and y are not equal
		if registers[x] != registers[y] {
			ip += 2
		}
		break
	case 0xA000:
		// Store lowest 12 bits of instruction in memory register
		mem_register = (opcode & 0xFFF)
		break
	case 0xB000:
		// Jump to sum of lowest 12 bits of instruction and value in register 0
		ip = uint16(opcode&0xFFF + uint(registers[0]))
		break
	case 0xC000:
		// Set register x to AND of random number between
		// 0 and 255 (inclusive) and lowest byte of instruction
		randN := uint(rand.Intn(256))
		registers[x] = uint8(randN & (opcode & 0xFF))
		break
	case 0xD000:
		// Display n-byte sprite at co-ordinates (value of register x, value of register y)
		// n = lowest 4 bits of instruction, if a collision is detected, set register F to 1
		width := 8
		height := int(opcode & 0xF)
		registers[0xF] = 0
		for row := 0; row < height; row++ {
			sprite := memory[int(mem_register)+row]

			for col := 0; col < width; col++ {
				if (sprite & 0x80) > 0 {
					if setPixel(int(registers[x])+col, int(registers[y])+row) {
						registers[0xF] = 1
					}
				}

				sprite <<= 1

			}
		}
		break
	case 0xE000:
		switch opcode & 0xFF {
		case 0x9E:
			// Skip next instruction if key with keycode in register x is pressed
			if isKeyPressed(registers[x]) {
				ip += 2
			}
			break
		case 0xA1:
			// Skip next instruction if key with keycode in register x is not pressed
			if !isKeyPressed(registers[x]) {
				ip += 2
			}
			break
		}
		break
	case 0xF000:
		switch opcode & 0xFF {
		case 0x07:
			// Set register x to delayTimer
			registers[x] = uint8(delayTimer)
			break
		case 0x0A:
			// Wait for a key press and store the keycode in register x
			paused = true
			onNextPress = func(key uint8) {
				registers[x] = key
				paused = false
			}
			break
		case 0x15:
			// Set delayTimer to value in register x
			delayTimer = uint(registers[x])
			break
		case 0x18:
			// Set soundTimer to value in register x
			soundTimer = uint(registers[x])
			break
		case 0x1E:
			// Increment memory register by value in register x
			mem_register += uint(registers[x])
			break
		case 0x29:
			// Set memory register to location of sprite for digit for value in register x
			mem_register = uint(registers[x] * 5)
			break
		case 0x33:
			// Store decimal representation of value in register x in memory at
			// memory locations memory location, memory location + 1, memory location + 2
			// as hundreds, tens, and units digits
			memory[mem_register] = uint8(registers[x] / 100)
			memory[mem_register+1] = uint8((registers[x] % 100) / 10)
			memory[mem_register+2] = uint8(registers[x] % 10)
			break
		case 0x55:
			// Copy values in registers 0 to x in memory locations starting at location of memory register
			for regIndex := 0; regIndex <= int(x); regIndex++ {
				memory[mem_register+uint(regIndex)] = registers[regIndex]
			}
			break
		case 0x65:
			// Copy memory locations starting at location of memory register to values in registers 0 to x
			for regIndex := 0; regIndex <= int(x); regIndex++ {
				registers[regIndex] = memory[mem_register+uint(regIndex)]
			}
			break
		}
		break
	default:
		// Invalid opcode error
		panic("Invalid opcode")
	}
}

// Update timers every cycle
func updateTimers() {
	if delayTimer > 0 {
		delayTimer -= 1
	}
	if soundTimer > 0 {
		soundTimer -= 1
	}
}
