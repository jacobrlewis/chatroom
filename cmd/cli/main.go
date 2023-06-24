package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jacobrlewis/chatroom/pkg/shared"
	"log"
	"net/http"
	"os"
	"strings"
)

func get_server_url() string {
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

func init_connection(url string, username string) int {

	clientHello := shared.ClientHello{
		Username: username,
	}
	helloBytes, _ := json.Marshal(clientHello)

	resp, err := http.Post(url, "json", bytes.NewReader(helloBytes))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response from server")
		log.Fatal(err)
	}

	var serverHello shared.ServerHello
	err = json.Unmarshal(body, &serverHello)
	if err != nil {
		log.Fatal("Failed to unmarshall server welcome")
	}

	return serverHello.RoomCount
}

func read_msg() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("You: ")
	return reader.ReadString('\n')
}

func send_chats(room_url string, username string) {

	for {
		msg, err := read_msg()

		if err != nil {
			log.Println("Failed to read message")
			continue
		}
		msg = strings.TrimSpace(msg)

		msgStruct := shared.Msg{
			Username: username,
			Msg:      msg,
		}
		jsonStr, err := json.Marshal(msgStruct)
		if err != nil {
			log.Println("Message failed to encode")
			continue
		}

		jsonBody := []byte(jsonStr)
		bodyReader := bytes.NewReader(jsonBody)
		resp, err := http.Post(room_url, "text/plain", bodyReader)
		if err != nil {
			log.Println("Message failed to send")
			continue
		}
		resp.Body.Close()
	}
}

func main() {
	url := get_server_url()

	// TODO prompt for username
	username := "username"

	num_rooms := init_connection(url, username)
	fmt.Printf("There are %d rooms on this server.\n", num_rooms)

	// TODO: prompt for room_id to connect to
	room_id := "1"
	room_url := url + "/room/" + room_id

	send_chats(room_url, username)
}
