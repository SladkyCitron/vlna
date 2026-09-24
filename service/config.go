package service

type Config struct {
	Theme string `json:"theme"`
}

type ConfigService struct {
	cfg *Config
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		// Default config for now
		cfg: &Config{
			Theme: "dark",
		},
	}
}

func (s *ConfigService) GetConfig() *Config {
	return s.cfg
}

func (s *ConfigService) SetTheme(theme string) {
	s.cfg.Theme = theme
}
