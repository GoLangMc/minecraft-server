package ents

import (
	"minecraft-server/apis/base"
)

type Sender interface {
	base.Named

	SendMessage(message ...interface{})
}
