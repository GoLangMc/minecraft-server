package ents

import (
	"minecraft-server/apis/base"
	"minecraft-server/apis/data"
	impl "minecraft-server/impl/base"
	"minecraft-server/impl/prot/states"
)

type player struct {
	entityLiving

	name string

	online bool

	conn impl.Connection
}

func NewPlayer() player {
	return player{}
}

func (p *player) SendMessage(message ...interface{}) {
	packet := states.PacketOChatMessage{
		Message:         data.NewMessage(base.ConvertToString(message...)),
		MessagePosition: data.NormalChat,
	}

	p.conn.SendPacket(&packet)
}

func (p *player) GetIsOnline() bool {
	return p.online
}

func (p *player) SetIsOnline(state bool) {
	p.online = state
}
