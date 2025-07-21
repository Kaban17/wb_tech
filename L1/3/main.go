package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func worker(id int, jobs <-chan string) {
	for job := range jobs {
		fmt.Printf("Job %d recieved: %s\n", id, job)
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		return
	}
	num_workers, err := strconv.Atoi(os.Args[1])
	if err != nil || num_workers <= 0 {
		fmt.Println("Invalid number of workers")
		return
	}
	jobs := make(chan string)

	for i := range num_workers {
		go worker(i, jobs)
	}
	cnt := 0
	for range num_workers * num_workers {
		msg := fmt.Sprintf("Job %d", cnt)
		jobs <- msg
		cnt++
		time.Sleep(time.Second)
	}
}
