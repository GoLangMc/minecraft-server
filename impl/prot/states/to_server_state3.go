package states

import (
	"minecraft-server/apis/data"
	"minecraft-server/apis/game"
	"minecraft-server/impl/base"
	"minecraft-server/impl/data/plugin"
)

type PacketIKeepAlive struct {
	KeepAliveID int64
}

func (p *PacketIKeepAlive) UUID() int32 {
	return 0x0F
}

func (p *PacketIKeepAlive) Pull(reader base.Buffer, conn base.Connection) {
	p.KeepAliveID = reader.PullI64()
}

type PacketIChatMessage struct {
	Message string
}

func (p *PacketIChatMessage) UUID() int32 {
	return 0x03
}

func (p *PacketIChatMessage) Pull(reader base.Buffer, conn base.Connection) {
	p.Message = reader.PullTxt()
}

type PacketITeleportConfirm struct {
	TeleportID int32
}

func (p *PacketITeleportConfirm) UUID() int32 {
	return 0x00
}

func (p *PacketITeleportConfirm) Pull(reader base.Buffer, conn base.Connection) {
	p.TeleportID = reader.PullVrI()
}

type PacketIQueryBlockNBT struct {
	TransactionID int32
	Position      data.PositionI
}

func (p *PacketIQueryBlockNBT) UUID() int32 {
	return 0x01
}

func (p *PacketIQueryBlockNBT) Pull(reader base.Buffer, conn base.Connection) {
	p.TransactionID = reader.PullVrI()
	p.Position = reader.PullPos()
}

type PacketISetDifficulty struct {
	Difficult game.Difficulty
}

func (p *PacketISetDifficulty) UUID() int32 {
	return 0x02
}

func (p *PacketISetDifficulty) Pull(reader base.Buffer, conn base.Connection) {
	p.Difficult = game.DifficultyValueOf(reader.PullByt())
}

type PacketIPluginMessage struct {
	Message plugin.Message
}

func (p *PacketIPluginMessage) UUID() int32 {
	return 0x0B
}

func (p *PacketIPluginMessage) Pull(reader base.Buffer, conn base.Connection) {
	channel := reader.PullTxt()
	message := plugin.GetMessageForChannel(channel)

	if message == nil {
		return // log unregistered channel?
	}

	message.Pull(reader)

	p.Message = message
}
