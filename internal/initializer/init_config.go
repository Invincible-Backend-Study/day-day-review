package initializer

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type DiscordConfig struct {
	Token string `yaml:"token"`
	Guild string `yaml:"guild"`
}

func LoadDiscordConfig(filePath string) (*DiscordConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open discord.yaml: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Failed to close discord.yaml:", err)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read discord.yaml: %w", err)
	}
	var config DiscordConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal discord.yaml: %w", err)
	}
	return &config, nil
}
