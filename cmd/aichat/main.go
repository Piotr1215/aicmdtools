package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/piotr1215/aicmdtools/internal/aichat"
)

var version = "v0.0.1"

func main() {

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	err := aichat.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}
}
