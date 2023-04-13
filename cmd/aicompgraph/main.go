package main

import (
	"fmt"
	"os"

	"github.com/piotr1215/aicmdtools/internal/aicompgraph"
)

func main() {

	err := aicompgraph.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}

}
