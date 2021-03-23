package game

import (
	m "footboard_server/models"
	u "footboard_server/models/utils"
)

func (game *Game) handleMove(client *m.Client, jsonReq map[string]interface{}) {
	x := int(jsonReq["val"].(map[string]interface{})["x"].(float64))
	y := int(jsonReq["val"].(map[string]interface{})["y"].(float64))

	// Check if the game is in proper state.
	if !(game.State == gameStateRunning) {
		client.Connection.WriteMessage(1, u.JsonedErr(errorGameGameIsNotRunning))
		return
	}

	isPlayer1 := client.Id == game.Player1.Id
	isPlayer2 := client.Id == game.Player2.Id

	// Check if the client is at least one of the players.
	if !isPlayer1 && !isPlayer2 {
		client.Connection.WriteMessage(1, u.JsonedErr(errorNonPlayerCannotMove))
		return
	}

	// Check if it's the client's turn.
	if (game.MovesPlayer1 && isPlayer2) || (!game.MovesPlayer1 && isPlayer1) {
		client.Connection.WriteMessage(1, u.JsonedErr(errorNotYourTurn))
		return
	}

	newMove := m.Move{
		SX:          game.Ball.X,
		SY:          game.Ball.Y,
		EX:          x,
		EY:          y,
		PerformedBy: client.Id,
	}

	game.Moves = append(game.Moves, newMove)

	u.LogV("handleMove", "new move: ("+string(x)+", "+string(y)+")")

	game.SendUpdateToEveryClient()
}
