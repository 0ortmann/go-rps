package main

import (
	"fmt"
	"os"
)

func main() {
	
	fmt.Println("Hello")
	passedArgs := os.Args[1:]
	fmt.Println("I got the following args passed: ", passedArgs)
}

