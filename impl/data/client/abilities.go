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

func (p *PlayerAbilities) Pull(reader base.Buffer) {
	flags := reader.PullByt()

	if flags&0x01 != 0 {
		p.Invulnerable = true
	}

	if flags&0x02 != 0 {
		p.Flying = true
	}

	if flags&0x04 != 0 {
		p.AllowFlight = true
	}

	if flags&0x08 != 0 {
		p.InstantBuild = true
	}
}
