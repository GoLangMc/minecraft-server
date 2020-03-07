package server

import "minecraft-server/impl/base"

// done

type PacketILoginStart struct {
	PlayerName string
}

func (p *PacketILoginStart) UUID() int32 {
	return 0x00
}

func (p *PacketILoginStart) Pull(reader base.Buffer, conn base.Connection) {
	p.PlayerName = reader.PullTxt()
}

type PacketIEncryptionResponse struct {
	Secret []byte
	Verify []byte
}

func (p *PacketIEncryptionResponse) UUID() int32 {
	return 0x01
}

func (p *PacketIEncryptionResponse) Pull(reader base.Buffer, conn base.Connection) {
	p.Secret = reader.PullUAS()
	p.Verify = reader.PullUAS()
}

type PacketILoginPluginResponse struct {
	Message int32
	Success bool
	OptData []byte
}

func (p *PacketILoginPluginResponse) UUID() int32 {
	return 0x02
}

func (p *PacketILoginPluginResponse) Pull(reader base.Buffer, conn base.Connection) {
	p.Message = reader.PullVrI()
	p.Success = reader.PullBit()
	p.OptData = reader.UAS()[reader.InI():reader.Len()]
}
