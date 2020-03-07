package level

import (
	apis_level "minecraft-server/apis/game/level"
	"minecraft-server/apis/uuid"
)

type level struct {
	name string
	uuid uuid.UUID

	chunks map[int64]*chunk
}

func NewLevel(name string) apis_level.Level {
	level := &level{
		name: name,
		uuid: uuid.NewUUID(),

		chunks: make(map[int64]*chunk),
	}

	return level
}

func (l *level) Name() string {
	return l.name
}

func (l *level) UUID() uuid.UUID {
	return l.uuid
}

func (l *level) Chunks() []apis_level.Chunk {
	chunks := make([]apis_level.Chunk, len(l.chunks), len(l.chunks))

	index := 0
	for _, chunk := range l.chunks {
		chunks[index] = chunk
		index++
	}

	return chunks
}

func (l *level) GetChunk(x, z int) apis_level.Chunk {
	return l.getChunk(x, z, true)
}

func (l *level) GetChunkIfLoaded(x, z int) apis_level.Chunk {
	return l.getChunk(x, z, false)
}

func (l *level) GetBlock(x, y, z int) apis_level.Block {
	return &block{
		x: x,
		y: y,
		z: z,

		slice: l.GetChunk(blockXZToChunkXZ(x, z)).GetSlice(blockYToSliceY(y)).(*slice),
	}
}

func (l *level) getChunk(x, z int, generate bool) apis_level.Chunk {
	idx := chunkIndex(x, z)

	cnk, con := l.chunks[idx]
	if con {
		return cnk
	}

	if generate {
		gen := newChunk(l, x, z)
		l.chunks[idx] = gen

		return gen
	}

	return nil
}

// generates chunks with the normal super-flat style
func GenSuperFlat(level apis_level.Level, size int) {
	for x := -size; x < size; x++ {
		for z := -size; z < size; z++ {
			chunk := level.GetChunk(x, z)

			for sliceY := 0; sliceY < apis_level.SliceC; sliceY++ {
				chunk.GetSlice(sliceY)
			}

			chunk.GetSlice(0).(*slice).layer(0, 7)
			chunk.GetSlice(0).(*slice).layer(1, 1)
			chunk.GetSlice(0).(*slice).layer(2, 1)
			chunk.GetSlice(0).(*slice).layer(3, 2)
		}
	}
}
