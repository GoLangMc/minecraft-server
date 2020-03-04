package game

import "fmt"

type Difficulty byte

const (
	PEACEFUL Difficulty = iota
	EASY
	NORMAL
	HARD
)

func (d Difficulty) String() string {
	switch d {
	case PEACEFUL:
		return "Peaceful"
	case EASY:
		return "Easy"
	case NORMAL:
		return "Normal"
	case HARD:
		return "Hard"
	default:
		panic(fmt.Errorf("no difficulty for id %d", byte(d)))
	}
}

func ValueOfDifficulty(d Difficulty) byte {
	return byte(d)
}

func DifficultyValueOf(id byte) Difficulty {
	switch id {
	case 0:
		return PEACEFUL
	case 1:
		return EASY
	case 2:
		return NORMAL
	case 3:
		return HARD
	default:
		panic(fmt.Errorf("no difficulty for id %d", id))
	}
}
