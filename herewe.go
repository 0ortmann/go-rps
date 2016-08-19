package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
	fmt.Println("Choose your playmode (single, server, client):")

	playmode := parseConsoleInput()

	switch playmode {
	case "server":
		startServer()
	case "client":
		startClient()
	case "single":
		startSinglePlayer()
	}
}

func startServer() {
	// todo
	fmt.Println("You started a server")
	fmt.Println("One player can now join your server, so you can play against each other")

	openWebserver()
}

func handler(writer http.ResponseWriter, request *http.Request) {

	fmt.Println("Hey, someone registered, now you have to play!")

	owner := getPlayerFromConsole()
	url := string(request.URL.Path[1:])
	opponentName := strings.Split(url, "/")[0]
	opponentFigure := strings.Split(url, "/")[1]
	opponent := player{opponentName, Figures[opponentFigure]}
	winner := determineWinner(&owner, &opponent)

	fmt.Println(owner.name, "picked", owner.figure.name)
	fmt.Println(opponent.name, "picked", opponent.figure.name)
	printWinnerText(winner)
	if winner != nil {
		fmt.Fprintf(writer, "Winner is %s", winner.name)
	} else {
		fmt.Fprintf(writer, "It's a tie!")
	}
}

func openWebserver() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5000", nil)
}

func startClient() {
	// todo
}

func startSinglePlayer() {
	computer := getComputerPlayer()
	human := getPlayerFromConsole()

	fmt.Println(human.name, "picked", human.figure.name)
	fmt.Println(computer.name, "picked", computer.figure.name)

	winner := determineWinner(&human, &computer)
	printWinnerText(winner)
}

func getFigureFromConsole() figure {

	var userFigure figure
	for userFigure == Figures["no"] {
		fmt.Println("Please draw (rock, scissor, paper)")
		userFigureName := parseConsoleInput()

		userFigure = Figures[userFigureName]
		if userFigure == Figures["no"] {
			fmt.Println("Choose again! Your choice has not been recognized")
		}
	}
	return userFigure
}

func getPlayerFromConsole() player {

	fmt.Println("Please choose a name")
	userName := parseConsoleInput()
	userFigure := getFigureFromConsole()

	return player{userName, userFigure}
}

func getComputerPlayer() player {
	rand.Seed(time.Now().UTC().UnixNano())
	keys := []string{}
	for k := range Figures {
		keys = append(keys, k)
	}
	randomKey := keys[rand.Intn(len(Figures))]
	return player{"Computer", Figures[randomKey]}
}

func determineWinner(player1, player2 *player) *player {
	switch {
	case player1.figure.name == player2.figure.successorName:
		return player1
	case player2.figure.name == player1.figure.successorName:
		return player2
	}
	return nil
}

func parseConsoleInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(input)
}

func printWinnerText(winner *player) {
	if winner != nil {
		fmt.Println("Winner is", winner.name, winner.figure.ascii)
	} else {
		fmt.Println("It's a tie!")
	}
}
