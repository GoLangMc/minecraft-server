package plugin

import "minecraft-server/impl/base"

type PluginMessage interface {
	Chan() string
	Data(writer base.Buffer, conn base.Connection)
}

type Brand struct {
	Name string
}

func (b *Brand) Chan() string {
	return "minecraft:brand"
}

func (b *Brand) Data(writer base.Buffer, conn base.Connection) {
	writer.PushTxt(b.Name)
}
