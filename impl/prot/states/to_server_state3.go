package states

import (
	"minecraft-server/apis/data"
	"minecraft-server/apis/game"
	"minecraft-server/impl/base"
	"minecraft-server/impl/data/client"
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

type PacketIClientStatus struct {
	Action client.StatusAction
}

func (p *PacketIClientStatus) UUID() int32 {
	return 0x04
}

func (p *PacketIClientStatus) Pull(reader base.Buffer, conn base.Connection) {
	p.Action = client.StatusAction(reader.PullVrI())
}

type PacketIClientSettings struct {
	Locale       string
	ViewDistance byte
	ChatMode     client.ChatMode
	ChatColors   bool // if false, strip messages of colors before sending
	SkinParts    client.SkinParts
	MainHand     client.MainHand
}

func (p *PacketIClientSettings) UUID() int32 {
	return 0x05
}

func (p *PacketIClientSettings) Pull(reader base.Buffer, conn base.Connection) {
	p.Locale = reader.PullTxt()
	p.ViewDistance = reader.PullByt()
	p.ChatMode = client.ChatMode(reader.PullVrI())
	p.ChatColors = reader.PullBit()

	parts := client.SkinParts{}
	parts.Pull(reader)

	p.SkinParts = parts
	p.MainHand = client.MainHand(reader.PullVrI())
}

type PacketIPlayerAbilities struct {
	Abilities   client.PlayerAbilities
	FlightSpeed float32
	GroundSpeed float32
}

func (p *PacketIPlayerAbilities) UUID() int32 {
	return 0x19
}

func (p *PacketIPlayerAbilities) Pull(reader base.Buffer, conn base.Connection) {
	abilities := client.PlayerAbilities{}
	abilities.Pull(reader)

	p.Abilities = abilities

	p.FlightSpeed = reader.PullF32()
	p.GroundSpeed = reader.PullF32()
}

type PacketIPlayerPosition struct {
	Position data.PositionF
	OnGround bool
}

func (p *PacketIPlayerPosition) UUID() int32 {
	return 0x11
}

func (p *PacketIPlayerPosition) Pull(reader base.Buffer, conn base.Connection) {
	p.Position = data.PositionF{
		X: reader.PullF64(),
		Y: reader.PullF64(),
		Z: reader.PullF64(),
	}

	p.OnGround = reader.PullBit()
}

type PacketIPlayerLocation struct {
	Location data.Location
	OnGround bool
}

func (p *PacketIPlayerLocation) UUID() int32 {
	return 0x12
}

func (p *PacketIPlayerLocation) Pull(reader base.Buffer, conn base.Connection) {
	p.Location = data.Location{
		PositionF: data.PositionF{
			X: reader.PullF64(),
			Y: reader.PullF64(),
			Z: reader.PullF64(),
		},
		RotationF: data.RotationF{
			AxisX: reader.PullF32(),
			AxisY: reader.PullF32(),
		},
	}

	p.OnGround = reader.PullBit()
}

type PacketIPlayerRotation struct {
	Rotation data.RotationF
	OnGround bool
}

func (p *PacketIPlayerRotation) UUID() int32 {
	return 0x13
}

func (p *PacketIPlayerRotation) Pull(reader base.Buffer, conn base.Connection) {
	p.Rotation = data.RotationF{
		AxisX: reader.PullF32(),
		AxisY: reader.PullF32(),
	}

	p.OnGround = reader.PullBit()
}
