package main

import (
	"flag"
	"simple-service/internal/server"
)

func main() {
	port := flag.Int("port", 5000, "specify the required port")
	flag.Parse()

	server.RunServer(*port)
}
