package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"bufio"
	"strings"
	"log"
)
var Figures = [3]string { "rock", "paper", "scissor" }

type player struct {
	name string
	figure string
}

func main() {

	
	computer := getComputerPlayer()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please choose a name")
	userName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	userName = strings.TrimSpace(userName)

	fmt.Println("Please draw (rock, scissor, paper)")
	userFigure, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	userFigure = strings.TrimSpace(userFigure)


	human := player{userName, userFigure}



	fmt.Println(human.name, "picked", human.figure)
	fmt.Println(computer.name, "picked", computer.figure)

	var winner = determineWinner(&human, &computer)
	if winner != nil {
		fmt.Println("Winner is", winner.name)
	} else {
		fmt.Println("It's a tie!")
	}
}

func getComputerPlayer() player {
	rand.Seed(time.Now().UTC().UnixNano())
	return player{"Computer", Figures[rand.Intn(len(Figures))]}
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
