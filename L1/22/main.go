package main

import (
	"fmt"
	"math/big"
)

func main() {
	// Создаем большие числа (примеры значений > 2^20)
	a := new(big.Int)
	b := new(big.Int)

	// Устанавливаем значения из строк (можно использовать числа любого размера)
	a.SetString("1000000000000000000000", 10) // 10^21
	b.SetString("500000000000000000000", 10)  // 5*10^20

	// Выполняем операции
	fmt.Println("Сложение:", add(a, b))
	fmt.Println("Вычитание:", sub(a, b))
	fmt.Println("Умножение:", mul(a, b))
	fmt.Println("Деление:", div(a, b))
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

// Деление больших чисел (с обработкой деления на ноль)
func div(a, b *big.Int) *big.Float {
	if b.Sign() == 0 {
		panic("Деление на ноль!")
	}

	// Конвертируем big.Int в big.Float для точного деления
	fa := new(big.Float).SetInt(a)
	fb := new(big.Float).SetInt(b)

	return new(big.Float).Quo(fa, fb)
}
