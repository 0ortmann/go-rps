package main

import (
	"fmt"
)

type Message struct {
	name string
	wait chan bool
}

func main() {
	c := fanIn(boring("foo"), boring("bar"))
	for i := 0; i < 10; i++ {
		msg1 := <-c
		msg2 := <-c
		fmt.Println(msg1.name)
		fmt.Println(msg2.name)
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("Main exit")
}

func fanIn(c1, c2 <-chan Message) <-chan Message {
	c := make(chan Message)
	go func() {
		for {
			select {
			case msg1 := <-c1:
				c <- msg1
			case msg2 := <-c2:
				c <- msg2
			}
		}
	}()
	return c
}

func boring(name string) <-chan Message {
	waitForIt := make(chan bool)
	c := make(chan Message)
	go func() {
		for i := 0; ; i++ {
			c <- Message{name: fmt.Sprintf("%s %d", name, i), wait: waitForIt}
			<-waitForIt
		}
	}()
	return c
}
