package main

import (
	"log"
	"modular_chassis/gateway/api/business"

	_ "modular_chassis/authentication/pkg/api"
	_ "modular_chassis/carpay/pkg/api"
)

func main() {
	err := business.GetHTTPServer().Run()
	if err != nil {
		log.Fatal(err)
	}
}
