package level

import (
	"github.com/golangmc/minecraft-server/apis/buff"
)

type Slice interface {
	buff.BufferPush

	Index() int

	Chunk() Chunk
	Level() Level

	// supports values x:[0:15] y:[0:15] z: [0:15]
	GetBlock(x, y, z int) Block
}
