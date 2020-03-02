package mode

import (
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"
	"minecraft-server/impl/data/status"
	"minecraft-server/impl/prot/states"
)

/**
 * status
 */

func HandleState1(watcher util.Watcher) {

	watcher.SubAs(func(packet *states.PacketIRequest, conn base.Connection) {
		response := states.PacketOResponse{Status: status.DefaultResponse()}
		conn.SendPacket(&response)
	})

	watcher.SubAs(func(packet *states.PacketIPing, conn base.Connection) {
		response := states.PacketOPong{Ping: packet.Ping}
		conn.SendPacket(&response)
	})

}
