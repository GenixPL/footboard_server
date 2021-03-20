package main

import (
	"fmt"
	"net/http"
)

func OnCreateNewGame(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/create-new-game triggered")

	CreateNewGame()

	fmt.Fprintf(w, "Dupa")
}
