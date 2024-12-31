package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/piotr1215/aicmdtools/internal/aicmd"
	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

var version = "v0.0.179"
var prompt_file = "prompt.txt"

// main is the entry point for the Goai command-line tool.
// It parses command-line flags and executes appropriate actions.
// If the "version" flag is set, it displays the version information and changelog.
// If the "version" flag is not set, it executes the command specified in the "prompt.txt" file.
// If an error occurs during execution, it prints an error message and exits with a non-zero status code.
func main() {
	versionFlag := flag.Bool("version", false, "Display version information")
	modelFlag := flag.Bool("model", false, "Display current model")
	flag.Parse()

	if *modelFlag {
		conf, _, err := config.ReadAndParseConfig("config.yaml", prompt_file)
		if err != nil {
			fmt.Printf("Error reading configuration: %v\n", err)
			os.Exit(-1)
		}
		fmt.Printf("Current model: %s\n", conf.Model)
		return
	}

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
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(-1)
	}
}
