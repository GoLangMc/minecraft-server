package level

import "minecraft-server/impl/base"

type Slice interface {
	base.BufferPush

	Index() int

	Chunk() Chunk
	Level() Level

	// supports values x:[0:15] y:[0:15] z: [0:15]
	GetBlock(x, y, z int) Block
}
