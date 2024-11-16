package config

import (
	"device-management/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func StartHttpServer() *gin.Engine {
	Server := router.NewRouter(Controllers...)
	log.Printf("starting application %s localhost:%d ....", ApplicationProperties.Server.AppName, ApplicationProperties.Server.Port)
	if err := Server.Run(fmt.Sprintf("localhost:%d", ApplicationProperties.Server.Port)); err != nil {
		log.Fatalf("error running server: %v", err)
	}
	return Server
}
