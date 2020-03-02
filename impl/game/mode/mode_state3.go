package mode

import (
	"minecraft-server/apis/logs"
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"
)

func HandleState3(watcher util.Watcher, logger *logs.Logging, join chan base.PlayerAndConnection) {

	go func() {
		for conn := range join {
			logger.InfoF("player %s logged in with uuid:%v", conn.Player.Name(), conn.Player.UUID())

			/*conn.SendPacket(&states.PacketOJoinGame{
				EntityID:    0,
				Hardcore:    false,
				GameMode:    game.CREATIVE,
				Dimension:   game.OVERWORLD,
				Difficulty:  game.PEACEFUL,
				MaxPlayers:  10,
				LevelType:   game.DEFAULT,
				ReduceDebug: false,
			})

			conn.SendPacket(&states.PacketOPluginMessage{
				Message: &plugin.Brand{
					Name: data.Translate("&b&lGo&2&lMC&r"),
				},
			})

			conn.SendPacket(&states.PacketOPlayerPositionAndLook{
				X:     0,
				Y:     100,
				Z:     0,
				Yaw:   0,
				Pitch: 0,
				Flags: 0,
				ID:    0,
			})*/
		}
	}()
}
