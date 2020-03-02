package game

type GameMode int

const (
	SURVIVAL GameMode = iota
	CREATIVE
	ADVENTURE
	SPECTATOR
)

func (g GameMode) Encoded(hardcore bool) byte {

	bit := 0
	if hardcore {
		bit = 0x8
	}

	return byte(g) | byte(bit)
}
