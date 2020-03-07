package client

import (
	"minecraft-server/apis/ents"
	"minecraft-server/impl/base"
)

type PlayerInfoAction int32

const (
	AddPlayer PlayerInfoAction = iota
	UpdateGameMode
	UpdateLatency
	UpdateDisplayName
	RemovePlayer
)

type PlayerInfo interface {
	base.BufferPush
}

type PlayerInfoAddPlayer struct {
	Player ents.Player
}

func (p *PlayerInfoAddPlayer) Push(writer base.Buffer) {
	prof := p.Player.GetProfile()
	writer.PushUID(prof.UUID)
	writer.PushTxt(prof.Name)

	writer.PushVrI(int32(len(prof.Properties)))

	for _, prop := range prof.Properties {
		writer.PushTxt(prop.Name)
		writer.PushTxt(prop.Value)

		if prop.Signature == nil {
			writer.PushBit(false)
		} else {
			writer.PushBit(true)
			writer.PushTxt(*prop.Signature)
		}
	}

	writer.PushVrI(int32(p.Player.GetGameMode()))

	writer.PushVrI(0) // update this to the player's actual ping

	writer.PushBit(false) // update this to be whether the player has a custom display name or not, write that name as json if they do
}
