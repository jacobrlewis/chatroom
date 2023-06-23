package main

import (
	"errors"
	"log"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("Got / request")
	io.WriteString(w, "Welcome message!\n")
}

func main() {
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Printf("unexpected error: %s", err)
		os.Exit(1)
	}
}