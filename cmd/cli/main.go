package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jacobrlewis/chatroom/pkg/shared"
	"log"
	"net/http"
	"strings"
)

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

func send_chats(room_url string, username string) {

	for {
		msg, err := ReadMsg()

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
	url := GetServerUrl()

	username := GetUsername()

	num_rooms := init_connection(url + "/welcome", username)
	fmt.Printf("There are %d rooms on this server.\n", num_rooms)

	room_id := GetRoomId()
	room_url := url + "/room/" + room_id
	fmt.Println(room_url)

	send_chats(room_url, username)
}
