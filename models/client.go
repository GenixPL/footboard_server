package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Connection *websocket.Conn

	Id          string `json:"id"`
	SecondsLeft uint16 `json:"secondsLeft"`
	StartedGame bool   `json:"startedGame"`
}

func NewClient(connection *websocket.Conn) Client {
	return Client{
		Connection:  connection,
		Id:          uuid.NewString(),
		SecondsLeft: 60 * 5,
		StartedGame: false,
	}
}
