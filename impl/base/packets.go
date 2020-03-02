package base

import (
	"fmt"

	"minecraft-server/apis/util"
)

type PacketState int

func (state PacketState) String() string {
	switch state {
	case Shake:
		return "Shake"
	case Status:
		return "Status"
	case Login:
		return "Login"
	case Play:
		return "Play"
	default:
		panic(fmt.Errorf("no state for value: %d", state))
	}
}

const (
	Shake PacketState = iota
	Status
	Login
	Play
)

func GetState(state int) PacketState {
	switch state {
	case 0:
		return Shake
	case 1:
		return Status
	case 2:
		return Login
	case 3:
		return Play
	default:
		panic(fmt.Errorf("no state for value: %d", state))
	}
}

func (state PacketState) Next() PacketState {
	switch state {
	case Shake:
		return Status
	case Status:
		return Login
	case Login:
		return Play
	case Play:
		return Shake
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
	Pull(reader Buffer, conn Connection)
}

type PacketO interface {
	Packet

	// encode the server_data from the packet into this writer
	Push(writer Buffer, conn Connection)
}

type Packets interface {
	util.Watcher

	/*GetPacketM(uuid int32, state PacketState) (pid int32, cont bool)*/

	GetPacketI(uuid int32, state PacketState) PacketI

	/*GetPacketO(uuid int32, state PacketState) PacketO*/
}
