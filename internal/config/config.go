package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `yaml:"environment" env:"ENV" required:"true"`
	Conn        string `yaml:"conn" env:"CONN" required:"true"`

	Grpc         TransportConfig `yaml:"grpc"`
	Https        TransportConfig `yaml:"https"`
	AccessToken  TokenConfig     `yaml:"access_token"`
	RefreshToken TokenConfig     `yaml:"refresh_token"`
	DefaultUsers []UserConfig    `yaml:"default_users"`
}

type TransportConfig struct {
	Port    int           `yaml:"port" env:"PORT"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT"`
	// Cert    string        `yaml:"cert" env:"CERT"`
	// Key     string        `yaml:"key" env:"KEY"`
}

type TokenConfig struct {
	Secret   string        `yaml:"secret" env:"TOKEN_SECRET" required:"true"`
	TTL      time.Duration `yaml:"ttl" env:"TOKEN_TTL" env-default:"1h"`
	Issuer   string        `yaml:"issuer" env:"TOKEN_ISSUER" env-default:"localhost"`
	Audience string        `yaml:"audience" env:"TOKEN_AUDIENCE" env-default:"localhost"`
	Subject  string        `yaml:"subject" env:"TOKEN_SUBJECT" env-default:"localhost"`
}

type UserConfig struct {
	Nickname string       `yaml:"nickname" env:"USER_NICKNAME" required:"true"`
	Password string       `yaml:"password" env:"USER_PASSWORD" required:"true"`
	Email    *string      `yaml:"email" env:"USER_EMAIL"`
	Phone    *string      `yaml:"phone" env:"USER_PHONE"`
	Roles    []RoleConfig `yaml:"roles"`
}

// TODO: sync with Role struct
type RoleConfig struct {
	Name string `yaml:"name" env:"ROLE_NAME" required:"true"`
}

func MustLoad() *Config {
	var configPath string
	flag.StringVar(&configPath, "config-path", "", "config file path")
	flag.Parse()

	if configPath == "" {
		configPath = "/etc/sso/ignore-config.yaml"
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
