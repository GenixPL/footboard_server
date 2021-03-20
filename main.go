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
	http.HandleFunc("/create-new-game", OnCreateNewGame)
	http.HandleFunc("/games", OnGetGames)
	// TODO: regular request for list of games
	// TODO: regular request for creating game
	// TODO: socket request for taking part in game

}

func main() {
	fmt.Println("Go Web!")

	setupRoutes()
	// StartPeriodicEmptyGamesRemoval()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
