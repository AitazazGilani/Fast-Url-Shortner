package main

import (
	"fmt"
	"github.com/AitazazGilani/Fast-Url-Shortner/backend/server"
)

func main() {
	fmt.Println("starting server at port 8080")
	server.StartServer()
}