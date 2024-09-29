package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "5252"
)

type Client struct {
	Username string
}

func startServer() error {
	clients := map[net.Conn]Client{}
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		// slog.Error("could not start server")
		// return error
		return errors.New("starting as an client")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("Error while accepting the connection")
		}
		clients[conn] = Client{
			Username: "test",
		}
		log.Printf("Client connected: %s\n", conn.RemoteAddr())

		go handleConnection(conn, clients)
	}
}

func handleConnection(conn net.Conn, clients map[net.Conn]Client) {
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			slog.Error("could not read from the connection")
		}
		log.Println(string(buffer[:n]))

		for client := range clients {
			if client != conn {
				_, err := client.Write(buffer[:n])
				if err != nil {
					slog.Error("could not write to client")
				}
			}
		}
	}
}

func startClient() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		slog.Error("could not connect to the server")
		return err
	}

	defer conn.Close()

	go func() {
		for {
			message := make([]byte, 1024)
			n, err := conn.Read(message)
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			fmt.Println(string(message[:n]))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		msg := scanner.Text()
		_, err := conn.Write([]byte(msg))
		if err != nil {
			slog.Error("could not write to the server")
			return err
		}
	}
}

func main() {
	err := startServer()
	if err != nil {
		log.Println(err)
		err := startClient()
		if err != nil {
			log.Println(err)
		}
	}

	log.Printf("Server started at %s:%s", HOST, PORT)
}
