package ents

import (
	"minecraft-server/apis/data/msgs"
	"minecraft-server/impl/prot/states"

	apis_base "minecraft-server/apis/base"
	impl_base "minecraft-server/impl/base"
)

type player struct {
	entityLiving

	name string

	online bool

	conn impl_base.Connection
}

func NewPlayer() player {
	return player{
		entityLiving: newEntityLiving(),
	}
}

func (p *player) SendMessage(message ...interface{}) {
	packet := states.PacketOChatMessage{
		Message:         *msgs.New(apis_base.ConvertToString(message...)),
		MessagePosition: msgs.NormalChat,
	}

	p.conn.SendPacket(&packet)
}

func (p *player) GetIsOnline() bool {
	return p.online
}

func (p *player) SetIsOnline(state bool) {
	p.online = state
}

func (p *player) SetConnection(conn impl_base.Connection) {
	p.conn = conn
}
