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

	"github.com/gorilla/websocket"
)

type Client struct {
	Host     string
	Username string
	Room     string
	Conn     *websocket.Conn
	Reader   *bufio.Reader
}

// connect to a server for the first time
func (client Client) initConnection() int {

	url := shared.GetInitUrl(client.Host)

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
	welcomeUrl := shared.GetRoomWelcomeUrl(client.Host, client.Room)

	clientHello := shared.ClientHello{
		Username: client.Username,
	}
	helloBytes, _ := json.Marshal(clientHello)

	resp, err := http.Post(welcomeUrl, "json", bytes.NewReader(helloBytes))
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
	wsUrl := shared.GetRoomWsUrl(client.Host, client.Room)
	headers := http.Header{}
	headers.Set("X-Client-Info", string(helloBytes))
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), headers)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server: ", err)
	}
	client.Conn = conn
	go client.receiveMessages()
	client.chatLoop()
}

// receive messages from the server
func (client Client) receiveMessages() {
	for {
		_, bytes, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from server:", err)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Server closed unexpectedly.")
			} else {
				log.Println("Server closed connection.")
			}
			// TODO close entire client
			return
		}
		var msg shared.Msg
		err = json.Unmarshal(bytes, &msg)
		if err != nil {
			log.Println("Failed to unmarshal JSON body", http.StatusBadRequest)
			return
		}

		fmt.Printf("Got message %v", msg)

		// TODO read current input buffer (if user has anything typed not yet sent)

		// clear current line and print new message
		clearCurrentLine := "\033[2K"
		fmt.Print(clearCurrentLine + "\r")
		fmt.Println(msg.Username + ": " + msg.Msg)

		// re-print prompt
		fmt.Print("You: ")
		// // TODO print anything user was typing
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
		msg, err := ReadMsg(client.Reader)
		if err != nil || msg == "" {
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

	client := Client{Username: username, Host: host, Room: "", Reader: bufio.NewReader(os.Stdin)}

	numRooms := client.initConnection()
	fmt.Printf("There are %d rooms on this server.\n", numRooms)

	client.Room = GetRoomId()

	client.joinRoom()
}
