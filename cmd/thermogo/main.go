package main

import (
	"flag"
	"github.com/miterion/thermogo/pkg/server"
)

func main() {
	var port = flag.Int("port", 3000, "Port number for the server")
	flag.Parse()
	server.Run(*port)

}
