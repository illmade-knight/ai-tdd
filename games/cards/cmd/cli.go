package cmd

import (
	"bufio"
	"games/user/store"
	"io"
	"os"
	"strings"
)

type CLI struct {
	s  store.PlayerStore
	in io.Reader
}

func NewParquetCLI() *CLI {
	ps, err := store.NewParquetPlayerStore("tmp.pq")
	if err != nil {
		panic(err)
	}
	cli := &CLI{ps, os.Stdin}
	return cli
}

func (cli CLI) PlayCards() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.s.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
