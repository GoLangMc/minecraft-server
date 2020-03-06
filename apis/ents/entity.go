package ents

import "minecraft-server/apis/base"

type Entity interface {
	Sender
	base.Unique

	EntityUUID() int64
}
