package base

import "minecraft-server/apis/uuid"

type Unique interface {
	UUID() uuid.UUID
}
