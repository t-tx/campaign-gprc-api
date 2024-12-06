package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Log      `yaml:"logger"`
		Database `yaml:"database"`
	}
	Database struct {
		Path string `env-required:"true" yaml:"path" env:"DB_PATH"`
	}
	App struct {
		Name         string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		JwtSecretKey string `env-required:"true" yaml:"jwt_secret_key"    env:"JWT_SECRET_KEY"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}
)

var (
	globalConf *Config
)

func Get() *Config {
	return globalConf
}

func Load(path string) (*Config, error) {
	config := &Config{}
	if len(path) < 1 {
		return nil, errors.New("config file path is empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	globalConf = config
	return config, nil
}
