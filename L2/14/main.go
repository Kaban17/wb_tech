package main

import (
	"fmt"
	"time"
)

// or реализует функцию, которая объединяет несколько каналов.
// Возвращаемый канал закрывается, как только закрывается любой
// из исходных каналов.
func or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		c := make(chan any)
		close(c)
		return c
	case 1:
		return channels[0]
	default:
		m := len(channels) / 2
		return mergeTwo(or(channels[:m]...), or(channels[m:]...))
	}
}

func mergeTwo(a, b <-chan any) <-chan any {
	c := make(chan any)
	go func() {
		defer close(c)
		select {
		case <-a:
		case <-b:
		}
	}()
	return c
}

func main() {
	sig := func(after time.Duration) <-chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("Готово после %v\n", time.Since(start))
}
