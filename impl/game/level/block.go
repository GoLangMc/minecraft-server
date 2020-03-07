package level

import (
	apis_level "github.com/golangmc/minecraft-server/apis/game/level"
)

type block struct {
	// these should always be level coordinates
	x int
	y int
	z int

	slice *slice
}

func (b *block) X() int {
	return b.x
}

func (b *block) Y() int {
	return b.y
}

func (b *block) Z() int {
	return b.z
}

func (b *block) Chunk() apis_level.Chunk {
	return b.slice.chunk
}

func (b *block) Level() apis_level.Level {
	return b.slice.chunk.level
}

func (b *block) GetBlockType() (value int) {
	value = b.slice.sliceBlockGet(sliceIndex(blockLevelToSlice(b.x, b.y, b.z)))
	return
}

func (b *block) SetBlockType(value int) {
	value = b.slice.sliceBlockSet(sliceIndex(blockLevelToSlice(b.x, b.y, b.z)), value)
	return
}
