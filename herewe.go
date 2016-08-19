package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
)
var Figures = [3]string { "rock", "paper", "scissor" }

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	playerFigure := os.Args[1] // bam if not provided

	randFigure := Figures[rand.Intn(len(Figures))]


	fmt.Println("You (player 1)", playerFigure)
	fmt.Println("Computer (player 2)", randFigure)


	fmt.Println("winner is player", determineWinner(playerFigure, randFigure))
}

func determineWinner(figures ...string) string {

	successors := map[string]string {
		"paper": "scissor",
		"scissor": "rock",
		"rock": "paper",
	}

	if successors[figures[0]] == figures[1] {
		return "2"
	} else if successors[figures[1]] == figures[0] {
		return "1"
	} else {
		return "draw"
	}
}
