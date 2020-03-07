package level

import (
	"github.com/golangmc/minecraft-server/apis/buff"
	apis_level "github.com/golangmc/minecraft-server/apis/game/level"
	"github.com/golangmc/minecraft-server/impl/base"
)

type slice struct {
	index int

	chunk *chunk

	values *base.Compacter
}

func newSlice(chunk *chunk, index int) *slice {
	slice := &slice{
		index: index,

		chunk: chunk,

		values: base.NewCompacter(apis_level.BitsPerBlock, apis_level.SliceS),
	}

	return slice
}

func (s *slice) Index() int {
	return s.index
}

func (s *slice) Chunk() apis_level.Chunk {
	return s.chunk
}

func (s *slice) Level() apis_level.Level {
	return s.chunk.level
}

func (s *slice) GetBlock(x, y, z int) apis_level.Block {
	if x < 0 || x > 15 {
		panic("invalid x value for slice get block")
	}
	if y < 0 || y > 15 {
		panic("invalid y value for slice get block")
	}
	if z < 0 || z > 15 {
		panic("invalid z value for slice get block")
	}

	return &block{
		x: (s.chunk.x << 0x04) | x,
		y: (apis_level.SliceH * s.index) + y,
		z: (s.chunk.z << 0x04) | z,

		slice: s,
	}
}

func (s *slice) Push(writer buff.Buffer) {
	writer.PushI16(apis_level.SliceS) // full slice

	writer.PushByt(apis_level.BitsPerBlock)

	// the server is using the direct palette

	writer.PushVrI(int32(len(s.values.Values)))

	for _, value := range s.values.Values {
		writer.PushI64(value)
	}
}

func (s *slice) sliceBlockGet(index int) int {
	return s.values.Get(index)
}

func (s *slice) sliceBlockSet(index int, value int) int {
	return s.values.Set(index, value)
}

func (s *slice) fill(value int) {
	for y := 0; y < apis_level.SliceH; y++ {
		s.layer(y, value)
	}
}

func (s *slice) layer(index int, value int) {
	for x := 0; x < apis_level.ChunkW; x++ {
		for z := 0; z < apis_level.ChunkL; z++ {
			s.sliceBlockSet(sliceIndex(x, index, z), value)
		}
	}
}
