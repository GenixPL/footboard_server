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

type Game struct {
	// All Clients getting updates about this game (includes player1 and player2).
	clients []Client

	// Clients that take part in this game.
	player1 Client
	player2 Client

	movesPlayer1 bool

	ball Ball

	moves []Move

	gameState string
}

func NewGame() Game {
	return Game{
		clients: []Client{},
		player1: Client{
			id: "",
		},
		player2: Client{
			id: "",
		},
		movesPlayer1: true,
		ball: Ball{
			// TODO give proper values
			x: 10,
			y: 10,
		},
		moves:     []Move{},
		gameState: waitingForPlayers,
	}
}
