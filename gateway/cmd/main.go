package main

import (
	"log"
	"modular_chassis/gateway/api/services"

	_ "modular_chassis/carpay/pkg/api"
)

func main() {
	err := services.GetHTTPServer().Run()
	if err != nil {
		log.Fatal(err)
	}
}
