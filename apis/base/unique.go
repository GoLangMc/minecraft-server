package base

import "github.com/satori/go.uuid"

type Unique interface {
	UUID() uuid.UUID
}
