package level

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	"github.com/golangmc/minecraft-server/apis/data/tags"
)

type Chunk interface {
	buff.BufferPush

	ChunkX() int
	ChunkZ() int

	Slices() []Slice

	Level() Level

	// supports values y:[0:15]
	GetSlice(y int) Slice

	// supports values x:[0:15] y:[0:255] z: [0:15]
	GetBlock(x, y, z int) Block

	HeightMapNbtCompound() *tags.NbtCompound
}
