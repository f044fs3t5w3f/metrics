package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadJsonConfig(t *testing.T) {
	jsonConfig, err := loadJsonConfig(
		getTestJsonFilePath(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, jsonConfig)
	assert.NotNil(t, jsonConfig["string_field"])
	assert.Equal(t, "json_value", *jsonConfig["string_field"])
}
func TestLoadJsonConfig_IncorrectFile(t *testing.T) {
	_, err := loadJsonConfig(
		"./incorrect.json",
	)
	assert.Error(t, err)
}
