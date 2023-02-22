package main

import (
	"log"
	"userhistory/pkg/api"
	"userhistory/pkg/routes"
)

func main() {
	go routes.Server()
	go api.RestServer()
	log.Println("running")
	<-make(chan bool)
}
