package ents

type entityLiving struct {
	entity

	health float64
}

func newEntityLiving() entityLiving {
	return entityLiving{entity: newEntity()}
}

func (e *entityLiving) GetHealth() float64 {
	return e.health
}

func (e *entityLiving) SetHealth(health float64) {
	e.health = health
}
