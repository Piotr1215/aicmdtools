package aicmd

import (
	"fmt"
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

	if result {
		t.Error("Expected result to be false, but it was true")
	}
}
