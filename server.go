package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type GameStore struct {
	games map[string]*Game
	mu    sync.RWMutex
}

type Game struct {
	Name    string
	Open    bool
	Players map[string]string
	Winners []string
	mu      sync.RWMutex
}

type PlayerAction struct {
	Game   string
	Player string
	Action string
}

func (gs *GameStore) Get(key string) *Game {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.games[key]
}

func (gs *GameStore) Set(key string, game *Game) bool {
	if gs.Get(key) != nil {
		return false
	}
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.games[key] = game
	return true
}

func (g *Game) GetAction(key string) (p string, ok bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	p, ok = g.Players[key]
	return
}

func (g *Game) SetAction(p string, a string) bool {
	if p == "" {
		return false
	}
	if _, ok := g.GetAction(p); ok {
		return false
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Players[p] = a
	return true
}

func (g *Game) Eval() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Open = false
	r := make(map[string][]string)
	for p, a := range g.Players {
		r[a] = append(r[a], p)
	}

	if len(r) == 1 || len(r) == 3 {
		// its a draw
		return
	}
	switch {
	case r["paper"] != nil && r["scissor"] != nil:
		g.Winners = r["scissor"]
	case r["scissor"] != nil && r["rock"] != nil:
		g.Winners = r["rock"]
	case r["rock"] != nil && r["paper"] != nil:
		g.Winners = r["paper"]
	}

	return

}

func NewGameStore() *GameStore {
	return &GameStore{
		games: make(map[string]*Game),
	}
}

var gs = NewGameStore()

func main() {
	http.HandleFunc("/create", checkPost(createHandler))
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/play", checkPost(actionHandler))
	http.HandleFunc("/eval", checkPost(evalHandler))
	http.ListenAndServe(":5000", nil)
}

func checkPost(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Wrong method", 405)
			return
		}
		f(w, req)
	}
}

func createHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var game Game
	err := decoder.Decode(&game)
	if err != nil {
		fmt.Fprint(w, "Fatal error parsing request")
		return
	}
	game.Open = true
	game.Players = make(map[string]string)
	if !gs.Set(game.Name, &game) {
		http.Error(w, "Already exists", 409)
		return
	}
	fmt.Fprint(w, "OK")
}

func gameHandler(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("game")
	game := gs.Get(name)
	if game == nil {
		http.Error(w, "Game not found", 404)
		return
	}
	fmt.Fprintf(w, "The open status of the game %s is %b", name, game.Open)
}

func actionHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var pa PlayerAction
	err := decoder.Decode(&pa)
	if err != nil {
		fmt.Fprintf(w, "Alles im Arsch Bruder!")
		return
	}
	game := gs.Get(pa.Game)
	if game == nil {
		http.Error(w, "Game not found", 404)
		return
	}
	if !game.Open {
		http.Error(w, "Game already closed", 409)
		return
	}
	if !game.SetAction(pa.Player, pa.Action) {
		http.Error(w, "Player name already exists or empty", 409)
		return
	}
	fmt.Fprintf(w, "Game %s, %s/%s", game.Name, pa.Player, pa.Action)
}

func evalHandler(w http.ResponseWriter, req *http.Request) {
	type Eval struct {
		Game string
	}
	decoder := json.NewDecoder(req.Body)
	var e Eval
	err := decoder.Decode(&e)
	if err != nil {
		fmt.Fprint(w, "Alles kapr0tt")
		return
	}
	game := gs.Get(e.Game)
	if game == nil {
		http.Error(w, "Game not found", 404)
		return
	}
	game.Eval()
	fmt.Fprintf(w, "Game winner(s) of %s: %s\nParticipants: %s", game.Name, game.Winners, game.Players)
}