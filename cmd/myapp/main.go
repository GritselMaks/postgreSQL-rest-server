package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/GritselMaks/postgreSQL-api-server/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file") //parse variable configPath
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil { // Start API server
		log.Fatal(err)
	}
}
