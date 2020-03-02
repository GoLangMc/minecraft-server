package apis

import (
	"minecraft-server/apis/cmds"
	"minecraft-server/impl"
)

var mc *minecraftServer

type minecraftServer struct {
	Server *impl.Server
}

func MinecraftServer() *minecraftServer {
	return mc
}

func CreateMinecraftServer(server *impl.Server) {
	mc = &minecraftServer{Server: server}
}

func RegisterServerCommand(command ...cmds.Command) {
	for _, c := range command {
		mc.Server.Command.RegisterCommand(c)
	}
}
