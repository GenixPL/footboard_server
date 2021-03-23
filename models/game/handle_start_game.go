package game

import (
	m "footboard_server/models"
	u "footboard_server/models/utils"
)

func (game *Game) handleStartGame(client *m.Client, jsonReq map[string]interface{}) {
	// Simply return if the client has already started.
	if client.StartedGame == true {
		return
	}

	// Check if game can actually be started.
	if game.state != gameStateHasTwoPlayers && game.state != gameStateOnePlayerStarted {
		client.Connection.WriteMessage(1, u.JsonedErr(errorGameCannotBeStartedYet))
		return
	}

	client.StartedGame = true

	// Update GameState.
	if game.Player1.StartedGame && !game.Player2.StartedGame {
		// Only 1st player has started.
		game.state = gameStateOnePlayerStarted
	} else if game.Player1.StartedGame && !game.Player2.StartedGame {
		// Only 2nd player has started.
		game.state = gameStateOnePlayerStarted
	} else if game.Player1.StartedGame && game.Player2.StartedGame {
		// Both players have started.
		game.state = gameStateRunning
	}

	game.SendUpdateToEveryClient()
}
