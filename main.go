package main

import (
	"device-management/config"
)

func main() {
	defer config.Unload()
	config.StartHttpServer()
}
