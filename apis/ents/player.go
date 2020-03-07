package ents

import "github.com/golangmc/minecraft-server/apis/game"

type Player interface {
	EntityLiving

	GetGameMode() game.GameMode
	SetGameMode(mode game.GameMode)

	GetIsOnline() bool
	SetIsOnline(state bool)

	GetProfile() *game.Profile
}
