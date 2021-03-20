package main

import (
	"fmt"
	"net/http"
)

func OnGetGames(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/games triggered")

	fmt.Fprintf(w, GamesToJsonStr())
}
