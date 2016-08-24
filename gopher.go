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

type Gopher struct {
	Name string
	Game string
	Wins int
}

type GameStore struct {
	games []string
	mu    sync.RWMutex
}

func NewGopher(name string) *Gopher {
	return &Gopher{
		Name: name,
	}
}

const playURL = "http://localhost:5000/play"
const createURL = "http://localhost:5000/create"
const evalURL = "http://localhost:5000/eval"
const totalGames = 1000

func main() {
	fin := make(chan int)
	for i := 0; i < totalGames; i++ {
		go initGophers("game-" + strconv.Itoa(i), fin)
	}
	for i := 0; i < totalGames; i++ {
		<-fin
	}
}

// Creates a random number of gopher between 1 and 10 and lets them play a game, then close that game
func initGophers(game string, fin chan int) {
	rand.Seed(time.Now().UTC().UnixNano())
	sendCreate(game)
	gCount := rand.Intn(10) + 1
	done := make(chan int)
	for i := 0; i < gCount; i++ {
		g := NewGopher("gopher-" + strconv.Itoa(i))
		go g.Play(game, done)
	}
	for i := 0; i < gCount; i++{
		<-done
	}
	sendEval(game)
	fin <- 1
}

func (g *Gopher) Play(game string, done chan int) {
	type Payload struct {
		Game   string
		Player string
		Action string
	}
	payload := &Payload{game, g.Name, g.ChooseAction()}
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
}

func (b *Gopher) ChooseAction() string {
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