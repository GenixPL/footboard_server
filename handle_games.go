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
	var game *Game
	for i := 0; i < len(Games); i++ {
		if Games[i].Id == gameId {
			game = &Games[i]
			break
		}
	}

	if game == nil {
		msg := "{\"error\": \"no_such_game\", \"newGame\": null}"
		connection.WriteMessage(1, []byte(msg))
		return
	}

	game.AddClient(connection)

	gameJsonString, err := game.ToJsonString()
	if err != nil {
		msg := "{\"error\": \"invalid_game_json\", \"newGame\": null}"
		connection.WriteMessage(1, []byte(msg))
		return
	}

	msg := "{\"error\": null, \"newGame\": " + gameJsonString + "}"
	connection.WriteMessage(1, []byte(msg))
}
