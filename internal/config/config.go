package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	PostgresDSN string `yaml:"postgres_dsn" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := flag.String("config-path", "", "path file configure")

	flag.Parse()

	if *configPath == "" {
		log.Fatal("--config-path is not set")
	}

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", *configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}

	return &cfg
}
