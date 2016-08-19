package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var Figures = map[string]figure{
	"rock":    figure{"rock", "asciigen(rock)", "paper"},
	"paper":   figure{"paper", "asciigen(paper)", "scissor"},
	"scissor": figure{"scissor", "asciigen(scissor)", "rock"},
}

type figure struct {
	name          string
	ascii         string
	successorName string
}

type player struct {
	name   string
	figure figure
}

func main() {

	computer := getComputerPlayer()
	human := getPlayerFromConsole()

	fmt.Println(human.name, "picked", human.figure.name)
	fmt.Println(computer.name, "picked", computer.figure.name)

	var winner = determineWinner(&human, &computer)
	if winner != nil {
		fmt.Println("Winner is", winner.name, winner.figure.ascii)
	} else {
		fmt.Println("It's a tie!")
	}
}

func getFigure() figure {
	var userFigureName string
	var userFigure figure
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please draw (rock, scissor, paper)")
	userFigureName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	userFigureName = strings.TrimSpace(userFigureName)
	userFigure = Figures[userFigureName] // this may break
	return userFigure
}

func getPlayerFromConsole() player {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please choose a name")
	userName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	userName = strings.TrimSpace(userName)
	var userFigure figure
	for userFigure == Figures["no"] {
		userFigure = getFigure()
		if userFigure == Figures["no"] {
			fmt.Println("Choose again! Your choice has not been recognized")
		}
	}
	return player{userName, userFigure}
}

func getComputerPlayer() player {
	rand.Seed(time.Now().UTC().UnixNano())
	keys := []string{}
	for k := range Figures {
		keys = append(keys, k)
	}
	fmt.Println(keys)
	randomKey := keys[rand.Intn(len(Figures))]
	return player{"Computer", Figures[randomKey]}
}

func determineWinner(player1, player2 *player) *player {

	if player1.figure.name == player2.figure.successorName {
		return player1
	} else if player2.figure.name == player1.figure.successorName {
		return player2
	} else {
		return nil
	}
}
