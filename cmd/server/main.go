package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jacobrlewis/chatroom/pkg/shared"
	"log"
	"net/http"
	"os"
	"io/ioutil"
)

var NUM_ROOMS = 2

// handle new connections to the server
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Could not read ClientHello")
		http.Error(w, "Invalid ClientHello", http.StatusBadRequest)
		return
	}
	var hello shared.ClientHello
	err = json.Unmarshal(body, &hello)
	if err != nil {
		log.Println("Could not unmarshall hello")
		http.Error(w, "Invalid ClientHello", http.StatusBadRequest)
		return
	}

	log.Printf("%s connected\n", hello.Username)

	msg := shared.ServerHello{
		RoomCount: NUM_ROOMS,
		WelcomeMsg: fmt.Sprintf("Howdy %s, welcome to the server!", hello.Username),
	}
	msgBytes, _ := json.Marshal(msg)
	w.Write(msgBytes)
}

type ServerRoom struct {
	Id string
	RoomWelcomeMsg string
	Password string
	// connections
}

// start a chat room
func (room ServerRoom) beginRoom() {
	http.HandleFunc(shared.GetRoomWelcomeEndpoint(room.Id), room.roomWelcomehandler)
	http.HandleFunc(shared.GetRoomSendChatEndpoint(room.Id), room.receiveMessage)
}

// listen for incoming connections to a room
func (room ServerRoom) roomWelcomehandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Could not read ClientHello")
		http.Error(w, "Invalid ClientHello", http.StatusBadRequest)
		return
	}
	var hello shared.ClientHello
	err = json.Unmarshal(body, &hello)
	if err != nil {
		log.Println("Could not unmarshall hello")
		http.Error(w, "Invalid ClientHello", http.StatusBadRequest)
		return
	}

	log.Printf("%s connected to room %s\n", hello.Username, room.Id)

	msg := shared.ServerHello{
		RoomCount: 0,
		WelcomeMsg: room.RoomWelcomeMsg,
	}
	msgBytes, _ := json.Marshal(msg)
	w.Write(msgBytes)
}

// listen for incoming messages to a room
func (room ServerRoom) receiveMessage(w http.ResponseWriter, r *http.Request) {
	
    body, err := ioutil.ReadAll(r.Body)
	if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
	
	var msg shared.Msg
	err = json.Unmarshal(body, &msg)
	if err != nil {
        http.Error(w, "Failed to unmarshal JSON body", http.StatusBadRequest)
        return
    }

	log.Println(fmt.Sprintf("Room %s received message: %+v", room.Id, msg))
	io.WriteString(w, "Thanks for the message\n")
}

// Start server
func main() {
	http.HandleFunc(shared.GetInitEndpoint(), welcomeHandler)

	for i := 1; i <= NUM_ROOMS; i++ {
		str_id := fmt.Sprint(i)
        room := ServerRoom{
			Id: str_id, 
			RoomWelcomeMsg: fmt.Sprintf("Welcome to room %s!", str_id), 
			Password: "",
		}
		room.beginRoom()
    }

	log.Println("Server started")
	err := http.ListenAndServe("localhost:3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Printf("unexpected error: %s", err)
		os.Exit(1)
	}
}