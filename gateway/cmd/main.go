package main

import (
	"log"
	"modular_chassis/gateway/api/business"

	_ "modular_chassis/finanace/pkg/api"
	_ "modular_chassis/telecom/pkg/api"
)

func main() {
	err := business.GetHTTPServer().Run()
	if err != nil {
		log.Fatal(err)
	}
}
