package base

import "github.com/golangmc/minecraft-server/apis/uuid"

type Unique interface {
	UUID() uuid.UUID
}
