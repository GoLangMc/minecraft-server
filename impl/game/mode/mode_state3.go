package mode

import (
	"time"

	"minecraft-server/apis"
	"minecraft-server/apis/data"
	"minecraft-server/apis/game"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/task"
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"
	"minecraft-server/impl/data/plugin"
	"minecraft-server/impl/data/values"
	"minecraft-server/impl/prot/states"

	impl_event "minecraft-server/impl/game/event"
)

func HandleState3(watcher util.Watcher, logger *logs.Logging, tasking *task.Tasking, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) {

	tasking.EveryTime(10, time.Second, func(task *task.Task) {

		api := apis.MinecraftServer()

		// I hate this, add a functional method for player iterating
		for _, player := range api.Players() {

			// also probably add one that returns both the player and their connection
			conn := api.ConnByUUID(player.UUID())

			// keep player connection alive via keep alive
			conn.SendPacket(&states.PacketOKeepAlive{KeepAliveID: time.Now().UnixNano() / 1e6})
		}
	})

	watcher.SubAs(func(packet *states.PacketIKeepAlive, conn base.Connection) {
		logger.DataF("player %s is being kept alive", conn.Address())
	})

	go func() {
		for conn := range join {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnJoinEvent{Conn: conn})

			conn.SendPacket(&states.PacketOJoinGame{
				EntityID:      0,
				Hardcore:      false,
				GameMode:      game.CREATIVE,
				Dimension:     game.OVERWORLD,
				HashedSeed:    values.DefaultWorldHashedSeed,
				MaxPlayers:    10,
				LevelType:     game.DEFAULT,
				ViewDistance:  12,
				ReduceDebug:   false,
				RespawnScreen: false,
			})

			conn.SendPacket(&states.PacketOPluginMessage{
				Message: &plugin.Brand{
					Name: data.Translate("&b&lGo&2&lMC&r"),
				},
			})

			conn.SendPacket(&states.PacketOServerDifficulty{
				Difficulty: game.PEACEFUL,
				Locked:     true,
			})

			conn.SendPacket(&states.PacketOPlayerAbilities{
				Invulnerable: true,
				Flying:       true,
				AllowFlight:  true,
				InstantBuild: false,
				FlyingSpeed:  0.05, // default value
				FieldOfView:  0.1,  // default value
			})

			/*
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
