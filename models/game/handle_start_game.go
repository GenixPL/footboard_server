package game

import (
	m "footboard_server/models"
	u "footboard_server/utils"
)

func (game *Game) handleStartGame(client *m.Client, jsonReq map[string]interface{}) {
	// Simply return if the client has already started.
	if client.StartedGame == true {
		return
	}

	// Check if game can actually be started.
	if (game.State != gameStateHasTwoPlayers) && (game.State != gameStateOnePlayerStarted) {
		client.Connection.WriteMessage(1, u.JsonedErr(errorGameCannotBeStartedYet))
		return
	}

	client.StartedGame = true

	// Update GameState.
	if game.Player1.StartedGame && !game.Player2.StartedGame {
		// Only 1st player has started.
		game.State = gameStateOnePlayerStarted
	} else if game.Player1.StartedGame && !game.Player2.StartedGame {
		// Only 2nd player has started.
		game.State = gameStateOnePlayerStarted
	} else if game.Player1.StartedGame && game.Player2.StartedGame {
		// Both players have started.
		game.State = gameStateRunning
	}

	game.SendUpdateToEveryClient()
}
