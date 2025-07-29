package main

import (
	"fmt"
	"os"
	"wb_tech/l2_8/mytime"
)

func main() {
	t, err := mytime.Now()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Println("Current time:", t)
}
