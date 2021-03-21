package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var Games []Game = []Game{}

// Periodically removes games without clients.
func StartPeriodicEmptyGamesRemoval() {
	duration := time.Duration(5) * time.Minute

	ticker := time.NewTicker(duration)

	for range ticker.C {
		removeEmptyGames()
	}
}

// Removes games that don't have any clients.
func removeEmptyGames() {
	fmt.Println("Removing empty games...")
}

// Creates new game.
func CreateNewGame() Game {
	fmt.Println("Creating new game...")

	game := NewGame()

	Games = append(Games, game)

	return game
}

func GamesToJsonStr() string {
	str := "["

	for i := 0; i < len(Games); i++ {
		if i != 0 {
			str += ", "
		}

		gameJsonString, err := Games[i].ToJsonString()
		if err != nil {
			continue
		}

		str += gameJsonString
	}

	str += "]"

	return str
}

func AddClient(connection *websocket.Conn, gameId string) {
	for {
		messageType, p, err := connection.ReadMessage()
		if err != nil {
			fmt.Println("CLOSING")
			fmt.Println(err)
			return
		}

		fmt.Println(string(p), messageType)
	}
}
