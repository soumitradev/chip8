package main

import (
	"fmt"
	"math/rand"
	"os"
)

var memory [4096]uint8

var registers [16]uint8   // v
var mem_register uint = 0 // i

var delayTimer uint = 0
var soundTimer uint = 0

var ip uint = 0x200

var stack = []uint8{}

var paused = false

var speed = 10

func loadSprites() {
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
	for i := 0; i < len(sprites); i++ {
		memory[i] = sprites[i]
	}
}

func loadProgram(program []uint8) {
	for i := 0x200; i < len(program); i++ {
		memory[i] = program[i]
	}
}

func loadROM(romPath string) {
	romFile, err := os.ReadFile(romPath)
	if err != nil {
		panic(err)
	}
	loadProgram(romFile)
}

func cycle() {
	for i := 0; i < speed; i++ {
		if !paused {
			fmt.Printf("%d\n%v\n%v", ip, stack, registers)
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

func executeInstruction(opcode uint) {
	ip += 2

	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := (opcode & 0xF000)

	switch n {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			clearScreen()
			break
		case 0x00EE:
			ip, stack = pop(len(stack)-1, stack)
			break
		}
		break
	case 0x1000:
		ip = uint(opcode & 0xFFF)
		break
	case 0x2000:
		stack = append(stack, uint8(ip))
		ip = opcode & 0xFFF
		break
	case 0x3000:
		if registers[x] == uint8(opcode&0xFF) {
			ip += 2
		}
		break
	case 0x4000:
		if registers[x] != uint8(opcode&0xFF) {
			ip += 2
		}
		break
	case 0x5000:
		if registers[x] == registers[y] {
			ip += 2
		}
		break
	case 0x6000:
		registers[x] = uint8(opcode & 0xFF)
		break
	case 0x7000:
		registers[x] += uint8(opcode & 0xFF)
		break
	case 0x8000:
		switch opcode & 0xF {
		case 0x0:
			registers[x] = registers[y]
			break
		case 0x1:
			registers[x] |= registers[y]
			break
		case 0x2:
			registers[x] &= registers[y]
			break
		case 0x3:
			registers[x] ^= registers[y]
			break
		case 0x4:
			sum := registers[x] + registers[y]
			registers[0xF] = 0
			if sum > 0xFF {
				registers[0xF] = 1
			}
			registers[x] = sum
			break
		case 0x5:
			registers[0xF] = 0
			if registers[x] > registers[y] {
				registers[0xF] = 1
			}
			registers[x] -= registers[y]
			break
		case 0x6:
			registers[0xF] = registers[x] & 0x1
			registers[x] >>= 1
			break
		case 0x7:
			registers[0xF] = 0
			if registers[y] > registers[x] {
				registers[0xF] = 1
			}
			registers[x] = registers[y] - registers[x]
			break
		case 0xE:
			registers[0xF] = registers[x] & 0x80
			registers[x] <<= 1
			break
		}
		break
	case 0x9000:
		if registers[x] != registers[y] {
			ip += 2
		}
		break
	case 0xA000:
		mem_register = (opcode & 0xFFF)
		break
	case 0xB000:
		ip = opcode&0xFFF + uint(registers[0])
		break
	case 0xC000:
		randN := uint(rand.Intn(256))
		registers[x] = uint8(randN & (opcode & 0xFF))
		break
	case 0xD000:
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
			if isKeyPressed(registers[x]) {
				ip += 2
			}
			break
		case 0xA1:
			if !isKeyPressed(registers[x]) {
				ip += 2
			}
			break
		}
		break
	case 0xF000:
		switch opcode & 0xFF {
		case 0x07:
			registers[x] = uint8(delayTimer)
			break
		case 0x0A:
			paused = true
			onNextPress = func(key uint8) {
				registers[x] = key
				paused = false
			}
			break
		case 0x15:
			delayTimer = uint(registers[x])
			break
		case 0x18:
			soundTimer = uint(registers[x])
			break
		case 0x1E:
			mem_register += uint(registers[x])
			break
		case 0x29:
			mem_register = uint(registers[x] * 5)
			break
		case 0x33:
			memory[mem_register] = uint8(registers[x] / 100)
			memory[mem_register+1] = uint8((registers[x] % 100) / 10)
			memory[mem_register+2] = uint8(registers[x] % 10)
			break
		case 0x55:
			for regIndex := 0; regIndex <= int(x); regIndex++ {
				memory[mem_register+uint(regIndex)] = registers[regIndex]
			}
			break
		case 0x65:
			for regIndex := 0; regIndex <= int(x); regIndex++ {
				registers[regIndex] = memory[mem_register+uint(regIndex)]
			}
			break
		}
		break
	default:
		panic("Invalid opcode")
	}
}

func updateTimers() {
	if delayTimer > 0 {
		delayTimer -= 1
	}
	if soundTimer > 0 {
		soundTimer -= 1
	}
}
