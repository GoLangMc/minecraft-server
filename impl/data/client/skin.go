package client

import (
	"fmt"

	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/impl/mask"
)

type SkinParts struct {
	mask.Masking

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

func (d *SkinParts) Push(writer buff.Buffer) {
	flags := byte(0)

	d.Set(&flags, 0x01, d.Cape)
	d.Set(&flags, 0x02, d.Body)
	d.Set(&flags, 0x04, d.ArmL)
	d.Set(&flags, 0x08, d.ArmR)
	d.Set(&flags, 0x10, d.LegL)
	d.Set(&flags, 0x20, d.LegR)
	d.Set(&flags, 0x40, d.Head)

	writer.PushByt(flags)
}

func (d *SkinParts) Pull(reader buff.Buffer) {
	flags := reader.PullByt()

	d.Cape = d.Has(flags, 0x01)
	d.Body = d.Has(flags, 0x02)
	d.ArmL = d.Has(flags, 0x04)
	d.ArmR = d.Has(flags, 0x08)
	d.LegL = d.Has(flags, 0x10)
	d.LegR = d.Has(flags, 0x20)
	d.Head = d.Has(flags, 0x40)
}
