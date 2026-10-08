package service

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Xuanwo/go-locale"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/text/language"
)

type Config struct {
	Theme     string   `json:"theme"`
	Locale    string   `json:"locale"`
	Favorites []string `json:"favorites"`
}

func defaultConfig() *Config {
	tag, err := locale.Detect()
	if err != nil {
		slog.Warn("Failed to detect locale, defaulting to English", "error", err)
	}
	tag = language.English

	return &Config{
		Theme:  "dark",
		Locale: tag.String(),
	}
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
		slog.Warn("Config does not exist, creating one with default values", "path", configFilePath)
		if err := os.MkdirAll(configDirPath, 0755); err != nil {
			return err
		}
		file, err := os.Create(configFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		s.cfg = defaultConfig()

		if err := json.MarshalWrite(file, s.cfg, jsontext.Multiline(true)); err != nil {
			return err
		}

		return nil
	}

	slog.Info("Config exists, loading values", "path", configFilePath)

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

	slog.Info("Saving config", "path", configFilePath)

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

func (s *ConfigService) GetFavorites() (Stations, error) {
	if len(s.cfg.Favorites) == 0 {
		return Stations{}, nil
	}

	stations := make(Stations, 0, len(s.cfg.Favorites))
	for _, favorite := range s.cfg.Favorites {
		req, err := http.NewRequest(
			http.MethodGet,
			radioBrowserURL+"/json/stations/byuuid/"+favorite,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create network request for favorite %q: %w", favorite, err)
		}
		req.Header.Set("User-Agent", getUserAgent())

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("network request failed for favorite %q: %w", favorite, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code for favorite %q: %d", favorite, resp.StatusCode)
		}

		var favoriteStations Stations
		err = json.UnmarshalRead(resp.Body, &favoriteStations)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to decode response for favorite %q: %w", favorite, err)
		}
		stations = append(stations, favoriteStations...)
	}

	return stations, nil
}

func (s *ConfigService) SetFavorites(favorites Stations) {
	s.cfg.Favorites = make([]string, len(favorites))
	for i, station := range favorites {
		s.cfg.Favorites[i] = station.StationUUID
	}
}

func (s *ConfigService) SetTheme(theme string) {
	s.cfg.Theme = theme
}

func (s *ConfigService) SetLocale(locale string) {
	s.cfg.Locale = locale
}
