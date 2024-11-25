package main

import (
	"fmt"
	"os"

	"github.com/burgr033/t00lbox/internal/gui"
)

// main runs the program
func main() {
	program, err := gui.CreateProgram("resources.yml")
	if err != nil {
		fmt.Println("Error initializing program:", err)
		os.Exit(1)
	}

	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
