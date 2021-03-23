package game

import (
	m "footboard_server/models"
	u "footboard_server/models/utils"

	"github.com/gorilla/websocket"
)

const (
	errorInvalidGameJson = "invalid_game_json"
)

func (game *Game) AddClient(connection *websocket.Conn) {
	newClient := m.NewClient(connection)
	game.Clients = append(game.Clients, newClient)

	gameJsonString, err := game.ToJsonString()
	if err != nil {
		u.LogE("AddClient", "getting game json failed")

		err = connection.WriteMessage(1, u.JsonedErrWithUid(errorInvalidGameJson, newClient.Id))
		if err != nil {
			u.LogE("AddClient", "sending msg failed")
		}
	}

	err = connection.WriteMessage(1, u.JsonedMsgWithUid(gameJsonString, newClient.Id))
	if err != nil {
		u.LogE("AddClient", "sending msg failed")
		return
	}

	game.SendUpdateToEveryClient()

	go game.handleMessages(&newClient)
}
