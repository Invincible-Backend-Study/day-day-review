package initializer

import (
	"day-day-review/internal/util"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
)

type DiscordConfig struct {
	Token string `yaml:"token"`
	Guild string `yaml:"guild"`
}

// LoadDiscordConfig discord.yaml 파일을 읽어와 DiscordConfig 구조체에 저장
func LoadDiscordConfig(filePath string) (*DiscordConfig, error) {
	content, err := util.LoadFile(filePath)
	if err != nil {
		log.Println("Failed to load discord.yaml:", err)
	}
	var config DiscordConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal discord.yaml: %w", err)
	}
	return &config, nil
}
