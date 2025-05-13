package main

import (
	"games/user/server"
	"games/user/store"
	"log"
	"net/http"
)

func main() {

	st, err := store.NewParquetPlayerStore("tmp.pq")
	if err != nil {
		panic(err)
	}

	sr := server.NewPlayerServer(st)

	log.Fatal(http.ListenAndServe(":5000", sr))

}
