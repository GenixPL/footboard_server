package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var MyUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func setupRoutes() {
	http.HandleFunc("/", OnHome)
	http.HandleFunc("/create", OnCreate)
	http.HandleFunc("/connect", OnConnect)
}

func main() {
	fmt.Println("Go Web!")

	setupRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
