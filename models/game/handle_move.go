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

	// Check if the client is one of the players.
	if !isPlayer1 && !isPlayer2 {
		client.Connection.WriteMessage(1, u.JsonedErr(errorNonPlayerCannotMove))
		return
	}

	// Check if it's the client's turn.
	if (game.MovesPlayer1 && isPlayer2) || (!game.MovesPlayer1 && isPlayer1) {
		client.Connection.WriteMessage(1, u.JsonedErr(errorNotYourTurn))
		return
	}

	// Check if the move is one of the valid moves.
	if !moveIsPossible(game.PossiblePoints, x, y) {
		client.Connection.WriteMessage(1, u.JsonedErr(errorInvalidMove))
		return
	}

	// FROM HERE WE KNOW THAT THE MOVE IS VALID

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

	game.Moves = append(game.Moves, newMove)
	game.Ball = m.Point{X: x, Y: y}
	game.PossiblePoints = getPossiblePoints(game.Moves, x, y)

	// Check if the point was already visited (if it was then it's still the player's turn).
	if !pointWasVisited(game.VisitedPoints, x, y) {
		if isPlayer1 {
			game.MovesPlayer1 = false
		} else {
			game.MovesPlayer1 = true
		}
	}

	if moveIsInLowerGoal(x, y) {
		game.State = gameStateFirstPlayerWon

	} else if moveIsInUpperGoal(x, y) {
		game.State = gameStateSecondPlayerWon

	} else if len(game.PossiblePoints) == 0 {
		if isPlayer1 {
			game.State = gameStateSecondPlayerWon
		} else {
			game.State = gameStateFirstPlayerWon
		}
	}

	game.SendUpdateToEveryClient()
}

func moveIsPossible(possiblePoints []m.Point, x int, y int) bool {
	for _, p := range possiblePoints {
		if p.X == x && p.Y == y {
			return true
		}
	}

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

func moveIsInUpperGoal(x int, y int) bool {
	return (y == 6) && (-1 <= x && x <= 1)
}

func moveIsInLowerGoal(x int, y int) bool {
	return (y == -6) && (-1 <= x && x <= 1)
}

func getPossiblePoints(moves []m.Move, x int, y int) []m.Point {
	possiblePoints := []m.Point{}

	// edge rows
	isLeftMostRow := (x == -4)
	isRightMostRow := (x == 4)
	isTopMostRow := (y == 5)
	isBottomMostRow := (y == -5)

	// top goal line
	isLeftTopGoalLine := (x == -1 && y == 5)
	isCenterTopGoalLine := (x == 0 && y == 5)
	isRightTopGoalLine := (x == 1 && y == 5)

	// bottom goal line
	isLeftBottomGoalLine := (x == -1 && y == -5)
	isCenterBottomGoalLine := (x == 0 && y == -5)
	isRightBottomGoalLine := (x == 1 && y == -5)

	// Top Left
	if !isLeftMostRow && !isTopMostRow {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y + 1})
	} else if isCenterTopGoalLine || isRightTopGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y + 1})
	}

	// Top
	if !isLeftMostRow && !isTopMostRow && !isRightMostRow {
		possiblePoints = append(possiblePoints, m.Point{X: x, Y: y + 1})
	} else if isCenterTopGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x, Y: y + 1})
	}

	// Top Right
	if !isRightMostRow && !isTopMostRow {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y + 1})
	} else if isLeftTopGoalLine || isCenterTopGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y + 1})
	}

	// Left
	if !isLeftMostRow && !isTopMostRow && !isBottomMostRow {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y})
	} else if isCenterTopGoalLine || isRightTopGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y})
	} else if isCenterBottomGoalLine || isRightBottomGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y})
	}

	// Right
	if !isTopMostRow && !isRightMostRow && !isBottomMostRow {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y})
	} else if isLeftTopGoalLine || isCenterTopGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y})
	} else if isLeftBottomGoalLine || isCenterBottomGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y})
	}

	// Bottom Left
	if !isLeftMostRow && !isBottomMostRow {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y - 1})
	} else if isRightBottomGoalLine || isCenterBottomGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x - 1, Y: y - 1})
	}

	// Bottom
	if !isLeftMostRow && !isRightMostRow && !isBottomMostRow {
		possiblePoints = append(possiblePoints, m.Point{X: x, Y: y - 1})
	} else if isCenterBottomGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x, Y: y - 1})
	}

	// Bottom Right
	if !isRightMostRow && !isBottomMostRow {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y - 1})
	} else if isLeftBottomGoalLine || isCenterBottomGoalLine {
		possiblePoints = append(possiblePoints, m.Point{X: x + 1, Y: y - 1})
	}

	// Remove already performed moves.
	filteredPoints := []m.Point{}
	for _, p := range possiblePoints {
		newMove := m.Move{
			SP: m.Point{X: x, Y: y},
			EP: p,
		}

		if !moveExists(moves, newMove) {
			filteredPoints = append(filteredPoints, p)
		}
	}

	return filteredPoints
}
