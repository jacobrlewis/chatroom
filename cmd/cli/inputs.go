package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var HOST_ENV = "GO_CHAT_HOST"
var USERNAME_ENV = "GO_CHAT_USERNAME"
var ROOM_ENV = "GO_CHAT_ROOM"

func GetServerUrl() string {

	if os.Getenv(HOST_ENV) != "" {
		fmt.Printf("Using %s : %s as host\n", HOST_ENV, os.Getenv(HOST_ENV))
		return os.Getenv(HOST_ENV)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Host (localhost): ")
	host, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read input")
	}
	host = strings.TrimSpace(host)

	if host == "" {
		fmt.Println("Defaulting to localhost")
		host = "localhost"
	}

	port := "3333"
	host = host + ":" + port
	return host
}

func GetUsername() string {

	if os.Getenv(USERNAME_ENV) != "" {
		fmt.Printf("Using %s : %s as username\n", USERNAME_ENV, os.Getenv(USERNAME_ENV))
		return os.Getenv(USERNAME_ENV)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read username")
	}
	username = strings.TrimSpace(username)

	return username
}

func GetRoomId() string {

	if os.Getenv(ROOM_ENV) != "" {
		fmt.Printf("Using %s : %s as room\n", ROOM_ENV, os.Getenv(ROOM_ENV))
		return os.Getenv(ROOM_ENV)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter room id: ")
	id, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read server id")
	}
	id = strings.TrimSpace(id)

	return id
}

func ReadMsg(reader *bufio.Reader) (string, error) {
	fmt.Print("You: ")
	msg, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read message")
		return msg, err
	}
	msg = strings.TrimSpace(msg)
	return msg, err
}
