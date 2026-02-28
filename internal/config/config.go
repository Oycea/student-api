package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Локал по дефолту вообще не безопасно
	Env string `yaml:"env" env-default:"local"`
	// Для prod конфига параметр 'connect_db' указвать не нужно, он будет взят из окружения
	Storage    string `yaml:"connect_db" env:"CONNECT_DB" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Adress      string        `yaml:"adress" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	// Путь к файлу конфига берём из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config path not set")
	}

	// Проверяем существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config

	// Читаем конфиг
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config from: %s", err)
	}

	return &cfg
}
