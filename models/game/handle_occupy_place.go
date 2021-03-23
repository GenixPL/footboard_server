package game

import (
	m "footboard_server/models"
	u "footboard_server/models/utils"
)

func (game *Game) handleOccupyPlace(client *m.Client, jsonReq map[string]interface{}) {
	place := int(jsonReq["val"].(float64))

	// Check game state.
	if !(game.State == gameStateWaitingForPlayers || game.State == gameStateHasOnePlayer) {
		client.Connection.WriteMessage(1, u.JsonedErr(errorGameIsNotWaitingForPlayers))
		return
	}

	if place == 1 {
		if game.Player1 != nil {
			client.Connection.WriteMessage(1, u.JsonedErr(errorPlaceAlreadyOccupied))
			return
		}

		game.Player1 = client

		if game.Player2 != nil {
			game.State = gameStateHasTwoPlayers
		} else {
			game.State = gameStateHasOnePlayer
		}

		game.SendUpdateToEveryClient()
		return

	} else if place == 2 {
		if game.Player2 != nil {
			client.Connection.WriteMessage(1, u.JsonedErr(errorPlaceAlreadyOccupied))
			return
		}

		game.Player2 = client

		if game.Player1 != nil {
			game.State = gameStateHasTwoPlayers
		} else {
			game.State = gameStateHasOnePlayer
		}

		game.SendUpdateToEveryClient()
		return

	} else {
		client.Connection.WriteMessage(1, u.JsonedErr(errorReceivedInvalidJson))
		return
	}
}
