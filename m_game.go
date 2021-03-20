package main

// Game states
const (
	waitingForPlayers = "waiting_for_players"
	hasOnePlayer      = "has_one_player"
	hasTwoPlayers     = "has_two_players"
	onePlayerStarted  = "one_player_started"
	// We skip two_players_started because that's where the running state starts.
	running         = "running"
	paused          = "paused"
	firstPlayerWon  = "first_player_won"
	secondPlayerWon = "second_player_won"
)

// Game model.
//
// Shouldn't be created through constructor and NewGame() instead.
type Game struct {
	// All Clients getting updates about this game (includes player1 and player2).
	Clients []Client `json:"clients"`

	// Clients that take part in this game.
	Player1 Client `json:"player1"`
	Player2 Client `json:"player2"`

	MovesPlayer1 bool `json:"movesPlayer1"`

	Ball Ball `json:"ball"`

	Moves []Move `json:"moves"`

	GameState string `json:"gameState"`
}

func NewGame() Game {
	return Game{
		Clients: []Client{},
		Player1: Client{
			Id: "",
		},
		Player2: Client{
			Id: "",
		},
		MovesPlayer1: true,
		Ball: Ball{
			// TODO give proper values
			X: 10,
			Y: 10,
		},
		Moves:     []Move{},
		GameState: waitingForPlayers,
	}
}
