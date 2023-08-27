package main

import (
	"fmt"
	"shortUrl/model"
	"shortUrl/server"
)



func main() {
	fmt.Println("Setting up model and server")
	model.Setup()
	server.startServer()
}