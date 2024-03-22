package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment    string        `yaml:"environment" penv:"ENV" required:"true"`
	Conn           string        `yaml:"conn" penv:"CONN" required:"true"`
	GRPC           GRPCConfig    `yaml:"grpc"`
	MigrationsPath string        `yaml:"migrations_path" penv:"MIGRATIONS_PATH"`
	TokenTTL       time.Duration `yaml:"token_ttl" penv:"TOKEN_TTL" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" penv:"GRPC_PORT"`
	Timeout time.Duration `yaml:"timeout" penv:"GRPC_TIMEOUT"`
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
