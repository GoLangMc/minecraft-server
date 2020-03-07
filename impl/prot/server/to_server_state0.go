package server

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/base"
)

// done

type PacketIHandshake struct {
	version int32

	host string
	port uint16

	State base.PacketState
}

func (p *PacketIHandshake) UUID() int32 {
	return 0x00
}

func (p *PacketIHandshake) Pull(reader buff.Buffer, conn base.Connection) {
	p.version = reader.PullVrI()

	p.host = reader.PullTxt()
	p.port = reader.PullU16()

	state := reader.PullVrI()

	p.State = base.PacketStateValueOf(int(state))
}
