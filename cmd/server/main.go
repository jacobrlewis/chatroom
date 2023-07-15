package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"jacobrlewis/chatroom/pkg/shared"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
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

	log.Printf("%s connected to the server\n", hello.Username)

	msg := shared.ServerHello{
		RoomCount:  NUM_ROOMS,
		WelcomeMsg: fmt.Sprintf("Howdy %s, welcome to the server!", hello.Username),
	}
	msgBytes, _ := json.Marshal(msg)
	w.Write(msgBytes)
}

// Start server
func main() {
	http.HandleFunc(shared.GetInitEndpoint(), welcomeHandler)

	for i := 1; i <= NUM_ROOMS; i++ {
		str_id := fmt.Sprint(i)
		room := ServerRoom{
			Id:             str_id,
			RoomWelcomeMsg: fmt.Sprintf("Welcome to room %s!", str_id),
			Password:       "",
			Conns:          make(map[string]*websocket.Conn),
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
