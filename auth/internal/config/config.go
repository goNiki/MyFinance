package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/subosito/gotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	DBConfig   `yaml:"db"`
	HTTPServer `yaml:"http_server"`
	TokenTTL   time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DBConfig struct {
	Host     string `yaml:"host" enf-default:"localhost"`
	Port     int64  `yaml:"port" enf-default:"5432"`
	User     string `yaml:"user" eng-default:"postgres"`
	Password string `yaml:"password" eng-default:"password"`
	DBname   string `yaml:"dbname" eng-default:"auth_db"`
	SSLmode  string `yaml:"sslmode" eng-default:"disable"`
}

func MustLoad() *Config {

	gotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatal("cannot read config %w", err)
	}
	return &cfg

}
