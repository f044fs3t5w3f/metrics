package configuration

import (
	"encoding/json"
	"os"
)

func loadJSONConfig(configPath string) (map[string]*string, error) {
	jsonConfig := make(map[string]*string)
	if configPath == "" {
		return jsonConfig, nil
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return nil, err
	}

	return jsonConfig, nil
}
