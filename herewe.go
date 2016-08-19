package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
)
var Figures = [3]string { "rock", "paper", "scissor" }

type player struct {
	name string
	figure string
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	human := player{"You", os.Args[1]} // bam if not provided

	computer := player{"Computer", Figures[rand.Intn(len(Figures))]}


	fmt.Println(human.name, "picked", human.figure)
	fmt.Println(computer.name, "picked", computer.figure)

	var winner = determineWinner(&human, &computer)
	if winner != nil {
		fmt.Println("Winner is", winner.name)
	} else {
		fmt.Println("It's a tie!")
	}
}

func determineWinner(player1, player2 *player) *player {

	successors := map[string]string {
		"paper": "scissor",
		"scissor": "rock",
		"rock": "paper",
	}

	if successors[player1.figure] == player2.figure {
		return player2
	} else if successors[player2.figure] == player1.figure {
		return player1
	} else {
		return nil
	}
}
