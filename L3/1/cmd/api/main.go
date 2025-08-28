package main

import (
	"fmt"
	"wb_tech/l3_1/pkg/types"
)

func main() {
	fmt.Println("Hello, world from L3/1!")
	model := types.NewNotification("1", "Hello, world!")
	fmt.Println(model)

}
