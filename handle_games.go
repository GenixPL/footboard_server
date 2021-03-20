package main

import (
	"encoding/json"
	"fmt"
	"time"
)

var Games []Game = []Game{}

// Periodically removes games without clients.
func StartPeriodicEmptyGamesRemoval() {
	duration := time.Duration(5) * time.Minute

	ticker := time.NewTicker(duration)

	for range ticker.C {
		removeEmptyGames()
	}
}

// Removes games that don't have any clients.
func removeEmptyGames() {
	fmt.Println("Removing empty games...")
}

// Creates new game.
func CreateNewGame() {
	fmt.Println("Creating new game...")

	Games = append(Games, NewGame())
}

func GamesToJsonStr() string {
	str := "["

	for i := 0; i < len(Games); i++ {
		byteArray, err := json.Marshal(Games[i])
		if err != nil {
			fmt.Println("Error in Marshal, e: ", err)
			continue
		}

		if i != 0 {
			str += ", "
		}

		str += string(byteArray)
	}

	str += "]"

	return str
}
