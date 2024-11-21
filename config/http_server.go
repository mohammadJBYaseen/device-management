package config

import (
	"device-management/router"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func StartHttpServer() *gin.Engine {
	Server := router.NewRouter(Controllers...)
	log.Printf("starting application %s localhost:%d ....", ApplicationProperties.Server.AppName, ApplicationProperties.Server.Port)
	if err := Server.Run(fmt.Sprintf("%s:%d", ApplicationProperties.Server.Host, ApplicationProperties.Server.Port)); err != nil {
		log.Fatalf("error running server: %v", err)
	}
	return Server
}
