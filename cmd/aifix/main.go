package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/piotr1215/aicmdtools/internal/aifix"
	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

var version = "v0.0.2"
var promptFile = "aifix-prompt.txt"

func main() {
	versionFlag := flag.Bool("version", false, "Display version information")
	modelFlag := flag.Bool("model", false, "Display current model")
	helpFlag := flag.Bool("help", false, "Display help information")
	initShellFlag := flag.String("init-shell", "", "Initialize shell integration (bash, zsh, or fish)")
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	if *initShellFlag != "" {
		showShellInit(*initShellFlag)
		return
	}

	if *modelFlag {
		conf, _, err := config.ReadAndParseConfig("config.yaml", promptFile)
		if err != nil {
			fmt.Printf("Error reading configuration: %v\n", err)
			os.Exit(-1)
		}
		fmt.Printf("Current model: %s\n", conf.Model)
		return
	}

	if *versionFlag {
		fmt.Printf("aifix version: %s\n", version)
		changelog, err := utils.GenerateChangelog(exec.Command)
		if err != nil {
			fmt.Printf("Error generating changelog: %v\n", err)
		} else {
			fmt.Printf("\nChangelog:\n%s", changelog)
		}
		return
	}

	// Get manual error input if provided
	var manualError string
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		manualError = strings.Join(os.Args[1:], " ")
	}

	err := aifix.Execute(promptFile, manualError)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}
}

func showHelp() {
	help := `aifix - Instant Error Explanation and Fix Suggestion

USAGE:
  aifix [error message]          Analyze error and suggest fixes
  aifix -version                 Display version information
  aifix -model                   Display current AI model
  aifix -help                    Display this help message
  aifix -init-shell <shell>      Show shell integration setup

EXAMPLES:
  # Analyze last command error automatically
  $ go build
  # error: undefined: fmt.Println
  $ aifix

  # Provide error message directly
  $ aifix "Module not found: 'react-dom'"

  # Show current model
  $ aifix -model

SHELL INTEGRATION (Optional):
  For automatic error detection, add to your shell config:

  # For zsh users (~/.zshrc):
  $ aifix -init-shell zsh

  # For bash users (~/.bashrc):
  $ aifix -init-shell bash

  # For fish users (~/.config/fish/config.fish):
  $ aifix -init-shell fish

CONFIGURATION:
  Config file: ~/.config/aicmdtools/config.yaml
  - provider: AI provider (openai or anthropic)
  - model: Model to use
  - temperature: Response randomness (0-1)
  - max_tokens: Maximum response length

For more information, visit: https://github.com/piotr1215/aicmdtools
`
	fmt.Print(help)
}

func showShellInit(shell string) {
	switch shell {
	case "zsh":
		fmt.Println(`# aifix - ZSH Integration
# Add this to your ~/.zshrc

# Error capture for automatic detection
export AIFIX_CMD_FILE="/tmp/aifix_last_cmd_$$"
export AIFIX_ERROR_FILE="/tmp/aifix_last_error_$$"

# Redirect stderr to capture file AND terminal
exec 2> >(tee -a "$AIFIX_ERROR_FILE" >&2)

# Capture command and error
precmd() {
    local exit_code=$?
    # Don't capture aifix itself
    if [[ "$AIFIX_LAST_CMD" =~ ^(aifix|fix) ]]; then
        return
    fi

    if [ $exit_code -ne 0 ] && [ -n "$AIFIX_LAST_CMD" ]; then
        echo "$AIFIX_LAST_CMD" > "$AIFIX_CMD_FILE"
        export AIFIX_LAST_EXIT=$exit_code
    else
        rm -f "$AIFIX_CMD_FILE" "$AIFIX_ERROR_FILE" 2>/dev/null
        unset AIFIX_LAST_EXIT
    fi
}

preexec() {
    export AIFIX_LAST_CMD="$1"
    # Clear error file before new command
    > "$AIFIX_ERROR_FILE"
}

# Quick alias
alias fix='aifix'

# Optional: bind to Esc-Esc
aifix-command-line() {
    BUFFER="aifix"
    zle accept-line
}
zle -N aifix-command-line
bindkey '\e\e' aifix-command-line
`)
	case "bash":
		fmt.Println(`# Add to ~/.bashrc for automatic error detection

export AIFIX_CMD_FILE="/tmp/aifix_last_cmd_$$"
export AIFIX_ERROR_FILE="/tmp/aifix_last_error_$$"

# Redirect stderr to capture file AND terminal
exec 2> >(tee -a "$AIFIX_ERROR_FILE" >&2)

# Capture command before execution
trap 'AIFIX_LAST_CMD="$BASH_COMMAND"; > "$AIFIX_ERROR_FILE"' DEBUG

# Capture exit code after execution
PROMPT_COMMAND='
    exit_code=$?
    if [[ "$AIFIX_LAST_CMD" =~ ^(aifix|fix) ]]; then
        return
    fi
    if [ $exit_code -ne 0 ] && [ -n "$AIFIX_LAST_CMD" ]; then
        echo "$AIFIX_LAST_CMD" > "$AIFIX_CMD_FILE"
        export AIFIX_LAST_EXIT=$exit_code
    else
        rm -f "$AIFIX_CMD_FILE" "$AIFIX_ERROR_FILE" 2>/dev/null
        unset AIFIX_LAST_EXIT
    fi
'"${PROMPT_COMMAND:+; $PROMPT_COMMAND}"

# Quick alias
alias fix='aifix'
`)
	case "fish":
		fmt.Println(`# Add to ~/.config/fish/config.fish for automatic error detection

set -gx AIFIX_CMD_FILE "/tmp/aifix_last_cmd_"(echo %self)
set -gx AIFIX_ERROR_FILE "/tmp/aifix_last_error_"(echo %self)

# Redirect stderr to capture file AND terminal
function fish_prompt_setup
    exec 2> >(tee -a $AIFIX_ERROR_FILE >&2)
end
fish_prompt_setup

function aifix_preexec --on-event fish_preexec
    # Clear error file before new command
    echo -n > $AIFIX_ERROR_FILE
end

function aifix_capture --on-event fish_postexec
    set -g exit_code $status
    set -g AIFIX_LAST_CMD "$argv"

    # Skip aifix commands
    if string match -q -r '^(aifix|fix)' "$argv"
        return
    end

    if test $exit_code -ne 0
        echo "$argv" > $AIFIX_CMD_FILE
        set -gx AIFIX_LAST_EXIT $exit_code
    else
        rm -f $AIFIX_CMD_FILE $AIFIX_ERROR_FILE 2>/dev/null
        set -e AIFIX_LAST_EXIT
    end
end

# Quick alias
alias fix='aifix'
`)
	default:
		fmt.Printf("Unknown shell: %s\n", shell)
		fmt.Println("Supported shells: bash, zsh, fish")
		os.Exit(1)
	}
}
