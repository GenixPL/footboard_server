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
	errorInvalidJson          = "invalid_json"
	errorPlaceAlreadyOccupied = "place_already_occupied"
)

// Commands
const (
	commandOccupyPlace = "occupy_place"
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

	for i, client := range game.Clients {
		gameJsonString, err := game.ToJsonString()
		if err != nil {
			continue
		}

		msg := "{\"error\": null, \"game\": " + gameJsonString + "}"
		err2 := client.connection.WriteMessage(1, []byte(msg))
		if err2 != nil {
			clientsToRemove = append(clientsToRemove, i)
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
	}

	game.Clients = append(game.Clients, newClient)

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
		msg := getErrorJsonString(errorInvalidJson)
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
}

// ====== COMMANDS HANDLING

func (game *Game) handleOccupyPlace(client *Client, jsonReq map[string]interface{}) {
	place := jsonReq["val"]

	if place == 1.0 {
		if game.Player1 != nil {
			msg := getErrorJsonString(errorPlaceAlreadyOccupied)
			client.connection.WriteMessage(1, []byte(msg))
			return
		}

		game.Player1 = client
		game.InformEveryClient()
		return

	} else if place == 2.0 {
		if game.Player2 != nil {
			msg := getErrorJsonString(errorPlaceAlreadyOccupied)
			client.connection.WriteMessage(1, []byte(msg))
			return
		}

		game.Player2 = client
		game.InformEveryClient()
		return

	} else {
		msg := getErrorJsonString(errorInvalidJson)
		client.connection.WriteMessage(1, []byte(msg))
		return
	}
}

// ====== HELPERS

func getErrorJsonString(err string) string {
	return "{\"error\": \"" + err + "\", \"game\": null}"
}
