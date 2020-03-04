package apis

import (
	"sync"

	"minecraft-server/apis/cmds"
	"minecraft-server/apis/ents"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/task"
	"minecraft-server/apis/util"

	apis_base "minecraft-server/apis/base"
	impl_base "minecraft-server/impl/base"
)

type Server interface {
	apis_base.State

	Logging() *logs.Logging

	Command() *cmds.CommandManager

	Tasking() *task.Tasking

	Watcher() util.Watcher

	Players() []ents.Player

	PlayerConnection(player *ents.Player) *impl_base.Connection
}

var instance *Server
var syncOnce sync.Once

func MinecraftServer() Server {
	if instance == nil {
		panic("server is unavailable")
	}

	return *instance
}

func SetMinecraftServer(server Server) {
	syncOnce.Do(func() {
		instance = &server
	})
}
