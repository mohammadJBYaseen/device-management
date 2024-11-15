package main

import (
	"device-management/config"
	"device-management/controller"
	"device-management/router"
	"fmt"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	myConfig := config.Load()
	server := router.NewRouter(controller.NewDeviceController())
	fmt.Printf("starting application %s localhost:%d ....", myConfig.Server.AppName, myConfig.Server.Port)
	if err := server.Run(fmt.Sprintf("localhost:%d", myConfig.Server.Port)); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
