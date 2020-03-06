package level

type Block interface {
	X() int
	Y() int
	Z() int

	Chunk() Chunk
	Level() Level

	GetBlockType() int64
	SetBlockType(value int64)
}
