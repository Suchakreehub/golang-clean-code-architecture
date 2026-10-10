package config

import (
	"fmt"
	"log"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var Config = koanf.New(".")

// Pass Environment Variables to Go Application
func Load(envpath string) {
	if err := Config.Load(file.Provider(envpath), dotenv.Parser()); err != nil {
		log.Printf("loading config: %v", err)
	}

	Config.Load(env.Provider("GOAPP_", ".", func(s string) string {
		return s
	}), nil)

	fmt.Println("Loaded config application name = ", Config.String("GOAPP_NAME"))
}
