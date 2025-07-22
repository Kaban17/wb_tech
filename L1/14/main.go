package main

import (
	"fmt"
)

func detectType(v any) {
	switch val := v.(type) {
	case int:
		fmt.Printf("Тип: int, значение: %d\n", val)
	case string:
		fmt.Printf("Тип: string, значение: %q\n", val)
	case bool:
		fmt.Printf("Тип: bool, значение: %t\n", val)
	case chan int:
		fmt.Println("Тип: chan int")
	case chan string:
		fmt.Println("Тип: chan string")
	case chan bool:
		fmt.Println("Тип: chan bool")
	default:
		fmt.Printf("Неизвестный тип: %T\n", val)
	}
}

func main() {
	detectType(42)
	detectType("hello")
	detectType(true)
	detectType(make(chan int))
	detectType(make(chan string))
	detectType(make(chan bool))
	detectType(3.14)
}
