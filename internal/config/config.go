package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"dev" required:"true"`
	// StoragePath string `yaml:"storage_path" env-required:"true"` // for sqlite
	HTTPServer `yaml:"http_server"`
	DataBase   `yaml:"database"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"3s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-requierd:"true" env:"HTTP_SERVER_PASSWORD"`
}

type DataBase struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLmode  string `yaml:"sslmode" env-default:"disable"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file not exists:", configPath) // without `err`
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config:", err)
	}

	return &cfg
}

func GetDBURL() string {
	// TODO: should be move cfg outside the function
	cfg := MustLoad()
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DataBase.User,
		cfg.DataBase.Password,
		cfg.DataBase.Host,
		cfg.DataBase.Port,
		cfg.DataBase.DBName,
		cfg.DataBase.SSLmode)

}
