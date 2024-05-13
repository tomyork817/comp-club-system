package main

import (
	"comp-club-system/internal/parser"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("not enough args")
		return
	}
	filepath := args[0]
	c, err := parser.ReadComputerClub(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.RunIncomingEvents()
	c.Print()
}
