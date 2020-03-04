package mode

import (
	"minecraft-server/apis"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"

	impl_event "minecraft-server/impl/game/event"
)

func HandleState3(watcher util.Watcher, logger *logs.Logging, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) {

	go func() {
		for conn := range join {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnJoinEvent{Conn: conn})

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

	go func() {
		for conn := range quit {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnQuitEvent{Conn: conn})
		}
	}()
}
