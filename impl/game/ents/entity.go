package ents

import "github.com/satori/go.uuid"

type entity struct {
	name string
	uuid uuid.UUID
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
