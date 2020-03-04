package states

import "minecraft-server/impl/base"

type PacketIKeepAlive struct {
	KeepAliveID int64
}

func (p *PacketIKeepAlive) UUID() int32 {
	return 0x0F
}

func (p *PacketIKeepAlive) Pull(reader base.Buffer, conn base.Connection) {
	p.KeepAliveID = reader.PullI64()
}
