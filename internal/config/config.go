package config

import (
	"flag"
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/config/lib/processconfig"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment    string        `yaml:"environment" required:"true"`
	Conn           string        `yaml:"conn" required:"true"`
	GRPC           GRPCConfig    `yaml:"grpc"`
	MigrationsPath string        `yaml:"migrations_path"`
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
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

	processconfig.PrioritizeEnvs(&cfg)

	return &cfg
}
