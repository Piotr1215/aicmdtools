package config_test

import (
	"testing"

	"github.com/piotr1215/aicmdtools/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestReadAndParseConfig(t *testing.T) {
	conf, prompt, err := config.ReadAndParseConfig("config.yaml", "comp-graph-prompt.txt")
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.NotEmpty(t, prompt)
}
