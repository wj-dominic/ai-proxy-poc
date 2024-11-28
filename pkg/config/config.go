package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "config.yaml"

type Config struct {
	LB struct {
		Port      int    `yaml:"port"`
		Algorithm string `yaml:"algorithm"`
	} `yaml:"lb"`
	Nodes []struct {
		ID     string `yaml:"id"`
		URL    string `yaml:"url"`
		MaxRPM int    `yaml:"max_rpm"`
		MaxBPM int64  `yaml:"max_bpm"`
	} `yaml:"nodes"`
}

func LoadConfig() (*Config, error) {
	// TODO: 환경변수 등 사용해서 config 읽도록 처리 필요
	bytes, err := os.ReadFile(defaultConfigPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
