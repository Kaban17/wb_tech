package main

import "fmt"

// Родительская структура
type Human struct {
	Name string
	Age  int
}

// Метод Human
func (h Human) Greet() {
	fmt.Printf("Привет, меня зовут %s и мне %d лет.\n", h.Name, h.Age)
}

// Структура-наследник (композиция)
type Action struct {
	Human // встраивание структуры Human
	Power int
}

// Дополнительный метод Action
func (a Action) Act() {
	fmt.Printf("%s действует с силой %d!\n", a.Name, a.Power)
}

func main() {
	a := Action{
		Human: Human{
			Name: "Кабан",
			Age:  20,
		},
		Power: 100,
	}

	a.Greet() // Метод унаследован от Human
	a.Act()   // Метод свой собственный
}
