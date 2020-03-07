package ents

import (
	"github.com/golangmc/minecraft-server/apis/data/msgs"
	"github.com/golangmc/minecraft-server/apis/ents"
	"github.com/golangmc/minecraft-server/apis/game"

	"github.com/golangmc/minecraft-server/impl/prot/client"

	apis_base "github.com/golangmc/minecraft-server/apis/base"
	impl_base "github.com/golangmc/minecraft-server/impl/base"
)

type player struct {
	entityLiving

	prof *game.Profile

	online bool

	conn impl_base.Connection

	mode game.GameMode
}

func NewPlayer(prof *game.Profile, conn impl_base.Connection) ents.Player {
	player := &player{
		prof:         prof,
		entityLiving: newEntityLiving(),
	}

	player.SetName(prof.Name)
	player.SetUUID(prof.UUID)

	player.SetConn(conn)

	return player
}

func (p *player) SendMessage(message ...interface{}) {
	packet := client.PacketOChatMessage{
		Message:         *msgs.New(apis_base.ConvertToString(message...)),
		MessagePosition: msgs.NormalChat,
	}

	p.conn.SendPacket(&packet)
}

func (p *player) GetGameMode() game.GameMode {
	return p.mode
}

func (p *player) SetGameMode(mode game.GameMode) {
	p.mode = mode
}

func (p *player) GetIsOnline() bool {
	return p.online
}

func (p *player) SetIsOnline(state bool) {
	p.online = state
}

func (p *player) GetProfile() *game.Profile {
	return p.prof
}

func (p *player) SetConn(conn impl_base.Connection) {
	p.conn = conn
}
