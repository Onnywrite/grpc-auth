package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment     string        `yaml:"environment" env:"ENV" required:"true"`
	Conn            string        `yaml:"conn" env:"CONN" required:"true"`
	GRPC            GRPCConfig    `yaml:"grpc"`
	MigrationsPath  string        `yaml:"migrations_path" env:"MIGRATIONS_PATH"`
	TokenTTL        time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"5m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl: " env:"REFRESH_TOKEN_TTL" env-default:"720h"`
	IdTokenTTL      time.Duration `yaml:"id_token_ttl: " env:"ID_TOKEN_TTL" env-default:"720h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env:"GRPC_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT"`
}

func MustLoad() *Config {
	var configPath string
	flag.StringVar(&configPath, "config-path", "", "config file path")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	return MustLoadByPath(configPath)
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("config could not be loaded: " + err.Error())
	}

	return &cfg
}
