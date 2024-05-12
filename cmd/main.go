package main

import (
	"comp-club-system/internal/parser"
	"fmt"
)

func main() {
	c, err := parser.ReadComputerClub("/home/nikita/comp-club-system/club.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.RunIncomingEvents()
	c.Print()
}
