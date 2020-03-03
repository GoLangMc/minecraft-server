package apis

import (
	"sync"

	"minecraft-server/apis/base"
	"minecraft-server/apis/cmds"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/task"
)

type Server interface {
	base.State

	Logging() *logs.Logging

	Command() *cmds.CommandManager

	Tasking() *task.Tasking
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
