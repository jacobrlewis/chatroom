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
		RoomCount: 1,
		WelcomeMsg: fmt.Sprintf("Howdy %s", hello.Username),
	}
	msgBytes, _ := json.Marshal(msg)
	w.Write(msgBytes)
}

func receiveMessage(w http.ResponseWriter, r *http.Request) {
	
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

	log.Println("Received message: ")
	log.Println(fmt.Sprintf("%+v\n", msg))
	io.WriteString(w, "Thanks for the message\n")
}

func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/room/1", receiveMessage)

	log.Println("Server started")
	err := http.ListenAndServe("localhost:3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Printf("unexpected error: %s", err)
		os.Exit(1)
	}
}