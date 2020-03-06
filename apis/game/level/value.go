package level

const (
	ChunkW = 16
	ChunkH = 256
	ChunkL = 16

	SliceC = 16
	SliceH = ChunkH / SliceC

	SliceS = ChunkW * ChunkL * SliceH

	BitsPerBlock = 14
	MaxPaletteID = (1 << BitsPerBlock) - 1
)
