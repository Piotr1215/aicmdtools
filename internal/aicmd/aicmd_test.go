package aicmd

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/piotr1215/aicmdtools/internal/config"
)

type MockExecutor struct {
	Err error
}

func (m *MockExecutor) Execute(command string) error {
	return m.Err
}

func TestShouldExecuteCommand(t *testing.T) {
	config := &config.Config{
		Safety: true,
	}

	// Custom reader to simulate user input
	input := strings.NewReader("n\n")

	result := shouldExecuteCommand(config, input)
	fmt.Printf("Result: %v\n", result) // Add this line to print the result

	if result != CmdDoNothing {
		t.Error("Expected command to be executed")
	}
}

func TestDefaultExecutor_Execute(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		e       *DefaultExecutor
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &DefaultExecutor{}
			if err := e.Execute(tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("DefaultExecutor.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_shouldExecuteCommand(t *testing.T) {
	type args struct {
		config *config.Config
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want CommandDecision
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldExecuteCommand(tt.args.config, tt.args.reader); got != tt.want {
				t.Errorf("shouldExecuteCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyCommandToClipboard(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := copyCommandToClipboard(tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("copyCommandToClipboard() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_appendCommandToHistory(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appendCommandToHistory(tt.args.command)
		})
	}
}

func TestExecute(t *testing.T) {
	type args struct {
		prompt_file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Execute(tt.args.prompt_file); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
