package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"strings"
)

func GetServerUrl() string {
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
	url := "http://" + host + ":" + port
	return url
}

func GetUsername() string {
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
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter room id: ")
	id, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read server id")
	}
	id = strings.TrimSpace(id)

	return id
}

func ReadMsg() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("You: ")
	return reader.ReadString('\n')
}