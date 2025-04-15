package main

import (
	"games/user/server"
	"log"
	"net/http"
)

func main() {
	s := server.PlayerServer{}
	log.Fatal(http.ListenAndServe(":5000", s.Handler))
}
