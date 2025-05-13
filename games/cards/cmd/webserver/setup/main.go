package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	testServer := "http://localhost:5000"

	playerName := "PlayerOne"
	putUrl := fmt.Sprintf("%s/players/%s/wins", testServer, playerName)

	// Record win 1
	r, err := http.NewRequest(http.MethodPut, putUrl, nil)

	if err != nil {
		fmt.Println(err)
	}

	_, err = http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("added")

}
