package event

import "minecraft-server/apis/ents"

type PlayerJoinEvent struct {
	Player ents.Player
}

type PlayerQuitEvent struct {
	Player ents.Player
}
