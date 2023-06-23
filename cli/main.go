package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
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
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(string(body))

}