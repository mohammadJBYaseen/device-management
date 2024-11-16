package config

import (
	"device-management/router"
	"fmt"
	"log"
)

func StartHttpServer() {
	Server := router.NewRouter(Controllers...)
	log.Printf("starting application %s localhost:%d ....", ApplicationProperties.Server.AppName, ApplicationProperties.Server.Port)
	if err := Server.Run(fmt.Sprintf("localhost:%d", ApplicationProperties.Server.Port)); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
