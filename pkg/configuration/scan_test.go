package configuration

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	StringField string        `env:"STRING_FIELD" flag:"s" jsonConfig:"string_field" default:"default_value"`
	IntField    int           `env:"INT_FIELD" jsonConfig:"int_field" default:"50"`
	Int64Field  int64         `env:"INT64_FIELD" jsonConfig:"int64_field" default:"100"`
	BoolField   bool          `env:"BOOL_FIELD" jsonConfig:"bool_field"`
	Duration    time.Duration `env:"DURATION_FIELD" default:"3s"`
	private     string
}

func TestScanConfig_env(t *testing.T) {
	var config testConfig
	os.Setenv("STRING_FIELD", "value")
	defer os.Unsetenv("STRING_FIELD")
	err := ScanConfig(&config, []string{})
	assert.NoError(t, err)
	assert.Equal(t, "value", config.StringField)
}

func TestScanConfig_flag(t *testing.T) {
	var config testConfig

	err := ScanConfig(&config, []string{"-s=value2"})
	require.NoError(t, err)

	assert.Equal(t, "value2", config.StringField)
}

func TestScanConfig_default(t *testing.T) {
	var config testConfig

	err := ScanConfig(&config, []string{})
	require.NoError(t, err)

	assert.Equal(t, "default_value", config.StringField)
}

func TestScanConfig_jsonConfig(t *testing.T) {
	var config testConfig
	os.Setenv("CONFIG", getTestJSONFilePath())
	defer os.Unsetenv("CONFIG")

	err := ScanConfig(&config, []string{})
	require.NoError(t, err)

	assert.Equal(t, "json_value", config.StringField)
}

func TestScanConfig_EnvFlagPriority(t *testing.T) {
	var config testConfig
	os.Setenv("STRING_FIELD", "value")
	defer os.Unsetenv("STRING_FIELD")

	err := ScanConfig(&config, []string{"-s=value2"})
	require.NoError(t, err)

	assert.NotEqual(t, "value2", config.StringField, "Incorrect priority: flag must be less important than ENV")
	assert.Equal(t, "value", config.StringField, "Incorrect value")
}

func TestScanConfig_FlagConfigPriority(t *testing.T) {
	var config testConfig
	os.Setenv("CONFIG", getTestJSONFilePath())
	defer os.Unsetenv("CONFIG")
	err := ScanConfig(&config, []string{"-s=value2"})
	require.NoError(t, err)

	assert.NotEqual(t, "json_value", config.StringField, "Incorrect priority: config must be less important than flag")
	assert.Equal(t, "value2", config.StringField, "Incorrect value")
}

func TestScanConfig_int(t *testing.T) {
	var config testConfig
	os.Setenv("INT_FIELD", "50")
	defer os.Unsetenv("INT_FIELD")
	err := ScanConfig(&config, []string{})
	assert.NoError(t, err)
	assert.Equal(t, 50, config.IntField)
}

func TestScanConfig_int64(t *testing.T) {
	var config testConfig
	os.Setenv("INT64_FIELD", "1000")
	defer os.Unsetenv("INT64_FIELD")
	err := ScanConfig(&config, []string{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), config.Int64Field)
}

func TestScanConfig_bool(t *testing.T) {
	var config testConfig
	os.Setenv("BOOL_FIELD", "true")
	defer os.Unsetenv("BOOL_FIELD")
	err := ScanConfig(&config, []string{})
	assert.NoError(t, err)
	assert.Equal(t, true, config.BoolField)
}

func TestScanConfig_boolFalse(t *testing.T) {
	var config testConfig
	os.Setenv("BOOL_FIELD", "false")
	defer os.Unsetenv("BOOL_FIELD")
	err := ScanConfig(&config, []string{})
	assert.NoError(t, err)
	assert.Equal(t, false, config.BoolField)
}

func TestScanConfig_duration(t *testing.T) {
	var config testConfig
	os.Setenv("DURATION_FIELD", "1s")
	defer os.Unsetenv("DURATION_FIELD")
	err := ScanConfig(&config, []string{})
	assert.NoError(t, err)
	assert.Equal(t, 1*time.Second, config.Duration)
}

func getTestJSONFilePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to retrieve caller information")
	}

	dir := filepath.Dir(filename)
	return filepath.Join(dir, "./test_config.json")
}
