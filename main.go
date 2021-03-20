package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var MyUpgrader = websocket.Upgrader{}

// ====== ROUTES ======

func OnHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/ triggered")

	fmt.Fprintf(w, "Welcome Adventurer!")
}

func OnCreateNewGame(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/create-new-game triggered")

	game := CreateNewGame()

	gameJsonString, err := game.ToJsonString()
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	}

	fmt.Fprintf(w, gameJsonString)
}

func OnGetGames(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/games triggered")

	fmt.Fprintf(w, GamesToJsonStr())
}

// ======

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
