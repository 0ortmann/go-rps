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
	go player("p1", table)
	go player("p2", table)
	table <- new(Ball)
	time.Sleep(time.Second)
	<-table
	panic("show me the stacks")
}

func player(msg string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		fmt.Println(msg, ball.hits)
		table <- ball
	}
}