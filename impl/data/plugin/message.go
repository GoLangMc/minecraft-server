package plugin

import "minecraft-server/impl/base"

type Message interface {
	Chan() string

	Push(writer base.Buffer)
	Pull(reader base.Buffer)
}

var registry = createMessageRegistry()

type MessageRegistry struct {
	channels map[string]func() Message
}

func createMessageRegistry() MessageRegistry {
	registry := MessageRegistry{make(map[string]func() Message)}

	registry.channels["minecraft:brand"] = func() Message {
		return &Brand{}
	}

	return registry
}

func GetMessageForChannel(channel string) Message {
	creator := registry.channels[channel]
	if creator == nil {
		return nil
	}

	return creator()
}

// look, they're like cute little packets :D

type Brand struct {
	Name string
}

func (b *Brand) Chan() string {
	return "minecraft:brand"
}

func (b *Brand) Push(writer base.Buffer) {
	writer.PushTxt(b.Name)
}

func (b *Brand) Pull(reader base.Buffer) {
	b.Name = reader.PullTxt()
}
