package ents

type EntityLiving interface {
	Entity

	GetHealth() float64
	SetHealth(health float64)
}
