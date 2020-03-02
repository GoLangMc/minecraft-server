package ents

import (
	"github.com/satori/go.uuid"
)

type Entity interface {
	Sender

	UUID() uuid.UUID
}
