package main

import (
	"flag"

	"github.com/pondparinya/go-gateways/internal/config"
	"github.com/pondparinya/go-gateways/pkg/log"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./configs", "path to the config file")
var flagStage = flag.String("stage", "local", "set working environment")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "APP_VERSION", Version)

	if err := config.LoadConfigs("APP", *flagConfig, *flagStage, config.APP); err != nil {
		panic(err)
	}

	logger.Info(config.APP.Port)

}
