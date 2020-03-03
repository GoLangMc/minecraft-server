package main

import (
	"flag"

	"github.com/fatih/color"

	"minecraft-server/impl"
	"minecraft-server/impl/conf"
)

func main() {
	color.NoColor = false

	server := impl.NewServer(MergeWithFlags(conf.DefaultServerConfig))
	server.Load()
}

func MergeWithFlags(c conf.ServerConfig) conf.ServerConfig {
	host := flag.String("host",
		conf.DefaultServerConfig.Network.Host,
		"the address this server will bind to")

	port := flag.Int("port",
		conf.DefaultServerConfig.Network.Port,
		"the port this server will bind to")

	flag.Parse()

	if *host != conf.DefaultServerConfig.Network.Host {
		c.Network.Host = *host
	}

	if *port != conf.DefaultServerConfig.Network.Port {
		c.Network.Port = *port
	}

	return c
}
