package aicmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/nlp"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

func shouldExecuteCommand(config *config.Config, reader io.Reader) bool {
	if !config.Safety {
		return true
	}

	fmt.Print("Execute the command? [Enter/n] ==> ")
	var answer string
	_, _ = fmt.Fscanln(reader, &answer)

	return strings.ToUpper(answer) != "N"
}
func Execute(prompt_file string) error {

	conf, prompt, err := config.ReadAndParseConfig("config.yaml", prompt_file)
	if err != nil {
		fmt.Printf("Error reading and parsing configuration: %v\n", err)
		os.Exit(-1)
	}
	operating_system, shell := utils.DetectOSAndShell()
	prompt = utils.ReplacePlaceholders(prompt, operating_system, shell)

	client := nlp.CreateOpenAIClient(*conf)

	aiClient := nlp.GoaiClient{
		Client: client,
		Prompt: prompt,
	}

	if len(os.Args) < 2 {
		fmt.Println("No user prompt specified.")
		os.Exit(-1)
	}

	userPrompt := strings.Join(os.Args[1:], " ")

	response, err := aiClient.ProcessCommand(userPrompt)
	if err != nil {
		fmt.Printf("Error processing command: %v\n", err)
		return err
	}

	command := response.Choices[0].Message.Content
	fmt.Printf("Command: %s\n", command)

	execute := true
	if conf.Safety {
		execute = shouldExecuteCommand(conf, os.Stdin)
	}

	if execute {
		var cmd *exec.Cmd
		// Use "sh -c" for Unix-like systems and "cmd /C" for Windows
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
