package apis

import (
	"sync"

	"github.com/satori/go.uuid"

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

	ConnByUUID(uuid uuid.UUID) impl_base.Connection

	PlayerByUUID(uuid uuid.UUID) ents.Player

	PlayerByConn(conn impl_base.Connection) ents.Player

	ServerVersion() string
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
