package main

import (
	"fmt"
	"games/cards/cmd"
)

func main() {
	fmt.Println("record a card win")
	fmt.Println("say 'foobar wins' and the given name has a score recorded")
	fmt.Println()

	cli := cmd.NewParquetCLI()
	cli.PlayCards()
}
