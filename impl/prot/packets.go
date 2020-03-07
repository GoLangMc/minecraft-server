package prot

import (
	"github.com/golangmc/minecraft-server/apis/logs"
	"github.com/golangmc/minecraft-server/apis/task"
	"github.com/golangmc/minecraft-server/apis/util"
	"github.com/golangmc/minecraft-server/impl/base"
	"github.com/golangmc/minecraft-server/impl/game/mode"
	"github.com/golangmc/minecraft-server/impl/prot/server"
)

type packets struct {
	util.Watcher

	logger  *logs.Logging
	packetI map[base.PacketState]map[int32]func() base.PacketI // UUID to I server_data

	join chan base.PlayerAndConnection
	quit chan base.PlayerAndConnection
}

func NewPackets(tasking *task.Tasking, join chan base.PlayerAndConnection, quit chan base.PlayerAndConnection) base.Packets {
	packets := &packets{
		Watcher: util.NewWatcher(),

		logger:  logs.NewLogging("protocol", logs.EveryLevel...),
		packetI: createPacketI(),
	}

	mode.HandleState0(packets)
	mode.HandleState1(packets)
	mode.HandleState2(packets, join)
	mode.HandleState3(packets, packets.logger, tasking, join, quit)

	return packets
}

func (p *packets) GetPacketI(uuid int32, state base.PacketState) base.PacketI {
	creator := p.packetI[state][uuid]
	if creator == nil {
		return nil
	}

	return creator()
}

func createPacketI() map[base.PacketState]map[int32]func() base.PacketI {
	return map[base.PacketState]map[int32]func() base.PacketI{
		base.SHAKE: {
			0x00: func() base.PacketI {
				return &server.PacketIHandshake{}
			},
		},
		base.STATUS: {
			0x00: func() base.PacketI {
				return &server.PacketIRequest{}
			},
			0x01: func() base.PacketI {
				return &server.PacketIPing{}
			},
		},
		base.LOGIN: {
			0x00: func() base.PacketI {
				return &server.PacketILoginStart{}
			},
			0x01: func() base.PacketI {
				return &server.PacketIEncryptionResponse{}
			},
			0x02: func() base.PacketI {
				return &server.PacketILoginPluginResponse{}
			},
		},
		base.PLAY: {
			0x00: func() base.PacketI {
				return &server.PacketITeleportConfirm{}
			},
			0x01: func() base.PacketI {
				return &server.PacketIQueryBlockNBT{}
			},
			0x02: func() base.PacketI {
				return &server.PacketISetDifficulty{}
			},
			0x03: func() base.PacketI {
				return &server.PacketIChatMessage{}
			},
			0x04: func() base.PacketI {
				return &server.PacketIClientStatus{}
			},
			0x05: func() base.PacketI {
				return &server.PacketIClientSettings{}
			},
			0x0B: func() base.PacketI {
				return &server.PacketIPluginMessage{}
			},
			0x0F: func() base.PacketI {
				return &server.PacketIKeepAlive{}
			},
			0x11: func() base.PacketI {
				return &server.PacketIPlayerPosition{}
			},
			0x12: func() base.PacketI {
				return &server.PacketIPlayerLocation{}
			},
			0x13: func() base.PacketI {
				return &server.PacketIPlayerRotation{}
			},
			0x19: func() base.PacketI {
				return &server.PacketIPlayerAbilities{}
			},
		},
	}
}
