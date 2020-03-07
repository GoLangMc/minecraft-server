package values

import (
	"encoding/binary"

	"github.com/golangmc/minecraft-server/apis/base"
)

const (
	// ticks per second
	TPS = 20
	// milliseconds per tick
	MPT = 1_000 / TPS
)

var DefaultWorldHashedSeed = int64(binary.LittleEndian.Uint64(base.JavaSHA256HashLong(int64(base.JavaStringHashCode("North Carolina")))))
