package mode

import (
	"fmt"
	"time"

	"minecraft-server/apis"
	"minecraft-server/apis/data"
	"minecraft-server/apis/data/chat"
	"minecraft-server/apis/data/msgs"
	"minecraft-server/apis/game"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/task"
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"
	"minecraft-server/impl/data/client"
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

	watcher.SubAs(func(packet *states.PacketIPluginMessage, conn base.Connection) {
		api := apis.MinecraftServer()

		player := api.PlayerByConn(conn)
		if player == nil {
			return // log no player found?
		}

		api.Watcher().PubAs(impl_event.PlayerPluginMessagePullEvent{
			Conn: base.PlayerAndConnection{
				Connection: conn,
				Player:     player,
			},
			Channel: packet.Message.Chan(),
			Message: packet.Message,
		})
	})

	watcher.SubAs(func(packet *states.PacketIChatMessage, conn base.Connection) {
		api := apis.MinecraftServer()

		who := api.PlayerByConn(conn)
		out := msgs.
			New(who.Name()).SetColor(chat.White).
			Add(":").SetColor(chat.Gray).
			Add(" ").
			Add(chat.Translate(packet.Message)).SetColor(chat.White).
			AsText() // why not just use translate?

		api.Broadcast(out)
	})

	go func() {
		for conn := range join {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnJoinEvent{Conn: conn})

			conn.SendPacket(&states.PacketOJoinGame{
				EntityID:      int32(conn.EntityUUID()),
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
					Name: chat.Translate(fmt.Sprintf("&b%s&r &a%s&r", "GoLangMc", apis.MinecraftServer().ServerVersion())),
				},
			})

			conn.SendPacket(&states.PacketOServerDifficulty{
				Difficulty: game.PEACEFUL,
				Locked:     true,
			})

			conn.SendPacket(&states.PacketOPlayerAbilities{
				Abilities: client.PlayerAbilities{
					Invulnerable: true,
					Flying:       true,
					AllowFlight:  true,
					InstantBuild: false,
				},
				FlyingSpeed: 0.05, // default value
				FieldOfView: 0.1,  // default value
			})

			conn.SendPacket(&states.PacketOHeldItemChange{
				Slot: client.SLOT_0,
			})

			conn.SendPacket(&states.PacketODeclareRecipes{})

			conn.SendPacket(&states.PacketOPlayerLocation{
				ID: 0,
				Location: data.Location{
					PositionF: data.PositionF{
						X: 0,
						Y: 100,
						Z: 0,
					},
					RotationF: data.RotationF{
						AxisX: 0,
						AxisY: 0,
					},
				},
				Relative: client.Relativity{},
			})

			conn.SendPacket(&states.PacketOPlayerInfo{
				Action: client.AddPlayer,
				Values: []client.PlayerInfo{
					&client.PlayerInfoAddPlayer{Player: conn.Player},
				},
			})
		}
	}()

	go func() {
		for conn := range quit {
			apis.MinecraftServer().Watcher().PubAs(impl_event.PlayerConnQuitEvent{Conn: conn})
		}
	}()
}
