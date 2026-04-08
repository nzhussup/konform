package main

import (
	"fmt"
	"log"

	"github.com/nzhussup/konform"
)

type Config struct {
	AppName     string `validate:"required,min=3,max=100"`
	Version     string `validate:"required"`
	Description string `validate:"required,max_len=10"`
	Database    struct {
		Host     string `validate:"required"`
		Port     int    `validate:"required,min=1,max=65535,len=5"`
		Username string `validate:"required, min_len=3, max_len=50"`
		Password string `validate:"required,min_len=15"`
		Name     string `validate:"required"`
	}
	Features struct {
		EnableLogging bool
		EnableMetrics bool
	} `validate:"required"`
}

func main() {
	cfg := Config{}

	if err := konform.Load(&cfg, konform.FromJSONFile("config.json")); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%+v\n", cfg)
}
