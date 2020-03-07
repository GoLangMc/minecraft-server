package ents

import (
	"github.com/golangmc/minecraft-server/apis/base"
)

type Sender interface {
	base.Named

	SendMessage(message ...interface{})
}
