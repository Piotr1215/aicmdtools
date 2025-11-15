package aifix

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/piotr1215/aicmdtools/internal/nlp"
	"github.com/piotr1215/aicmdtools/internal/utils"
)

const (
	maxHistoryLines = 100
	maxErrorLines   = 50
)

type ErrorContext struct {
	Command    string
	Error      string
	Shell      string
	OS         string
	ExitCode   int
	Timestamp  string
	RecentCmds []string
}

// GetShellHistoryFile returns the history file path for the current shell
func GetShellHistoryFile(shell string) string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	homeDir := usr.HomeDir

	switch shell {
	case "zsh":
		return filepath.Join(homeDir, ".zsh_history")
	case "bash":
		return filepath.Join(homeDir, ".bash_history")
	case "fish":
		return filepath.Join(homeDir, ".local/share/fish/fish_history")
	default:
		// Try bash as fallback
		bashHistory := filepath.Join(homeDir, ".bash_history")
		if _, err := os.Stat(bashHistory); err == nil {
			return bashHistory
		}
		// Try zsh as fallback
		zshHistory := filepath.Join(homeDir, ".zsh_history")
		if _, err := os.Stat(zshHistory); err == nil {
			return zshHistory
		}
		return ""
	}
}

// GetLastCommand retrieves the last command from shell history, skipping aifix commands
func GetLastCommand(shell string) (string, error) {
	historyFile := GetShellHistoryFile(shell)
	if historyFile == "" {
		return "", fmt.Errorf("could not determine shell history file")
	}

	file, err := os.Open(historyFile)
	if err != nil {
		return "", fmt.Errorf("error opening history file: %v", err)
	}
	defer file.Close()

	var commands []string
	scanner := bufio.NewScanner(file)

	// For zsh history format: : timestamp:0;command
	// For bash: just command
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var cmd string
		// Handle zsh extended history format
		if strings.HasPrefix(line, ":") {
			parts := strings.SplitN(line, ";", 2)
			if len(parts) == 2 {
				cmd = parts[1]
			}
		} else {
			cmd = line
		}

		if cmd != "" {
			commands = append(commands, strings.TrimSpace(cmd))
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading history file: %v", err)
	}

	// Find last command that's not aifix/fix
	for i := len(commands) - 1; i >= 0; i-- {
		cmd := commands[i]
		// Skip aifix, fix, and variations
		if !strings.HasPrefix(cmd, "aifix") &&
			!strings.HasPrefix(cmd, "fix") &&
			cmd != "aifix" &&
			cmd != "fix" {
			return cmd, nil
		}
	}

	return "", fmt.Errorf("no commands found in history (excluding aifix)")
}

// GetRecentCommands retrieves the last N commands from shell history
func GetRecentCommands(shell string, count int) ([]string, error) {
	historyFile := GetShellHistoryFile(shell)
	if historyFile == "" {
		return nil, fmt.Errorf("could not determine shell history file")
	}

	file, err := os.Open(historyFile)
	if err != nil {
		return nil, fmt.Errorf("error opening history file: %v", err)
	}
	defer file.Close()

	var commands []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var cmd string
		// Handle zsh extended history format
		if strings.HasPrefix(line, ":") {
			parts := strings.SplitN(line, ";", 2)
			if len(parts) == 2 {
				cmd = parts[1]
			}
		} else {
			cmd = line
		}

		if cmd != "" {
			commands = append(commands, strings.TrimSpace(cmd))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading history file: %v", err)
	}

	// Return last N commands
	if len(commands) > count {
		return commands[len(commands)-count:], nil
	}
	return commands, nil
}

// CaptureLastError attempts to capture the last command's error output
// This is a best-effort approach since we can't reliably capture stderr without shell integration
func CaptureLastError(command string, shell string, os string) (string, error) {
	// Try to execute the command to reproduce the error
	// This is safe because we're not actually changing anything, just reading output
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if os == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		// Run with interactive shell flags to load aliases
		if shell == "zsh" {
			cmd = exec.CommandContext(ctx, shell, "-i", "-c", command)
		} else if shell == "bash" {
			cmd = exec.CommandContext(ctx, shell, "-i", "-c", command)
		} else {
			cmd = exec.CommandContext(ctx, shell, "-c", command)
		}
	}

	output, err := cmd.CombinedOutput()

	// Debug: log what we got
	if err != nil {
		// Command failed, which is expected
		return string(output), nil
	}

	// Command succeeded - check if there's still error-like output
	outputStr := string(output)
	if len(outputStr) > 0 {
		// If there's output, assume it's an error message
		return outputStr, nil
	}

	// No output and no error
	return "", fmt.Errorf("command succeeded with no error")
}

// TruncateError limits error output to prevent overwhelming the AI
func TruncateError(errorText string, maxLines int) string {
	lines := strings.Split(errorText, "\n")
	if len(lines) <= maxLines {
		return errorText
	}

	// Take first 20 lines and last 20 lines to capture both root cause and final error
	firstPart := strings.Join(lines[:20], "\n")
	lastPart := strings.Join(lines[len(lines)-20:], "\n")

	return fmt.Sprintf("%s\n\n... [%d lines omitted] ...\n\n%s",
		firstPart, len(lines)-40, lastPart)
}

// FormatErrorContext creates a formatted context string for the AI
func FormatErrorContext(ctx ErrorContext) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Command: %s\n", ctx.Command))
	sb.WriteString(fmt.Sprintf("Shell: %s\n", ctx.Shell))
	sb.WriteString(fmt.Sprintf("OS: %s\n", ctx.OS))

	if len(ctx.RecentCmds) > 0 {
		sb.WriteString("\nRecent command history:\n")
		for i, cmd := range ctx.RecentCmds {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, cmd))
		}
	}

	sb.WriteString("\nError output:\n")
	sb.WriteString(ctx.Error)

	return sb.String()
}

// Execute is the main entry point for the aifix command
func Execute(promptFile string, manualError string) error {
	conf, prompt, err := config.ReadAndParseConfig("config.yaml", promptFile)
	if err != nil {
		return fmt.Errorf("error reading configuration: %v", err)
	}

	operatingSystem, shell := utils.DetectOSAndShell()
	prompt = utils.ReplacePlaceholders(prompt, operatingSystem, shell)

	// Create AI client based on provider
	var aiClient nlp.GAIClient
	if conf.Provider == "anthropic" {
		anthropicClient := nlp.CreateAnthropicClient(*conf)
		aiClient = &nlp.AnthropicClient{
			Client: anthropicClient,
			Prompt: prompt,
		}
	} else {
		openaiClient := nlp.CreateOpenAIClient(*conf)
		aiClient = &nlp.GoaiClient{
			Client: openaiClient,
			Prompt: prompt,
		}
	}

	var errorContext ErrorContext
	errorContext.Shell = shell
	errorContext.OS = operatingSystem
	errorContext.Timestamp = time.Now().Format(time.RFC3339)

	// Check if manual error was provided
	if manualError != "" {
		// User provided error directly
		errorContext.Command = "N/A (manual error input)"
		errorContext.Error = manualError
	} else {
		// Try shell integration first (environment files)
		cmdFile := os.Getenv("AIFIX_CMD_FILE")
		lastExit := os.Getenv("AIFIX_LAST_EXIT")

		if cmdFile != "" && lastExit != "0" && lastExit != "" {
			// Read last command from file
			cmdBytes, err := os.ReadFile(cmdFile)
			if err == nil && len(cmdBytes) > 0 {
				errorContext.Command = strings.TrimSpace(string(cmdBytes))

				// Get recent commands for context
				recentCmds, err := GetRecentCommands(shell, 3)
				if err == nil && len(recentCmds) > 1 {
					errorContext.RecentCmds = recentCmds[:len(recentCmds)-1]
				}

				// Re-run the command to capture its error output
				errorOutput, err := CaptureLastError(errorContext.Command, shell, operatingSystem)
				if err != nil {
					return fmt.Errorf("failed to capture error output: %v", err)
				}
				errorContext.Error = TruncateError(errorOutput, maxErrorLines)
				errorContext.ExitCode = 1
			} else {
				return fmt.Errorf("could not read command from file")
			}
		} else {
			// Fallback: Try to detect from shell history
			lastCmd, err := GetLastCommand(shell)
			if err != nil {
				return fmt.Errorf("no recent error detected. Usage: aifix [error message]\nError: %v", err)
			}
			errorContext.Command = lastCmd

			// Get recent commands for context
			recentCmds, err := GetRecentCommands(shell, 3)
			if err == nil && len(recentCmds) > 1 {
				// Exclude the last command itself
				errorContext.RecentCmds = recentCmds[:len(recentCmds)-1]
			}

			// Try to capture error by re-running the command
			errorOutput, err := CaptureLastError(lastCmd, shell, operatingSystem)
			if err != nil {
				return fmt.Errorf("last command succeeded. No errors to analyze\n\nTip: For automatic error detection, run: aifix -init-shell zsh")
			}
			errorContext.Error = TruncateError(errorOutput, maxErrorLines)
		}
	}

	// Format context for AI
	contextStr := FormatErrorContext(errorContext)

	// Process with AI
	response, err := aiClient.ProcessCommand(contextStr, *conf)
	if err != nil {
		return fmt.Errorf("error processing with AI: %v", err)
	}

	// Extract and display response
	result := response.Choices[0].Message.Content
	result = strings.TrimSpace(result)

	fmt.Println(result)

	return nil
}
