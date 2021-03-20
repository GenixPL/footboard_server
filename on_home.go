package main

import (
	"fmt"
	"net/http"
)

func OnHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/ triggered")

	fmt.Fprintf(w, "Welcome Adventurer!")
}
