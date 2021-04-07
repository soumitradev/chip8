package main

import (
	"fmt"
)

// Pops an item off a list and returns both
func pop(i int, xs []uint8) (uint, []uint8) {
	y := xs[i]
	ys := append(xs[:i], xs[i+1:]...)
	return uint(y), ys
}

// hexdump prints a hex dump of the data given
func hexdump(dump []uint8) {
	// Print columns for hex dump
	fmt.Print("\n\n     ")
	for i := 0; i < 16; i++ {
		fmt.Printf(" %02X ", i)
	}

	// Check number of rows and extra rows (rows at end of file that countain less than 16 bytes)
	N := len(dump) / 16
	ext := len(dump) % 16

	// Print normal rows that contain 16 bytes
	for j := 0; j < N; j++ {
		fmt.Printf("\n%02X0  ", j)
		for k := 0; k < 16; k++ {
			fmt.Printf(" %02X ", dump[16*j+k])
		}
	}

	// Print last extra row
	if ext > 0 {
		fmt.Printf("\n%02X0  ", N)
		for k := 0; k < ext; k++ {
			fmt.Printf(" %02X ", dump[16*N+k])
		}
	}
	// Newline at the end
	fmt.Println()
}
