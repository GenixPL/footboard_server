package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var MyUpgrader = websocket.Upgrader{}

func setupRoutes() {
	http.HandleFunc("/", OnHome)

	// TODO: regular request for list of games
	// TODO: regular request for creating game
	// TODO: socket request for taking part in game
	http.HandleFunc("/create", OnCreate)

}

func main() {
	fmt.Println("Go Web!")

	setupRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
