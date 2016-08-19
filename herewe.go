package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
)

func main() {
	
	playerFigure := os.Args[1] // bam if not provided

	figures := []string { "rock", "paper", "scissor" }

	rand.Seed(time.Now().UTC().UnixNano())

	randFigure := figures[rand.Intn(len(figures))]
	fmt.Println("You (player 1)", playerFigure)
	fmt.Println("Computer (player 2)", randFigure)


	fmt.Println("winner is player", determineWinner(playerFigure, randFigure))
}

func determineWinner(figure1, figure2 string) string {

	successors := map[string]string {
		"paper": "scissor",
		"scissor": "rock",
		"rock": "paper",
	}

	if successors[figure1] == figure2 {
		return "2"
	} else if successors[figure2] == figure1 {
		return "1"
	} else {
		return "draw"
	}
}

