package main

import (
	"fmt"
	"os"

	"github.com/piotr1215/aicmdtools/internal/aicompgraph"
)

var version = "v0.0.1"

func main() {

	err := aicompgraph.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}

}
