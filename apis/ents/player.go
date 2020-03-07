package ents

import "minecraft-server/apis/game"

type Player interface {
	EntityLiving

	GetIsOnline() bool
	SetIsOnline(state bool)

	GetProfile() *game.Profile
}
