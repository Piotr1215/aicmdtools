package aicompgraph

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/nlp"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

var version = "v0.0.1"
var prompt_file = "comp-graph-prompt.txt"

func Execute() error {

	versionFlag := flag.Bool("version", false, "Display version information")
	fileFlag := flag.String("f", "", "Path to YAML file")

	if *versionFlag {
		fmt.Printf("aicompgraph version: %s\n", version)
		changelog, err := utils.GenerateChangelog(exec.Command)
		if err != nil {
			fmt.Printf("Error generating changelog: %v\n", err)
		} else {
			fmt.Printf("\nChangelog:\n%s", changelog)
		}
		return err
	}

	flag.Parse()
	if *fileFlag == "" {
		fmt.Println("Error: No YAML file path specified. Use the -f flag to provide a file path.")
		os.Exit(-1)
	}

	yamlFilePath := *fileFlag
	yamlFileContent, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		os.Exit(-1)
	}

	userPrompt := string(yamlFileContent)

	configReader := &utils.FileReader{
		FilePathFunc: func() string { return config.ConfigFilePath("config.yaml") },
	}
	configContent := configReader.ReadFile()
	conf := config.ParseConfig(configContent)

	promptReader := &utils.FileReader{
		FilePathFunc: func() string { return config.ConfigFilePath(prompt_file) },
	}
	prompt := promptReader.ReadFile()
	operating_system, shell := utils.DetectOSAndShell()
	prompt = utils.ReplacePlaceholders(prompt, operating_system, shell)

	client := nlp.CreateOpenAIClient(conf)

	aiClient := nlp.GoaiClient{
		Client: client,
		Prompt: prompt,
	}

	response, err := aiClient.ProcessCommand(userPrompt)
	if err != nil {
		fmt.Printf("Error processing command: %v\n", err)
		return err
	}

	command := response.Choices[0].Message.Content
	fmt.Printf("Command: %s\n", command)

	return nil
}
