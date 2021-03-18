package main

import (
	"fmt"
	"net/http"
)

func OnHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("home")
}
