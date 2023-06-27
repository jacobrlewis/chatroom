package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jacobrlewis/chatroom/pkg/shared"
	"log"
	"net/http"
)

type Client struct {
	Host string
	Username string
	Room string
}

// connect to a server for the first time
func (client Client) initConnection() int {

	url := client.Host + shared.GetInitEndpoint()

	clientHello := shared.ClientHello{
		Username: client.Username,
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

	fmt.Println(serverHello.WelcomeMsg)
	return serverHello.RoomCount
}

// join a room
func (client Client) joinRoom() {

	welcome_url := client.Host + shared.GetRoomWelcomeEndpoint(client.Room)

	clientHello := shared.ClientHello{
		Username: client.Username,
	}
	helloBytes, _ := json.Marshal(clientHello)

	resp, err := http.Post(welcome_url, "json", bytes.NewReader(helloBytes))
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

	fmt.Println(serverHello.WelcomeMsg)
}

// send a chat to a room
func sendChat(room_url string, msgStruct shared.Msg) {
	
	jsonStr, err := json.Marshal(msgStruct)
	if err != nil {
		log.Println("Message failed to encode")
		return
	}

	jsonBody := []byte(jsonStr)
	bodyReader := bytes.NewReader(jsonBody)
	resp, err := http.Post(room_url, "json", bodyReader)
	if err != nil {
		log.Println("Message failed to send")
		return
	}
	resp.Body.Close()
}

// endlessly send messages
func (client Client) chatLoop() {
	
	room_url := client.Host + shared.GetRoomSendChatEndpoint(client.Room)
	for {
		msg, err := ReadMsg()
		if err != nil {
			continue
		}
		msgStruct := shared.Msg{
			Username: client.Username,
			Msg:      msg,
		}
		// send without waiting
		go sendChat(room_url, msgStruct)
	}
}

func main() {
	host := GetServerUrl()
	username := GetUsername()

	client := Client{Username: username, Host: host, Room: ""}

	num_rooms := client.initConnection()
	fmt.Printf("There are %d rooms on this server.\n", num_rooms)

	client.Room = GetRoomId()

	client.joinRoom()

	
	client.chatLoop()
}
