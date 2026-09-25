package service

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Config struct {
	Theme string `json:"theme"`
}

var defaultConfig = &Config{
	Theme: "dark",
}

type ConfigService struct {
	cfg *Config
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) getPath() (string, error) {
	configDirPath := os.Getenv("VLNA_CONFIG_DIR")
	if configDirPath == "" {
		ucd, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		configDirPath = filepath.Join(ucd, "vlna")
	}
	return filepath.Join(configDirPath, "config.json"), nil
}

func (s *ConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	configFilePath, err := s.getPath()
	if err != nil {
		return err
	}

	configDirPath := filepath.Dir(configFilePath)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// config file does not exist, create one with defaults
		if err := os.MkdirAll(configDirPath, 0755); err != nil {
			return err
		}
		file, err := os.Create(configFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		s.cfg = defaultConfig

		if err := json.MarshalWrite(file, s.cfg, jsontext.Multiline(true)); err != nil {
			return err
		}

		return nil
	}

	b, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &s.cfg); err != nil {
		return err
	}

	return nil
}

func (s *ConfigService) SaveConfig() error {
	configFilePath, err := s.getPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.MarshalWrite(file, s.cfg, jsontext.Multiline(true)); err != nil {
		return err
	}

	return nil
}

func (s *ConfigService) GetConfig() *Config {
	return s.cfg
}

func (s *ConfigService) SetTheme(theme string) {
	s.cfg.Theme = theme
}
