package base

import (
	"fmt"

	"minecraft-server/apis/buff"
	"minecraft-server/apis/util"
)

type PacketState int

const (
	SHAKE PacketState = iota
	STATUS
	LOGIN
	PLAY
)

func ValueOfPacketState(s PacketState) int {
	return int(s)
}

func PacketStateValueOf(s int) PacketState {
	switch s {
	case 0:
		return SHAKE
	case 1:
		return STATUS
	case 2:
		return LOGIN
	case 3:
		return PLAY
	default:
		panic(fmt.Errorf("no state for value: %d", s))
	}
}

func (state PacketState) String() string {
	switch state {
	case SHAKE:
		return "Shake"
	case STATUS:
		return "Status"
	case LOGIN:
		return "Login"
	case PLAY:
		return "Play"
	default:
		panic(fmt.Errorf("no state for value: %d", state))
	}
}

func (state PacketState) Next() PacketState {
	switch state {
	case SHAKE:
		return STATUS
	case STATUS:
		return LOGIN
	case LOGIN:
		return PLAY
	case PLAY:
		return SHAKE
	default:
		panic(fmt.Errorf("no state for value: %d", state))
	}
}

type Packet interface {
	// the uuid of this packet
	UUID() int32
}

type PacketI interface {
	Packet

	// decode the server_data from the reader into this packet
	Pull(reader buff.Buffer, conn Connection)
}

type PacketO interface {
	Packet

	// encode the server_data from the packet into this writer
	Push(writer buff.Buffer, conn Connection)
}

type Packets interface {
	util.Watcher

	/*GetPacketM(uuid int32, state PacketState) (pid int32, cont bool)*/

	GetPacketI(uuid int32, state PacketState) PacketI

	/*GetPacketO(uuid int32, state PacketState) PacketO*/
}
