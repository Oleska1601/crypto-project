package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		Web  `yaml:"web"`
		DB   `yaml:"db"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	HTTP struct {
		Port string `yaml:"port"`
	}

	Log struct {
		Level string `yaml:"level"`
	}

	Web struct {
		Path string `yaml:"path"`
	}

	DB struct {
		Path string `yaml:"path"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig("config/config.yml", cfg)
	if err != nil {
		return nil, err
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
