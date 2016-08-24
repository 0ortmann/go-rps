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

func NewGameStore() *GameStore {
	return &GameStore{}
}

var gs = NewGameStore()

const playURL = "http://localhost:5000/play"
const createURL = "http://localhost:5000/create"
const evalURL = "http://localhost:5000/eval"

func main() {
	initGameStore()

	for i := 0; i < len(gs.games); i++ {
		initBots()
	}
	time.Sleep(time.Second)
}

// Creates a random number of bots between 1 and 10 and lets them play a game, then close that game
func initBots() {
	rand.Seed(time.Now().UTC().UnixNano())
	game := gs.Get()
	sendCreate(game)
	for i := 0; i < rand.Intn(10); i++ {
		b := NewBot("bot-" + strconv.Itoa(i))
		b.Play()
	}
	time.Sleep(time.Second)
	sendEval(game)
	gs.CloseGame()
}

func initGameStore() {
	for i := 0; i < 5; i++ {
		name := "game-" + strconv.Itoa(i+1)
		fmt.Println("Created game: ", name)
		gs.Set(name)
	}
}

func (gs *GameStore) CloseGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.games[0] = gs.games[len(gs.games)-1]
	gs.games = gs.games[:len(gs.games)-1]
}

func (gs *GameStore) Get() string {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.games[0]
}

func (gs *GameStore) Set(game string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.games = append(gs.games, game)
}

func (b *Bot) Play() {
	type Payload struct {
		Game   string
		Player string
		Action string
	}
	payload := &Payload{gs.Get(), b.Name, b.ChooseAction()}
	jsonStr, _ := json.Marshal(payload)

	fmt.Println(string(jsonStr))

	resp, err := http.Post(playURL, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("I faild to play :( )", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
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
		fmt.Println("I faild to eval for game", game)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func sendCreate(game string) {
	fmt.Println("Shall send create for ", game)
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
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
