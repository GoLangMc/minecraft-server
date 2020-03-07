package base

import "github.com/golangmc/minecraft-server/apis/ents"

type PlayerAndConnection struct {
	Connection
	ents.Player
}
