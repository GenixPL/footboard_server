package game

import (
	"encoding/json"
	"errors"
	m "footboard_server/models"
	u "footboard_server/utils"

	"github.com/google/uuid"
)

// Game states
const (
	gameStateWaitingForPlayers = "waiting_for_players"
	gameStateHasOnePlayer      = "has_one_player"
	gameStateHasTwoPlayers     = "has_two_players"
	gameStateOnePlayerStarted  = "one_player_started"
	// We skip two_players_started because that's where the gameStateRunning state starts.
	gameStateRunning         = "running"
	gameStatePaused          = "paused"
	gameStateFirstPlayerWon  = "first_player_won"
	gameStateSecondPlayerWon = "second_player_won"
)

// Errors
const (
	errorReceivedInvalidJson        = "received_invalid_json"
	errorPlaceAlreadyOccupied       = "place_already_occupied"
	errorGameIsNotWaitingForPlayers = "game_is_not_waiting_for_players"
	errorCouldntParseGameToJson     = "couldnt_parse_game_to_json"
	errorGameCannotBeStartedYet     = "game_cannot_be_started_yet"
	errorGameGameIsNotRunning       = "game_is_not_running"
	errorNonPlayerCannotMove        = "non_player_cannot_move"
	errorNotYourTurn                = "not_your_turn"
	errorInvalidMove                = "invalid_move"
)

// Commands
const (
	commandOccupyPlace = "occupy_place"
	commandStartGame   = "start_game"
	commandMove        = "move"
)

// Game model.
//
// Shouldn't be created through constructor, but NewGame() instead.
type Game struct {
	Id string `json:"id"`

	// All Clients getting updates about this game (includes player1 and player2).
	Clients []m.Client `json:"clients"`

	// Clients that take part in this game.
	Player1 *m.Client `json:"player1"`
	Player2 *m.Client `json:"player2"`

	// If true then it means that it's first player's turn.
	// If false then second player's.
	MovesPlayer1 bool `json:"movesPlayer1"`

	// Describes current position of the ball.
	Ball m.Point `json:"ball"`

	// List of already visited VisitedPoints.
	VisitedPoints []m.Point `json:"visitedPoints"`

	PossiblePoints []m.Point `json:"possiblePoints"`

	// List of Moves that were performed.
	Moves []m.Move `json:"moves"`

	// String secribing current state of this game.
	State string `json:"gameState"`
}

// Creates new Game object with proper initial values.
func NewGame() Game {
	return Game{
		Id:           uuid.NewString(),
		Clients:      []m.Client{},
		MovesPlayer1: true,
		Ball: m.Point{
			X: 0,
			Y: 0,
		},
		VisitedPoints: []m.Point{
			{
				X: 0,
				Y: 0,
			},
		},
		PossiblePoints: []m.Point{
			{X: -1, Y: -1},
			{X: 0, Y: -1},
			{X: 1, Y: -1},
			{X: -1, Y: 0},
			{X: 1, Y: 0},
			{X: -1, Y: 1},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
		},
		Moves: []m.Move{},
		State: gameStateWaitingForPlayers,
	}
}

// Returns string consisting of the Game object encoded as JSON.
func (game Game) ToJsonString() (string, error) {
	byteArray, err := json.Marshal(game)
	if err != nil {
		u.LogE("ToJsonString", "Error in Marshal, e: "+err.Error())
		return "", errors.New("Couldn't parse Game")
	}

	return string(byteArray), nil
}

// Reads incoming messages and removes client if it disconnects.
func (game *Game) handleMessages(client *m.Client) {
	for {
		_, p, err := client.Connection.ReadMessage()
		if err != nil {
			u.LogE("handleMessages", "read msg resulted in err: "+err.Error())
			game.RemoveClientWithId(client.Id)
			return
		}

		game.handleMessage(client, p)
	}
}

func (game *Game) handleMessage(client *m.Client, msg []byte) {
	u.LogV("handleMessage", "received new message: "+string(msg))

	var jsonReq map[string]interface{}
	err := json.Unmarshal(msg, &jsonReq)
	if err != nil {
		u.LogV("handleMessage", "coudln't parse the message, err: "+err.Error())
		err = client.Connection.WriteMessage(1, []byte(u.JsonedErr(errorReceivedInvalidJson)))
		if err != nil {
			u.LogE("handleMessage", "send msg failed, err: "+err.Error())
		}
		return
	}

	command := jsonReq["command"]

	// Check if command has proper type.
	switch command.(type) {
	case string:
		break
	default:
		err := client.Connection.WriteMessage(1, u.JsonedErr(errorReceivedInvalidJson))
		if err != nil {
			u.LogE("handleMessage", "send msg failed, err: "+err.Error())
		}
		return
	}

	if command == commandOccupyPlace {
		game.handleOccupyPlace(client, jsonReq)
		return
	}

	if command == commandStartGame {
		game.handleStartGame(client, jsonReq)
		return
	}

	if command == commandMove {
		game.handleMove(client, jsonReq)
		return
	}
}
