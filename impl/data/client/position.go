package client

import "minecraft-server/impl/base"

type Relativity struct {
	X bool
	Y bool
	Z bool

	AxisX bool
	AxisY bool
}

func (r *Relativity) Push(writer base.Buffer) {
	mask := byte(0)

	if r.X {
		mask |= 0x01
	}
	if r.Y {
		mask |= 0x02
	}
	if r.Z {
		mask |= 0x04
	}

	// the fact that these are flipped deeply bothers me.
	if r.AxisY {
		mask |= 0x08
	}
	if r.AxisX {
		mask |= 0x10
	}

	writer.PushByt(mask)
}
