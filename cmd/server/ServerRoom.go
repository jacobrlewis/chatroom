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
	// connections
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to ws. ", err)
		return
	}

	helloBytes := r.Header.Get("X-Client-Info")

	var hello shared.ClientHello
	err = json.Unmarshal([]byte(helloBytes), &hello)
	if err != nil {
		log.Println("Could not unmarshall ClientHello. ", err)
		http.Error(w, "Invalid ClientHello", http.StatusBadRequest)
		return
	}

	log.Printf("Room %s started web socket connection with %s", room.Id, hello.Username)

	// TODO handle duplicate username
	room.Conns[hello.Username] = conn
	go room.wsListen(conn)

}

func (room ServerRoom) wsListen(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from client:", err)
			break
		}

		var msg shared.Msg
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Failed to unmarshal JSON body", http.StatusBadRequest)
			return
		}
		log.Println(fmt.Sprintf("Room %s received message: %+v", room.Id, msg))
	}
}
