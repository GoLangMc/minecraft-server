package ents

type Player interface {
	EntityLiving

	GetIsOnline() bool
	SetIsOnline(state bool)
}
