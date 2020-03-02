package mode

import (
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"
	"minecraft-server/impl/prot/states"
)

/**
 * handshake
 */

func HandleState0(watcher util.Watcher) {

	watcher.SubAs(func(packet *states.PacketIHandshake, conn base.Connection) {
		conn.SetState(packet.State)
	})

}
