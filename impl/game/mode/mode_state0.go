package mode

import (
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/prot/server"
)

/**
 * handshake
 */

func HandleState0(watcher util.Watcher) {

	watcher.SubAs(func(packet *server.PacketIHandshake, conn base.Connection) {
		conn.SetState(packet.State)
	})

}
