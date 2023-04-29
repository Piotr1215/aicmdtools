package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/piotr1215/aicmdtools/internal/aicmd"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

var version = "v0.0.107"
var prompt_file = "prompt.txt"

func main() {
	versionFlag := flag.Bool("version", false, "Display version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Goai version: %s\n", version)
		changelog, err := utils.GenerateChangelog(exec.Command)
		if err != nil {
			fmt.Printf("Error generating changelog: %v\n", err)
		} else {
			fmt.Printf("\nChangelog:\n%s", changelog)
		}
		return
	}
	err := aicmd.Execute(prompt_file)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}
}
