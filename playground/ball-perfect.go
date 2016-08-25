package main

import (
	"fmt"
	"time"
)

type Ball struct {
	hits int
}
func main() {
	table := make(chan *Ball)
	quit1 := player("p1", table)
	quit2 := player("p2", table)
	table <- new(Ball)
	time.Sleep(time.Second)
	<-table
	quit1 <- true
	quit2 <- true
	time.Sleep(time.Millisecond) // <-- extra wait, so the last runnable goes away, feels not good though
	panic("show me the stacks")
}

func player(msg string, table chan *Ball) chan bool {
	quit := make(chan bool)
	go func() {
		for {
			select {
			case ball := <-table:
				ball.hits++
				fmt.Println(msg, ball.hits)
				table <- ball
			case <-quit: 
				return
			}
		}
	}()
	return quit
}