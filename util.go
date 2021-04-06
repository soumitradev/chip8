package main

func pop(i int, xs []uint8) (uint, []uint8) {
	y := xs[i]
	ys := append(xs[:i], xs[i+1:]...)
	return uint(y), ys
}
