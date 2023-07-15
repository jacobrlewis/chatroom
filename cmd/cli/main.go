package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"jacobrlewis/chatroom/pkg/shared"
	"log"
	"net/http"
)

type Client struct {
	Host     string
	Username string
	Room     string
	Conn     *websocket.Conn
}

// connect to a server for the first time
func (client Client) initConnection() int {

	url := shared.GetInitUrl(client.Host)
	fmt.Println(url)

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

	// hellos
	welcome_url := shared.GetRoomWelcomeUrl(client.Host, client.Room)
	fmt.Println(welcome_url)

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

	// open websocket
	ws_url := shared.GetRoomWsUrl(client.Host, client.Room)
	headers := http.Header{}
	headers.Set("X-Client-Info", string(helloBytes))
	conn, _, err := websocket.DefaultDialer.Dial(ws_url.String(), headers)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server: ", err)
	}
	client.Conn = conn
	go client.receiveMessages()
	client.chatLoop()
}

func (client Client) receiveMessages() {
	for {
		_, bytes, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from server:", err)
			return
		}
		var msg shared.Msg
		err = json.Unmarshal(bytes, &msg)
		if err != nil {
			log.Println("Failed to unmarshal JSON body", http.StatusBadRequest)
			return
		}
		fmt.Println(msg.Username + ": " + msg.Msg)
	}
}

// send a chat to a room
func sendChat(conn *websocket.Conn, msgStruct shared.Msg) {

	jsonStr, err := json.Marshal(msgStruct)
	if err != nil {
		log.Println("Message failed to encode")
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonStr)
	if err != nil {
		log.Println("Error sending message to server:", err)
		return
	}
}

// endlessly send messages
func (client Client) chatLoop() {

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
		go sendChat(client.Conn, msgStruct)
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
}
