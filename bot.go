package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Bot struct {
	Name string
	Game string
	Wins int
}

type GameStore struct {
	games []string
	mu    sync.RWMutex
}

func NewBot(name string) *Bot {
	return &Bot{
		Name: name,
	}
}

const playURL = "http://localhost:5000/play"
const createURL = "http://localhost:5000/create"
const evalURL = "http://localhost:5000/eval"

func main() {
	total := 1000
	fin := make(chan int)
	for i := 0; i < total; i++ {
		initBots("game-" + strconv.Itoa(i), fin)
	}
	var count int
	for {
		count += <-fin
		if count == total {
			return 
		}
	}
}

// Creates a random number of bots between 1 and 10 and lets them play a game, then close that game
func initBots(game string, fin chan int) {
	go func() {
		rand.Seed(time.Now().UTC().UnixNano())
		sendCreate(game)
		bCount := rand.Intn(10) + 1
		done := make(chan int)
		for i := 0; i < bCount; i++ {
			b := NewBot("bot-" + strconv.Itoa(i))
			b.Play(game, done)
		}
		var doneCount int
		for {
			doneCount += <-done
			if doneCount == bCount {
				sendEval(game)
				fin <- 1
				return
			}
		}
	}()
}

func (b *Bot) Play(game string, done chan int) {
	go func() {
		type Payload struct {
			Game   string
			Player string
			Action string
		}
		payload := &Payload{game, b.Name, b.ChooseAction()}
		jsonStr, _ := json.Marshal(payload)

		resp, err := http.Post(playURL, "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Println("I failed to play :( )", err)
			done <- 1
			return
		}
		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))
		resp.Body.Close()
		done <- 1
	}()
}

func (b *Bot) ChooseAction() string {
	rand.Seed(time.Now().UTC().UnixNano())
	a := [3]string{"rock", "paper", "scissor"}
	return a[rand.Intn(3)]
}

func sendEval(game string) {
	type Payload struct {
		Game string
	}
	payload := &Payload{game}
	jsonStr, _ := json.Marshal(payload)
	resp, err := http.Post(evalURL, "application/jsonStr", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("I failed to eval for game", game)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	resp.Body.Close()
}

func sendCreate(game string) {
	type Payload struct {
		Name string
	}
	payload := &Payload{game}
	jsonStr, _ := json.Marshal(payload)
	resp, err := http.Post(createURL, "application/jsonStr", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("I faild to create game", game)
		return
	}
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	resp.Body.Close()
}