package main

import (
	"fmt"
	"net/http"
  "math/rand"
  "io/ioutil"
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

gs := NewGameStore()
playURL := "localhost:5000/play"

func main() {

}

func NewGameStore()  *GameStore{
  return &GameStore{}
}


func init()  {

}

func (gs *GameStore) CloseGame(g string) {
	gs.mu.Lock()
  defer gs.mu.Unlock()
  delete(gs.games[g])
}

func (gs *GameStore) Get()  {
  gs.mu.RLock()
  defer gs.mu.RUnlock()
  return gs.games[0]
}

func (b *Bot) Play() {
  g := gs.Get()
  json := []byte(`{"game" : g, "player": b.Name, "action" : b.ChooseAction()}`)
  resp, err := http.Post(playURL, "application/json", bytes.NewBuffer(json))
    if err != nil {
    fmt.Println("I faild to play :( )")
    return
  }
  defer resp.Body.Close()

  body,_ := ioutil.ReadAll(resp.Body)
  fmt.Println(body)

}

func (b *Bot) ChooseAction() string {
  rand.Seed(time.Now().UTC().UnixNano())
  a := [3]string{"rock","paper", "scissor"}
  return a[rand.Intn(3)]
}
