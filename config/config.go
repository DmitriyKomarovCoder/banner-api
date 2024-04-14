package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Http            http          `yaml:"http"`
	Log             logCustom     `yaml:"log_file"`
	PG              postgres      `yaml:"postgres"`
	Redis           redis         `yaml:"redis"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type http struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type logCustom struct {
	Path string `yaml:"path"`
}

type postgres struct {
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Port     int    `env:"DB_PORT"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	PoolMax  int32  `yaml:"pool_max"`
	URL      string
}

type redis struct {
	Address string `yaml:"host"`
	DB      int    `yaml:"db"`
}

func NewConfig(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return &cfg, err
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		fmt.Println(err)
		return &cfg, err
	}

	cfg.PG.URL = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PG.Host, cfg.PG.Port, cfg.PG.User, cfg.PG.Password, cfg.PG.Name)

	log.Println("Parsed Configuration")
	log.Println(cfg)
	return &cfg, nil
}
