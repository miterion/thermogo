package cmd

import (
	"flag"
	"github.com/miterion/thermogo/pkg/server"
)

func Execute() {
	var port = flag.Int("port", 3000, "Port number for the server")
	flag.Parse()
	server.Run(*port)

}
