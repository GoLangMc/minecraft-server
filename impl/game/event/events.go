package event

import "minecraft-server/impl/base"

type PlayerConnJoinEvent struct {
	Conn base.PlayerAndConnection
}

type PlayerConnQuitEvent struct {
	Conn base.PlayerAndConnection
}
