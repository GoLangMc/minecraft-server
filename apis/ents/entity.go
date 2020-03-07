package ents

import "github.com/golangmc/minecraft-server/apis/base"

type Entity interface {
	Sender
	base.Unique

	EntityUUID() int64
}
