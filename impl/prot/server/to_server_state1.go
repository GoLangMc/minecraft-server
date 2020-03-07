package server

import "minecraft-server/impl/base"

// done

type PacketIRequest struct {
}

func (p *PacketIRequest) UUID() int32 {
	return 0x00
}

func (p *PacketIRequest) Pull(reader base.Buffer, conn base.Connection) {
	// no fields
}

type PacketIPing struct {
	Ping int64
}

func (p *PacketIPing) UUID() int32 {
	return 0x01
}

func (p *PacketIPing) Pull(reader base.Buffer, conn base.Connection) {
	p.Ping = reader.PullI64()
}
