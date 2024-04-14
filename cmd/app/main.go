package main

import (
	"fmt"
	"log"

	"github.com/DmitriyKomarovCoder/banner-api/config"
	"github.com/DmitriyKomarovCoder/banner-api/internal/app"
	"github.com/joho/godotenv"
)

const path = "config/config.yaml"

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file")
		return
	}

	cfg, err := config.NewConfig(path)
	if err != nil {
		log.Fatalf("Config error %s", err)
	}

	app.Run(cfg)
}
