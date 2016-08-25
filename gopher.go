package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Gopher struct {
	Name string
	Game string
}

type GameResult struct {
	Name      string
	Open      bool
	Players   map[string]string
	Winners   []string
	WinAction string
}

type Stats struct {
	Figure string
	Amount int
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
	stats := make(chan *Stats)
	for i := 0; i < totalGames; i++ {
		go startGame("game-"+strconv.Itoa(i), stats)
		if i%100 == 0 {
			// dont fuck with the ulimits, dont burn the fan...
			time.Sleep(time.Millisecond * 50)
		}
	}
	res := make(map[string]int)
	for i := 0; i < totalGames; i++ {
		stat := <-stats
		res[stat.Figure] += stat.Amount
	}
	fmt.Printf("Total win stats for %d games\n", totalGames)
	for f, w := range res {
		fmt.Println(f, w)
	}
}

// Creates a random number of gophers between 2 and 10 and lets them play a game, then close that game
func startGame(game string, stats chan *Stats) {
	rand.Seed(time.Now().UTC().UnixNano())
	sendCreate(game)
	gCount := rand.Intn(9) + 2
	done := make(chan int)
	for i := 0; i < gCount; i++ {
		g := NewGopher(game + "-gopher-" + strconv.Itoa(i))
		go g.Play(game, done)
	}
	for i := 0; i < gCount; i++ {
		<-done
	}
	r := sendEval(game)
	aggregateStats(r, stats)
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
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	resp.Body.Close()
	done <- 1
}

func (b *Gopher) ChooseAction() string {
	rand.Seed(time.Now().UTC().UnixNano())
	a := [3]string{"rock", "paper", "scissor"}
	return a[rand.Intn(3)]
}

func sendEval(game string) *GameResult {
	type Payload struct {
		Game string
	}
	p := &Payload{game}
	jsonStr, _ := json.Marshal(p)
	resp, err := http.Post(evalURL, "application/jsonStr", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("I failed to eval for game", game)
		return nil
	}
	var r GameResult
	d := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = d.Decode(&r)
	if err != nil {
		fmt.Println("Cannot parse eval response for game", game)
		return nil
	}
	return &r

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
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

func aggregateStats(r *GameResult, stats chan *Stats) {
	if r == nil {
		stats <- &Stats{Figure: "error", Amount: 1}
		return
	}
	if len(r.Winners) == 0 {
		stats <- &Stats{Figure: "tie", Amount: 1}
	} else {
		stats <- &Stats{Figure: r.WinAction, Amount: 1}
	}
}
