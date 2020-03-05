package client

import "minecraft-server/impl/base"

type PlayerAbilities struct {
	Invulnerable bool
	Flying       bool
	AllowFlight  bool
	InstantBuild bool // creative??
}

func (p *PlayerAbilities) Push(writer base.Buffer) {
	flags := byte(0)

	if p.Invulnerable {
		flags |= 0x01
	}
	if p.Flying {
		flags |= 0x02
	}
	if p.AllowFlight {
		flags |= 0x04
	}
	if p.InstantBuild {
		flags |= 0x08
	}

	writer.PushByt(flags)
}
