package main

import "fmt"

func setBit(n int64, i uint) int64 {
	return n | (1 << i)
}
func clearBit(n int64, i uint) int64 {
	return n &^ (1 << i)
}
func main() {
	var num int64 = 10
	fmt.Printf("Decimal: %d. Binary: %b\n", num, num)
	var bitPos uint = 1
	fmt.Printf("Decimal: %d. Binary: %b\n", clearBit(num, bitPos), clearBit(num, bitPos))
	fmt.Printf("Decimal: %d. Binary: %b\n", setBit(num, bitPos), setBit(num, bitPos))
}
