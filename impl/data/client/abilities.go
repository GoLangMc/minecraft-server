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
		flags |= 1
	}
	if p.Flying {
		flags |= 2
	}
	if p.AllowFlight {
		flags |= 3
	}
	if p.InstantBuild {
		flags |= 4
	}

	writer.PushByt(flags)
}
