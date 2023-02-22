package main

import (
	"userhistory/pkg/api"
	"userhistory/pkg/routes"
)

func main() {
	//corre el servidor
	go routes.Server()
	//corre el servidor REST
	go api.RestServer()
	//mantiene las goroutinas vivas
	<-make(chan bool)
}
