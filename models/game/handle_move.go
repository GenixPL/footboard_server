package game

import (
	m "footboard_server/models"
	u "footboard_server/utils"
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

	// Check if haven't moved outside of boundaries.
	if !moveIsValid(x, y) {
		client.Connection.WriteMessage(1, u.JsonedErr(errorInvalidMove))
		return
	}

	newMove := m.Move{
		SP: m.Point{
			X: game.Ball.X,
			Y: game.Ball.Y,
		},
		EP: m.Point{
			X: x,
			Y: y,
		},
		PerformedBy: client.Id,
	}

	// Check if such move has already been done
	if moveExists(game.Moves, newMove) {
		client.Connection.WriteMessage(1, u.JsonedErr(errorInvalidMove))
		return
	}

	game.Moves = append(game.Moves, newMove)
	game.Ball = m.Point{
		X: x,
		Y: y,
	}

	// Check if the point has already been visited.
	if !pointWasVisited(game.VisitedPoints, x, y) {
		if isPlayer1 {
			game.MovesPlayer1 = false
		} else {
			game.MovesPlayer1 = true
		}
	}

	u.LogV("handleMove", "new move: ("+string(rune(x))+", "+string(rune(y))+")")

	// Check possible moves (coveers corner cases).
	if len(possibleMoves(game.Moves, x, y)) == 0 {
		if isPlayer1 {
			game.State = gameStateSecondPlayerWon
		} else {
			game.State = gameStateFirstPlayerWon
		}
	}

	// Check if Player1 has won.
	if moveIsInLowerGoal(x, y) {
		game.State = gameStateFirstPlayerWon
	}

	// Check if Player2 has won.
	if moveIsInUpperGoal(x, y) {
		game.State = gameStateSecondPlayerWon
	}

	game.SendUpdateToEveryClient()
}

func moveIsInUpperGoal(x int, y int) bool {
	return (y == 6) && (-1 <= x && x <= 1)
}

func moveIsInLowerGoal(x int, y int) bool {
	return (y == -6) && (-1 <= x && x <= 1)
}

func pointIsCorner(x int, y int) bool {
	return (x == -4 && y == -5) || (x == 4 && y == -5) || (x == -4 && y == 5) || (x == 4 && y == 5)
}

// TODO doesnt cover walking on walls
func moveIsValid(x int, y int) bool {
	// One of the goals.
	if moveIsInUpperGoal(x, y) || moveIsInUpperGoal(x, y) {
		return true
	}

	// Main field area.
	if (-4 <= x && x <= 4) && (-5 <= y && y <= 5) {
		return true
	}

	// Otherwise.
	return false
}

func moveExists(moves []m.Move, move m.Move) bool {
	for _, m := range moves {
		// If start and end point are the same.
		if move.EP.IsEqualToPoint(m.EP) && move.SP.IsEqualToPoint(m.SP) {
			return true
		}

		// If start and end point are switched.
		if move.EP.IsEqualToPoint(m.SP) && move.SP.IsEqualToPoint(m.EP) {
			return true
		}
	}

	return false
}

func pointWasVisited(vistedPoint []m.Point, x int, y int) bool {
	for _, p := range vistedPoint {
		if p.X == x && p.Y == y {
			return true
		}
	}

	return false
}

// TODO: test this
func possibleMoves(moves []m.Move, x int, y int) []m.Point {
	possiblePoints := []m.Point{}

	// Top Left
	if x >= -3 && y <= 4 {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y + 1})
	}

	// Top
	if (x != -4 && x != 4) && y <= 4 {
		possiblePoints = append(possiblePoints, m.Point{X: x, Y: y + 1})
	}

	// Top Right
	if x <= 3 && y <= 4 {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y + 1})
	}

	// Left
	if x >= -3 && (y != 5 && y != -5) {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y})
	}

	// Right
	if x <= 3 && (y != 5 && y != -5) {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y})
	}

	// Bottom Left
	if x >= -3 && y >= -4 {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y - 1})
	}

	// Bottom
	if (x != -4 && x != 4) && y >= 4 {
		possiblePoints = append(possiblePoints, m.Point{X: x, Y: y - 1})
	}

	// Bottom Right
	if x <= 3 && y >= -4 {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y - 1})
	}

	// The above cover main field point, so we have to add the near goal ones.
	// =================

	// Top left goal line.
	if x == -1 && y == 5 {
		// Top Right
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y + 1})
		// Right
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y})
	}

	// Top center goal line.
	if x == 0 && y == 5 {
		// Top Right
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y + 1})
		// Top Left
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y + 1})
		// Right
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y})
		// Left
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y})
	}

	// Top right goal line.
	if x == 1 && y == 5 {
		// Top Left
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y + 1})
		// Left
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y})
	}

	// Bottom left goal line.
	if x == -1 && y == -5 {
		// Bottom Right
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y - 1})
		// Right
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y})
	}

	// Bottom center goal line.
	if x == 0 && y == -5 {
		// Bottom Right
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y - 1})
		// Bottom Left
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y - 1})
		// Right
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y})
		// Left
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y})
	}

	// Bottom right goal line.
	if x == 1 && y == -5 {
		// Bottom Left
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y - 1})
		// Left
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y})
	}

	// Remove already performed moves.
	// =================

	wantedPoints := []m.Point{}
	for _, p := range possiblePoints {
		newMove := m.Move{
			SP: m.Point{X: x, Y: y},
			EP: p,
		}

		if !moveExists(moves, newMove) {
			wantedPoints = append(wantedPoints, p)
		}
	}

	return wantedPoints
}
