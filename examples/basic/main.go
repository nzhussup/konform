package main

import (
	"fmt"
	"log"

	"github.com/nzhussup/conform"
)

// This example demonstrates source precedence:
// YAML is loaded first, then ENV overrides matching fields.
type Config struct {
	App struct {
		Name        string
		Version     string
		Description string
		Author      string
		License     string
	}
	Database struct {
		Port        int
		URI         string
		PoolSize    int
		MaxLifetime string
	}
	Redis struct {
		Host string `env:"REDIS_HOST"`
		Port int    `env:"REDIS_PORT"`
	}
}

func main() {
	var cfg Config

	err := conform.Load(&cfg,
		conform.FromJSONFile("config.json"),
		conform.FromYAMLFile("config.yaml"),
		conform.FromEnv(),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)
}
