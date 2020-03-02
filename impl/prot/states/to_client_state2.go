package states

import (
	"minecraft-server/apis/data"
	"minecraft-server/impl/base"
)

type PacketODisconnect struct {
	Reason data.Message
}

func (p *PacketODisconnect) UUID() int32 {
	return 0x00
}

func (p *PacketODisconnect) Push(writer base.Buffer, conn base.Connection) {
	message := p.Reason

	writer.PushTxt(message.AsJson())
}

type PacketOEncryptionRequest struct {
	Server string // unused?
	Public []byte
	Verify []byte
}

func (p *PacketOEncryptionRequest) UUID() int32 {
	return 0x01
}

func (p *PacketOEncryptionRequest) Push(writer base.Buffer, conn base.Connection) {
	writer.PushTxt(p.Server)
	writer.PushArr(p.Public, true)
	writer.PushArr(p.Verify, true)
}

type PacketOLoginSuccess struct {
	PlayerUUID string
	PlayerName string
}

func (p *PacketOLoginSuccess) UUID() int32 {
	return 0x02
}

func (p *PacketOLoginSuccess) Push(writer base.Buffer, conn base.Connection) {
	writer.PushTxt(p.PlayerUUID)
	writer.PushTxt(p.PlayerName)
}
