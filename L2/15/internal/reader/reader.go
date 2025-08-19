package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func ReadInput() string {
	// Канал для сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Канал для результата
	result := make(chan string, 1)

	go func() {
		buf := bufio.NewReader(os.Stdin)
		input, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nПолучен EOF (Ctrl+D)")
				result <- ""
				os.Exit(0)
			}
			fmt.Println("Ошибка чтения:", err)
			result <- ""
			return
		}
		result <- input
	}()

	select {
	case sig := <-sigs:
		fmt.Println("\nПойман сигнал:", sig)
		return ""
	case input := <-result:
		return input
	}
}
