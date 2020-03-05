package client

import (
	"fmt"

	"minecraft-server/impl/base"
)

type SkinParts struct {
	Cape bool
	Head bool
	Body bool
	ArmL bool
	ArmR bool
	LegL bool
	LegR bool
}

func (d *SkinParts) String() string {
	return fmt.Sprintf("Cape:%t Head:%t Body:%t ArmL:%t ArmR:%t LegL:%t LegR:%t", d.Cape, d.Head, d.Body, d.ArmL, d.ArmR, d.LegL, d.LegR)
}

func (d *SkinParts) Push(writer base.Buffer) {
	mask := byte(0)

	if d.Cape {
		mask |= 0x01
	}
	if d.Body {
		mask |= 0x02
	}
	if d.ArmL {
		mask |= 0x04
	}
	if d.ArmR {
		mask |= 0x08
	}
	if d.LegL {
		mask |= 0x10
	}
	if d.LegR {
		mask |= 0x20
	}
	if d.Head {
		mask |= 0x40
	}

	writer.PushByt(mask)
}

func (d *SkinParts) Pull(reader base.Buffer) {
	mask := reader.PullByt()

	if mask&0x01 != 0 {
		d.Cape = true
	}
	if mask&0x02 != 0 {
		d.Body = true
	}
	if mask&0x04 != 0 {
		d.ArmL = true
	}
	if mask&0x08 != 0 {
		d.ArmR = true
	}
	if mask&0x10 != 0 {
		d.LegL = true
	}
	if mask&0x20 != 0 {
		d.LegR = true
	}
	if mask&0x40 != 0 {
		d.Head = true
	}
}
