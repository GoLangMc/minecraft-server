package states

import (
	"encoding/json"

	"minecraft-server/impl/base"
	"minecraft-server/impl/data/status"
)

type PacketOResponse struct {
	Status status.Response
}

func (p *PacketOResponse) UUID() int32 {
	return 0x00
}

func (p *PacketOResponse) Push(writer base.Buffer, conn base.Connection) {
	if text, err := json.Marshal(p.Status); err != nil {
		panic(err)
	} else {
		writer.PushTxt(string(text))
	}
}

type PacketOPong struct {
	Ping int64
}

func (p *PacketOPong) UUID() int32 {
	return 0x01
}

func (p *PacketOPong) Push(writer base.Buffer, conn base.Connection) {
	writer.PushI64(p.Ping)
}
