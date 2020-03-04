package states

import (
	"minecraft-server/apis/data"
	"minecraft-server/apis/game"
	"minecraft-server/impl/base"
	"minecraft-server/impl/data/plugin"
)

type PacketOChatMessage struct {
	Message         data.Message
	MessagePosition data.MessagePosition
}

func (p *PacketOChatMessage) UUID() int32 {
	return 0x0E
}

func (p *PacketOChatMessage) Push(writer base.Buffer, conn base.Connection) {
	message := p.Message

	if p.MessagePosition == data.HotBarText {
		message = data.NewMessage(message.AsText())
	}

	writer.PushTxt(message.AsJson())
	writer.PushByt(byte(p.MessagePosition))
}

type PacketOJoinGame struct {
	EntityID    int32
	Hardcore    bool
	GameMode    game.GameMode
	Dimension   game.Dimension
	Difficulty  game.Difficulty
	MaxPlayers  int
	LevelType   game.LevelType
	ReduceDebug bool
}

func (p *PacketOJoinGame) UUID() int32 {
	return 0x25
}

func (p *PacketOJoinGame) Push(writer base.Buffer, conn base.Connection) {
	writer.PushI32(p.EntityID)
	writer.PushByt(p.GameMode.Encoded(p.Hardcore /* pull this value from somewhere */))
	writer.PushI32(int32(p.Dimension))
	writer.PushByt(uint8(p.MaxPlayers))
	writer.PushTxt(p.LevelType.String())
	writer.PushVrI(12)
	writer.PushBit(p.ReduceDebug)
}

type PacketOPluginMessage struct {
	Message plugin.PluginMessage
}

func (p *PacketOPluginMessage) UUID() int32 {
	return 0x18
}

func (p *PacketOPluginMessage) Push(writer base.Buffer, conn base.Connection) {
	writer.PushTxt(p.Message.Chan())
	p.Message.Data(writer, conn)
}

type PacketOPlayerPositionAndLook struct {
	X float64
	Y float64
	Z float64

	Yaw   float32
	Pitch float32

	Flags byte

	ID int32
}

func (p *PacketOPlayerPositionAndLook) UUID() int32 {
	return 0x35
}

func (p *PacketOPlayerPositionAndLook) Push(writer base.Buffer, conn base.Connection) {
	writer.PushF64(p.X)
	writer.PushF64(p.Y)
	writer.PushF64(p.Z)

	writer.PushF32(p.Yaw)
	writer.PushF32(p.Pitch)

	writer.PushByt(p.Flags)

	writer.PushVrI(p.ID)
}

type PacketOKeepAlive struct {
	KeepAliveID int64
}

func (p *PacketOKeepAlive) UUID() int32 {
	return 0x21
}

func (p *PacketOKeepAlive) Push(writer base.Buffer, conn base.Connection) {
	writer.PushI64(p.KeepAliveID)
}
