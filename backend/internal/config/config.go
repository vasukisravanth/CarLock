package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Port     string `json:"port"`
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`
	MongoDB struct {
		URI    string `json:"uri"`
		DBName string `json:"db_name"`
	} `json:"mongodb"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		paths := []string{
			filepath.Join(cwd, "config.json"),
			filepath.Join(cwd, "..", "config.json"),
			filepath.Join(cwd, "..", "..", "config.json"),
		}

		exePath, err := os.Executable()
		if err == nil {
			exeDir := filepath.Dir(exePath)
			paths = append(paths,
				filepath.Join(exeDir, "config.json"),
				filepath.Join(exeDir, "..", "config.json"),
				filepath.Join(exeDir, "..", "..", "config.json"),
			)
		}

		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				configPath = p
				break
			}
		}
	}

	if configPath == "" {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
