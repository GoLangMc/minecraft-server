package event

import (
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/data/plugin"
)

type PlayerConnJoinEvent struct {
	Conn base.PlayerAndConnection
}

type PlayerConnQuitEvent struct {
	Conn base.PlayerAndConnection
}

type PlayerPluginMessagePullEvent struct {
	Conn base.PlayerAndConnection

	Channel string
	Message plugin.Message
}
