package telnet

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	types "wb_tech/l2_17/pkg"
)

func Connect(ctx context.Context, config *types.Config) error {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	fmt.Printf("Attempting to connect to %s...\n", addr)

	conn, err := net.DialTimeout("tcp", addr, config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	// Ensure the connection is closed when the function exits
	defer conn.Close()
	fmt.Printf("Connected successfully. Press Ctrl+C or Ctrl+D to exit.\n")

	var wg sync.WaitGroup
	defer wg.Wait() // Wait for all goroutines to finish before returning

	// --- 1. Goroutine for reading from the network and writing to stdout ---
	// This handles data coming from the server.
	wg.Add(1)
	go func() {
		defer wg.Done()
		// defer conn.Close() здесь избыточно, так как conn.Close() уже в defer главной функции Connect.
		// Однако, оставляем, чтобы гарантировать закрытие соединения в случае, если сервер
		// закрывает его первым (EOF), что является важным условием для завершения io.Copy.
		defer conn.Close()

		// Copy data from the network connection to the standard output (terminal).
		// io.Copy blocks until EOF (connection closed) or an error occurs.
		if _, err := io.Copy(os.Stdout, conn); err != nil && err != io.EOF {
			// Check if the context was cancelled; if so, this is often an expected shutdown.
			select {
			case <-ctx.Done():
				fmt.Printf("\r\nConnection reader closed due to cancellation.\n")
			default:
				// Only report unexpected errors
				fmt.Fprintf(os.Stderr, "\r\nError reading from connection: %v\n", err)
			}
		}
	}()

	// --- 2. Goroutine for reading from stdin and writing to the network ---
	// This handles user input going to the server.
	wg.Add(1)
	go func() {
		defer wg.Done()
		// defer conn.Close() гарантирует, что соединение будет закрыто,
		// если пользователь нажмет Ctrl+D (EOF для os.Stdin).
		defer conn.Close()

		// Copy data from the standard input (terminal) to the network connection.
		if _, err := io.Copy(conn, os.Stdin); err != nil && err != io.EOF {
			select {
			case <-ctx.Done():
				// User manually quit (Ctrl+C), or an external signal was received.
				// io.Copy might return an error if stdin is closed externally.
				fmt.Printf("\r\nInput writer closed due to cancellation.\n")
			default:
				// Only report unexpected errors
				fmt.Fprintf(os.Stderr, "\r\nError writing to connection: %v\n", err)
			}
		}
	}()

	// --- 3. Block and wait for context cancellation (Ctrl+C) or for goroutines to trigger conn.Close() ---
	// Context cancellation (Ctrl+C) is the explicit way to quit.
	// If one of the goroutines closes the connection, the other will fail and exit,
	// and the main routine will eventually exit after wg.Wait().
	<-ctx.Done()

	// Since we closed the connection in the goroutines upon copy termination,
	// if we reach here via context cancellation, we ensure both streams are done.
	return nil
}
