package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	DatabaseURL string `yaml:"database_url" env:"DATABASE_URL"`
	AliasLength int    `env:"ALIAS_LENGTH" env-default:"5"`
	HTTPServer  `yaml:"http_server"`
	Author      `yaml:"auth"`
}

type Author struct {
	User     string `yaml:"user" env:"AUTH_USER"`
	Password string `yaml:"password" env:"AUTH_PASSWORD"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	if cfg.AliasLength <= 0 {
		log.Fatal("ALIAS_LENGTH must be greater than 0")
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
		log.Printf("Using DATABASE_URL from environment")
	}

	if user := os.Getenv("AUTH_USER"); user != "" {
		cfg.Author.User = user
		log.Println("Using AUTH_USER from environment")
	}

	if password := os.Getenv("AUTH_PASSWORD"); password != "" {
		cfg.Author.Password = password
		log.Println("Using AUTH_PASSWORD from environment")
	}

	if cfg.Author.User == "" || cfg.Author.Password == "" {
		log.Fatal("Username and Password is empty")
	}

	return &cfg
}
