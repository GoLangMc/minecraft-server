package prot

import (
	"minecraft-server/apis/logs"
	"minecraft-server/apis/task"
	"minecraft-server/apis/util"
	"minecraft-server/impl/base"
	"minecraft-server/impl/game/mode"
	"minecraft-server/impl/prot/states"
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
				return &states.PacketIHandshake{}
			},
		},
		base.STATUS: {
			0x00: func() base.PacketI {
				return &states.PacketIRequest{}
			},
			0x01: func() base.PacketI {
				return &states.PacketIPing{}
			},
		},
		base.LOGIN: {
			0x00: func() base.PacketI {
				return &states.PacketILoginStart{}
			},
			0x01: func() base.PacketI {
				return &states.PacketIEncryptionResponse{}
			},
			0x02: func() base.PacketI {
				return &states.PacketILoginPluginResponse{}
			},
		},
		base.PLAY: {
			0x00: func() base.PacketI {
				return &states.PacketITeleportConfirm{}
			},
			0x01: func() base.PacketI {
				return &states.PacketIQueryBlockNBT{}
			},
			0x0F: func() base.PacketI {
				return &states.PacketIKeepAlive{}
			},
			0x03: func() base.PacketI {
				return &states.PacketIChatMessage{}
			},
		},
	}
}
