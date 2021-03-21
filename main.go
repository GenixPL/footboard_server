package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/websocket"
)

var MyUpgrader = websocket.Upgrader{}

// ====== ROUTES ======

func OnHome(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the handler
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		return
	}

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

func OnConnect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/connect triggered")

	ws, err := MyUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	AddClient(ws, path.Base(r.URL.Path))
}

// ======

func setupRoutes() {
	http.HandleFunc("/", OnHome)
	http.HandleFunc("/create-new-game", OnCreateNewGame)
	http.HandleFunc("/games", OnGetGames)
	http.HandleFunc("/connect/", OnConnect)
}

func main() {
	fmt.Println("Go Web!")

	setupRoutes()
	// StartPeriodicEmptyGamesRemoval()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
