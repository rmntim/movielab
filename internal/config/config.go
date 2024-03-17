package config

import (
	"errors"
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string           `yaml:"env" env-required:"true"`
	DB         DBConfig         `yaml:"db"`
	HTTPServer HTTPServerConfig `yaml:"http_server"`
}

type DBConfig struct {
	Address  string `yaml:"address" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	Database string `yaml:"name" env-default:"movielab"`
	SslMode  string `yaml:"sslmode" env-default:"disable"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		panic(err)
	}
	return config
}

func Load() (*Config, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return nil, errors.New("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("config file doesn't exist: " + configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, errors.New("failed to read config: " + err.Error())
	}

	return &config, nil
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "config file path")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
