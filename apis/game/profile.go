package game

import "github.com/golangmc/minecraft-server/apis/uuid"

type Profile struct {
	UUID uuid.UUID
	Name string

	Properties []*ProfileProperty
}

type ProfileProperty struct {
	Name      string
	Value     string
	Signature *string
}
