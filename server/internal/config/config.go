package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	Storage    `yaml:"storage"`
	HTTPServer `yaml:"http_server"`
	JWT        `yaml:"jwt"`
}

type Storage struct {
	DbHost  string `yaml:"db_host" env-required:"true"`
	DbUser  string `yaml:"db_user" env-required:"true"`
	DbPort  int    `yaml:"db_port" env-required:"true"`
	DbPass  string `yaml:"db_pass" env-required:"true"`
	SslMode string `yaml:"ssl_mode" env-default:"prefer"`
	DbName  string `yaml:"db_name" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type JWT struct {
	JWT_ACCESS_SECRET  string `yaml:"jwt_access_secret" env-required:"true"`
	JWT_REFRESH_SECRET string `yaml:"jwt_refresh_secret" env-required:"true"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cnf Config
	if err := cleanenv.ReadConfig(configPath, &cnf); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cnf
}
