package mode

import (
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"
	"minecraft-server/impl/data/status"
	"minecraft-server/impl/prot/client"
	"minecraft-server/impl/prot/server"
)

/**
 * status
 */

func HandleState1(watcher util.Watcher) {

	watcher.SubAs(func(packet *server.PacketIRequest, conn base.Connection) {
		response := client.PacketOResponse{Status: status.DefaultResponse()}
		conn.SendPacket(&response)
	})

	watcher.SubAs(func(packet *server.PacketIPing, conn base.Connection) {
		response := client.PacketOPong{Ping: packet.Ping}
		conn.SendPacket(&response)
	})

}
