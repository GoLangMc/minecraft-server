package base

import "minecraft-server/apis/ents"

type PlayerAndConnection struct {
	Connection
	ents.Player
}
