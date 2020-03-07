package ents

import "github.com/golangmc/minecraft-server/apis/uuid"

var entityCounter = int64(0)

type entity struct {
	name string
	uuid uuid.UUID

	entityID int64
}

func newEntity() entity {
	id := entityCounter
	entityCounter++

	return entity{entityID: id}
}

func (e *entity) Name() string {
	return e.name
}

func (e *entity) UUID() uuid.UUID {
	return e.uuid
}

func (e *entity) SetName(name string) {
	e.name = name
}

func (e *entity) SetUUID(uuid uuid.UUID) {
	e.uuid = uuid
}

func (e *entity) SendMessage(message ...interface{}) {
	// nothing
}

func (e *entity) EntityUUID() int64 {
	return e.entityID
}
