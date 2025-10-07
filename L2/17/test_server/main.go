package main

// test -server
import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const listenPort = "8080"

func handleConnection(conn net.Conn) {
	defer conn.Close()

	welcomeMsg := fmt.Sprintf("Hello, world from test server! Connected to port %s.\n", listenPort)
	if _, err := conn.Write([]byte(welcomeMsg)); err != nil {
		log.Printf("Ошибка при отправке приветствия: %v", err)
		return
	}

	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	log.Printf("Начата обработка соединения от %s...", conn.RemoteAddr())

	_, err := io.Copy(conn, conn)

	if err != nil && err != io.EOF {
		log.Printf("Соединение с %s закрыто с ошибкой: %v", conn.RemoteAddr(), err)
	} else {
		log.Printf("Соединение с %s закрыто клиентом (EOF)", conn.RemoteAddr())
	}
}

func main() {
	listener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер на порту %s: %v", listenPort, err)
	}
	defer listener.Close()
	log.Printf("TCP Echo-сервер запущен на порту :%s. Ожидание соединений...", listenPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка при приеме соединения: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}
