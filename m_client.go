package main

import "github.com/gorilla/websocket"

type Client struct {
	connection *websocket.Conn

	// If it's an empty string then the object should be treated as nil.
	Id string `json:"id"`

	SecondsLeft uint16 `json:"secondsLeft"`
}
