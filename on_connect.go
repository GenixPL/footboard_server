package main

import (
	"fmt"
	"net/http"
)

func OnConnect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connect")
}
