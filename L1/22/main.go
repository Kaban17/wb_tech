package main

import (
	"fmt"
	"math/big"
)

func main() {
	a := new(big.Int)
	b := new(big.Int)
	a.SetString("1000000000000000000000", 10) // 10^21
	b.SetString("500000000000000000000", 10)  // 5*10^20
	sum := add(a, b)
	diff := sub(a, b)
	prod := mul(a, b)
	quot := div(a, b)

	fmt.Println("Сложение:", sum)
	fmt.Println("Вычитание:", diff)
	fmt.Println("Умножение:", prod)
	fmt.Println("Деление:", quot)
}

// Сложение больших чисел
func add(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

// Вычитание больших чисел
func sub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

// Умножение больших чисел
func mul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

// Деление больших чисел
func div(a, b *big.Int) *big.Rat {
	// Используем Rat для точного представления дробных результатов
	return new(big.Rat).SetFrac(a, b)
}
