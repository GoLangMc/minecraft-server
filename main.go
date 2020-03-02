package main

import (
	"flag"

	"github.com/fatih/color"

	"minecraft-server/apis"
	"minecraft-server/impl"
)

func main() {
	color.NoColor = false

	var host string
	var port int

	parseParams(&host, &port)

	server := impl.NewServer(host, port)
	apis.CreateMinecraftServer(server)
	server.Load()
}

func parseParams(host *string, port *int) {
	tempHost := flag.String("host", "0.0.0.0", "the address this server will bind to")
	tempPort := flag.Int("port", 25565, "the port this server will bind to")

	flag.Parse()

	*host = *tempHost
	*port = *tempPort
}
