package main

import (
	"flag"

	"github.com/fatih/color"

	"github.com/golangmc/minecraft-server/impl"
	"github.com/golangmc/minecraft-server/impl/conf"
)

func main() {
	color.NoColor = false

	server := impl.NewServer(mergeWithFlags(conf.DefaultServerConfig))
	server.Load()
}

func mergeWithFlags(c conf.ServerConfig) conf.ServerConfig {
	host := flag.String("host",
		conf.DefaultServerConfig.Network.Host,
		"the address this server will bind to")

	port := flag.Int("port",
		conf.DefaultServerConfig.Network.Port,
		"the port this server will bind to")

	offline := flag.Bool("offline",
		conf.DefaultServerConfig.Offline,
		"do not authenticate against Mojang's session server")

	flag.Parse()

	if *host != conf.DefaultServerConfig.Network.Host {
		c.Network.Host = *host
	}

	if *port != conf.DefaultServerConfig.Network.Port {
		c.Network.Port = *port
	}

	if *offline != conf.DefaultServerConfig.Offline {
		c.Offline = *offline
	}

	return c
}
