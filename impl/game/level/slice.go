package level

import (
	apis_level "minecraft-server/apis/game/level"
	"minecraft-server/impl/base"
)

type slice struct {
	index int

	chunk *chunk

	values []int64
}

func newSlice(chunk *chunk, index int) *slice {
	slice := &slice{
		index: index,

		chunk: chunk,

		values: make([]int64, apis_level.SliceS, apis_level.SliceS),
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

func (s *slice) Push(writer base.Buffer) {
	writer.PushI16(apis_level.SliceS) // full slice

	writer.PushByt(apis_level.BitsPerBlock)

	// the server is using the direct palette

	writer.PushVrI(int32(len(s.values)))

	for _, value := range s.values {
		writer.PushI64(value)
	}
}

func (s *slice) sliceBlockGet(index int) int64 {
	bIndex := index * apis_level.BitsPerBlock

	sIndex := bIndex >> 0x06
	eIndex := (((index + 1) * apis_level.BitsPerBlock) - 1) >> 0x06

	uIndex := bIndex ^ (sIndex << 0x06)

	if sIndex == eIndex {
		return int64((uint64(s.values[sIndex]) >> uIndex) & apis_level.MaxPaletteID)
	}

	zIndex := 64 - uIndex

	return int64(uint64(s.values[sIndex]>>uIndex) | uint64((s.values[eIndex]<<zIndex))&apis_level.MaxPaletteID)
}

func (s *slice) sliceBlockSet(index int, value int64) int64 {
	bIndex := index * apis_level.BitsPerBlock

	sIndex := bIndex >> 0x06
	eIndex := (((index + 1) * apis_level.BitsPerBlock) - 1) >> 0x06

	uIndex := bIndex ^ (sIndex << 0x06)

	previousValue := int64((uint64(s.values[sIndex]) >> uIndex) & apis_level.MaxPaletteID)

	s.values[sIndex] = (s.values[sIndex] & ^(apis_level.MaxPaletteID << uIndex)) | ((value & apis_level.MaxPaletteID) << uIndex)

	if sIndex != eIndex {
		zIndex := 64 - uIndex
		pIndex := apis_level.BitsPerBlock - 1

		previousValue |= (s.values[eIndex] << zIndex) & apis_level.MaxPaletteID

		s.values[eIndex] = int64(((uint64(s.values[eIndex]) >> pIndex) << pIndex) | uint64((value&apis_level.MaxPaletteID)>>zIndex))
	}

	return previousValue
}

func (s *slice) fill(value int64) {
	for y := 0; y < apis_level.SliceH; y++ {
		s.layer(y, value)
	}
}

func (s *slice) layer(index int, value int64) {
	for x := 0; x < apis_level.ChunkW; x++ {
		for z := 0; z < apis_level.ChunkL; z++ {
			s.sliceBlockSet(sliceIndex(x, index, z), value)
		}
	}
}
