package apis

import (
	"sync"

	"github.com/golangmc/minecraft-server/apis/cmds"
	"github.com/golangmc/minecraft-server/apis/ents"
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/apis/uuid"

	apis_base "github.com/golangmc/minecraft-server/apis/base"
	impl_base "github.com/golangmc/minecraft-server/impl/base"
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

	Broadcast(message string)
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
