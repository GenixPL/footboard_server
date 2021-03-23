package main

import (
	"fmt"
	gm "footboard_server/models/game"

	"github.com/gorilla/websocket"
)

var Games []gm.Game = []gm.Game{}

// Creates new Game and adds it to the Games array.
func CreateNewGame() gm.Game {
	fmt.Println("Creating new game...")

	game := gm.NewGame()
	Games = append(Games, game)

	return game
}

// Returns string consting of Games as JSON array.
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

func AddClientToGame(connection *websocket.Conn, gameId string) {
	var game *gm.Game
	for i := 0; i < len(Games); i++ {
		if Games[i].Id == gameId {
			game = &Games[i]
			break
		}
	}

	// If client wants to connect to a game that doesn't exist.
	if game == nil {
		msg := "{\"error\": \"no_such_game\", \"game\": null}"
		connection.WriteMessage(1, []byte(msg))
		return
	}

	game.AddClient(connection)
}
