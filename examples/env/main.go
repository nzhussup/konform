package main

import (
	"fmt"
	"log"

	"github.com/nzhussup/conform"
)

type Config struct {
	Port     int    `default:"8080" env:"PORT"`
	LogLevel string `default:"info" env:"LOG_LEVEL"`
	DBURL    string `env:"DATABASE_URL" required:"true"`
}

func main() {
	var cfg Config

	if err := conform.Load(&cfg, conform.FromEnv()); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)
}
