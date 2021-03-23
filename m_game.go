package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
)

// Commands
const (
	commandOccupyPlace = "occupy_place"
	commandStartGame   = "start_game"
)

// Game model.
//
// Shouldn't be created through constructor and NewGame() instead.
type Game struct {
	Id string `json:"id"`

	// All Clients getting updates about this game (includes player1 and player2).
	Clients []Client `json:"clients"`

	// Clients that take part in this game.
	Player1 *Client `json:"player1"`
	Player2 *Client `json:"player2"`

	MovesPlayer1 bool `json:"movesPlayer1"`

	Ball Ball `json:"ball"`

	Moves []Move `json:"moves"`

	GameState string `json:"gameState"`
}

func NewGame() Game {
	return Game{
		Id:           uuid.NewString(),
		Clients:      []Client{},
		MovesPlayer1: true,
		Ball: Ball{
			// TODO give proper values
			X: 10,
			Y: 10,
		},
		Moves:     []Move{},
		GameState: gameStateWaitingForPlayers,
	}
}

func (game Game) ToJsonString() (string, error) {
	byteArray, err := json.Marshal(game)
	if err != nil {
		fmt.Println("Error in Marshal, e: ", err)
		return "", errors.New("Couldn't parse Game")
	}

	return string(byteArray), nil
}

func (game *Game) InformEveryClient() {
	clientsToRemove := []int{}

	gameJsonString, err := game.ToJsonString()

	for i, client := range game.Clients {
		fmt.Println("Informing client: ", client.Id)

		if err != nil {
			msg := getErrorJsonString(errorCouldntParseGameToJson)
			err2 := client.connection.WriteMessage(1, []byte(msg))
			if err2 != nil {
				clientsToRemove = append(clientsToRemove, i)
				continue
			}
		}

		msg := "{\"error\": null, \"game\": " + gameJsonString + "}"
		err2 := client.connection.WriteMessage(1, []byte(msg))
		if err2 != nil {
			clientsToRemove = append(clientsToRemove, i)
			continue
		}
	}

	for _, clientIndex := range clientsToRemove {
		game.RemoveClientUnderIndex(clientIndex)
	}
}

func (game *Game) AddClient(connection *websocket.Conn) {
	newClient := Client{
		connection:  connection,
		Id:          uuid.NewString(),
		SecondsLeft: 60 * 5,
		StartedGame: false,
	}

	game.Clients = append(game.Clients, newClient)

	gameJsonString, err := game.ToJsonString()
	if err != nil {
		msg := "{\"error\": \"invalid_game_json\", \"game\": null, \"your_id\": \"" + newClient.Id + "\"}"
		connection.WriteMessage(1, []byte(msg))
		return
	}

	msg := "{\"error\": null, \"game\": " + gameJsonString + ", \"your_id\": \"" + newClient.Id + "\"}"
	connection.WriteMessage(1, []byte(msg))

	game.InformEveryClient()

	go game.handleMessages(&newClient)
}

func (game *Game) RemoveClientUnderIndex(index int) {
	fmt.Println("Remove client at index: ", index)

	game.Clients[index].connection.Close()
	game.Clients = append(game.Clients[:index], game.Clients[index+1:]...)
}

func (game *Game) RemoveClient(client *Client) {
	fmt.Println("Remove client:", client)

	index := -1
	for i, c := range game.Clients {
		if c.Id == client.Id {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	if game.Player1.Id == client.Id {
		game.Player1 = nil
	} else if game.Player2.Id == client.Id {
		game.Player2 = nil
	}

	game.RemoveClientUnderIndex(index)
}

func (game *Game) handleMessages(client *Client) {
	for {
		_, p, err := client.connection.ReadMessage()
		if err != nil {
			fmt.Println(err)
			game.RemoveClient(client)
			return
		}

		game.handleMessage(client, string(p))
	}
}

func (game *Game) handleMessage(client *Client, msg string) {
	fmt.Println("New message: ", msg)

	var jsonReq map[string]interface{}
	err := json.Unmarshal([]byte(msg), &jsonReq)
	if err != nil {
		msg := getErrorJsonString(errorReceivedInvalidJson)
		client.connection.WriteMessage(1, []byte(msg))
		return
	}

	command := jsonReq["command"]
	val := jsonReq["val"]

	fmt.Println("command:", command)
	fmt.Println("val:", val)
	fmt.Println("val type:", reflect.TypeOf(val))

	if command == commandOccupyPlace {
		game.handleOccupyPlace(client, jsonReq)
		return
	}

	if command == commandStartGame {
		game.handleStartGame(client, jsonReq)
		return
	}
}

// ====== COMMANDS HANDLING

func (game *Game) handleOccupyPlace(client *Client, jsonReq map[string]interface{}) {
	place := jsonReq["val"]

	// ALREADY HAS PLAYERS
	if !(game.GameState == gameStateWaitingForPlayers || game.GameState == gameStateHasOnePlayer) {
		msg := getErrorJsonString(errorGameIsNotWaitingForPlayers)
		client.connection.WriteMessage(1, []byte(msg))
		return
	}

	if place == 1.0 {
		if game.Player1 != nil {
			msg := getErrorJsonString(errorPlaceAlreadyOccupied)
			client.connection.WriteMessage(1, []byte(msg))
			return
		}

		game.Player1 = client

		if game.Player2 != nil {
			game.GameState = gameStateHasTwoPlayers
		} else {
			game.GameState = gameStateHasOnePlayer
		}

		game.InformEveryClient()
		return

	} else if place == 2.0 {
		if game.Player2 != nil {
			msg := getErrorJsonString(errorPlaceAlreadyOccupied)
			client.connection.WriteMessage(1, []byte(msg))
			return
		}

		game.Player2 = client

		if game.Player1 != nil {
			game.GameState = gameStateHasTwoPlayers
		} else {
			game.GameState = gameStateHasOnePlayer
		}

		game.InformEveryClient()
		return

	} else {
		msg := getErrorJsonString(errorReceivedInvalidJson)
		client.connection.WriteMessage(1, []byte(msg))
		return
	}
}

func (game *Game) handleStartGame(client *Client, jsonReq map[string]interface{}) {
	if game.GameState != gameStateHasTwoPlayers && game.GameState != gameStateOnePlayerStarted {
		msg := getErrorJsonString(errorGameCannotBeStartedYet)
		client.connection.WriteMessage(1, []byte(msg))
		return
	}

	client.StartedGame = true

	if game.Player1.StartedGame && !game.Player2.StartedGame {
		// Only 1st player has started.
		game.GameState = gameStateOnePlayerStarted
	} else if game.Player1.StartedGame && !game.Player2.StartedGame {
		// Only 2nd player has started.
		game.GameState = gameStateOnePlayerStarted
	} else if game.Player1.StartedGame && game.Player2.StartedGame {
		// Both players have started.
		game.GameState = gameStateRunning
	}

	game.InformEveryClient()
}

// ====== HELPERS

func getErrorJsonString(err string) string {
	return "{\"error\": \"" + err + "\", \"game\": null}"
}
