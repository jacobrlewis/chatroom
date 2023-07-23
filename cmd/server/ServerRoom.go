package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jacobrlewis/chatroom/pkg/shared"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ServerRoom struct {
	Id             string
	RoomWelcomeMsg string
	Password       string
	Conns map[string]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// start a chat room
func (room ServerRoom) beginRoom() {
	http.HandleFunc(shared.GetRoomWelcomeEndpoint(room.Id), room.roomWelcomehandler)
	http.HandleFunc(shared.GetRoomWsEndpoint(room.Id), room.startWs)
}

// listen for incoming connections to a room
func (room ServerRoom) roomWelcomehandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Could not read ClientHello. ", err)
		http.Error(w, "Invalid ClientHello", http.StatusBadRequest)
		return
	}
	var hello shared.ClientHello
	err = json.Unmarshal(body, &hello)
	if err != nil {
		log.Println("Could not unmarshall ClientHello. ", err)
		http.Error(w, "Invalid ClientHello", http.StatusBadRequest)
		return
	}

	log.Printf("%s connected to room %s\n", hello.Username, room.Id)

	msg := shared.ServerHello{
		RoomCount:  0,
		WelcomeMsg: room.RoomWelcomeMsg,
	}
	msgBytes, _ := json.Marshal(msg)
	w.Write(msgBytes)
}

// start websocket connection with a client
func (room ServerRoom) startWs(w http.ResponseWriter, r *http.Request) {
	log.Printf("Room %s starting web socket connection", room.Id)
	
	helloBytes := r.Header.Get("X-Client-Info")
	var hello shared.ClientHello
	err := json.Unmarshal([]byte(helloBytes), &hello)
	if err != nil {
		log.Println("Could not unmarshall ClientHello. ", err)
		http.Error(w, "Invalid ClientHello", http.StatusBadRequest)
		return
	}

	if room.Conns[hello.Username] != nil {
		http.Error(w, "Duplicate username", http.StatusConflict)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to ws. ", err)
		return
	}
	log.Printf("Room %s started web socket connection with %s", room.Id, hello.Username)

	room.Conns[hello.Username] = conn
	room.sendJoinAlert(hello.Username)
	room.wsListen(conn, hello.Username)
}

// receive messages for one websocket connection
func (room ServerRoom) wsListen(conn *websocket.Conn, username string) {
	defer conn.Close()
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client:", err)
			log.Printf("Closing connection with %s", username)

			// remove conn from known connections
			delete(room.Conns, username)
			room.sendleaveAlert(username)

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("Client closed unexpectedly.")
			} else {
				log.Println("Client closed connection.")
			}
			break
		}

		var msg shared.Msg
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			log.Println("Failed to unmarshal JSON body", http.StatusBadRequest)
			continue
		}
		log.Println(fmt.Sprintf("Room %s received message: %+v", room.Id, msg))

		// validate values before sending to other clients
		msg.Join = false
		msg.Join = false
		room.sendMsg(username, msgBytes)
	}
}

// send message to all clients that a user joined
func (room ServerRoom) sendJoinAlert(username string) {

	msg := shared.Msg{
		Username: username,
		Join: true,
	}
	msgBytes, _ := json.Marshal(msg)
	room.sendMsg(username, msgBytes)
}

// send message to all clients that a user joined
func (room ServerRoom) sendleaveAlert(username string) {

	msg := shared.Msg{
		Username: username,
		Leave: true,
	}
	msgBytes, _ := json.Marshal(msg)
	room.sendMsg(username, msgBytes)
}

// send msg to all clients except sender
func (room ServerRoom) sendMsg(sender string, msg []byte) {
	for u, out := range room.Conns {
		if u != sender {
			go out.WriteMessage(websocket.TextMessage, msg)
		}
	}
}